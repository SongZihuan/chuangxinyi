package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type ReplyOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplyOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyOrderLogic {
	return &ReplyOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplyOrderLogic) ReplyOrder(req *types.AdminReplyOrder, r *http.Request) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	filename, ok := r.MultipartForm.Value["filename[]"]
	if !ok {
		filename = []string{}
	}

	if len(filename) > int(config.BackendConfig.MySQL.WorkOrderFileLimit) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooManyFile, "上传太多文件"),
		}, nil
	}

	if utils.HasDuplicate(filename) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.RepeatFileName, "上传文件有重复"),
		}, nil
	}

	file, ok := r.MultipartForm.File["file[]"]
	if !ok {
		file = []*multipart.FileHeader{}
	}

	if len(file) != len(filename) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadForm, "文件名数量和文件数量不匹配"),
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

	order, err := workOrderModel.FindOneByUidWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WorkOrderNotFound, "工单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if web.ID != warp.UserCenterWebsite && web.ID != order.FromId {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WorkOrderNotFound, "工单未找到"),
		}, nil
	}

	if order.Status != db.WorkOrderStatusFinish {
		order.Status = req.Status
		err = workOrderModel.UpdateCh(l.ctx, order)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
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

	communicate, err = workOrderCommunicateModel.FindOneWithoutDelete(l.ctx, communicate.Id) // 重新赋值，查找create_at
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	go func(order *db.WorkOrder, communicate *db.WorkOrderCommunicate, filename []string, fileByte [][]byte) {
		for i, b := range fileByte {
			err := oss.UploadWorkOrderFile(order, communicate, filename[i], b, true)
			if err != nil {
				continue
			}
		}

		db.NewWorkOrderCommunicate(communicate, order.Uid, order.FromId, mysql.MySQLConn)
	}(order, communicate, filename, fileByte)

	audit.NewAdminAudit(user.Id, "管理员回复（%d）工单（%s）成功", communicate.Id, order.Uid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
