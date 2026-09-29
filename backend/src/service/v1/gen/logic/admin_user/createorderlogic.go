package admin_user

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/google/uuid"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.AdminCreateOrder, r *http.Request) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	filename, ok := r.MultipartForm.Value["filename[]"]
	if !ok {
		filename = []string{}
	}

	if len(filename) > int(config.BackendConfig.MySQL.WorkOrderFileLimit) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooManyFile, "文件太多"),
		}, nil
	}

	if utils.HasDuplicate(filename) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.RepeatFileName, "包含重复的文件名"),
		}, nil
	}

	file, ok := r.MultipartForm.File["file[]"]
	if !ok {
		file = []*multipart.FileHeader{}
	}

	if len(file) != len(filename) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "文件个数和文件名个数不匹配"),
		}, nil
	}

	fileByte := make([][]byte, 0, len(filename))
	for _, fh := range file {
		err := func() error {
			f, err := fh.Open()
			if err != nil {
				return err
			}
			defer utils.Close(f)

			if fh.Size >= FileMaxSize {
				return errors.Errorf("file too big")
			}

			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			fileByte = append(fileByte, b)
			return nil
		}()

		if err != nil {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.BadFormFile, errors.WarpQuick(err), "读取文件错误"),
			}, nil
		}
	}

	workOrderModel := db.NewWorkOrderModel(mysql.MySQLConn)
	workOrderCommunicateModel := db.NewWorkOrderCommunicateModel(mysql.MySQLConn)

	orderIDUUID, success := redis.GenerateUUIDMore(l.ctx, "order", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		workOrderModel := db.NewWorkOrderModel(mysql.MySQLConn)
		_, err := workOrderModel.FindOneByUidWithoutDelete(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CreateWorkOrderFail, "生成工单ID错误"),
		}, nil
	}

	orderID := orderIDUUID.String()

	order := &db.WorkOrder{
		Uid:    orderID,
		UserId: srcUser.Id,
		Title:  req.Title,
		Status: db.WorkOrderStatusWaitReply,
		LastReplyAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	if web.ID == warp.UserCenterWebsite {
		order.From = config.BackendConfig.User.ReadableName
		order.FromId = warp.UserCenterWebsite
	} else {
		bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
		allow, err := bannedModel.CheckAllow(l.ctx, srcUser.Id, web.ID, db.AllowMsg)
		if err != nil || !allow {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "用户关闭了通信授权许可"),
			}, nil
		}

		order.From = web.Name
		order.FromId = web.ID
	}

	res1, err := workOrderModel.Insert(l.ctx, order)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}
	order.Id, err = res1.LastInsertId()
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	communicate := &db.WorkOrderCommunicate{
		OrderId: order.Id,
		Content: req.Content,
		From:    db.CommunicateFromAdmin,
	}
	res2, err := workOrderCommunicateModel.Insert(l.ctx, communicate)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}
	communicate.Id, err = res2.LastInsertId()
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	go func(order *db.WorkOrder, communicate *db.WorkOrderCommunicate, filename []string, fileByte [][]byte) {
		for i, b := range fileByte {
			err := oss.UploadWorkOrderFile(order, communicate, filename[i], b, true)
			if err != nil {
				logger.Logger.Error("mysql error: %s", err.Error())
				continue
			}
		}

		// 不需要更新ws事件，因为这是新建的，没人监听
	}(order, communicate, filename, fileByte)

	audit.NewAdminAudit(user.Id, "管理员创建工单（%s）成功", orderID)
	logger.Logger.WXInfo("管理员（%s）为用户（%s）创建工单（%s）：%s", user.Uid, srcUser.Uid, order.Title, communicate.Content)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
