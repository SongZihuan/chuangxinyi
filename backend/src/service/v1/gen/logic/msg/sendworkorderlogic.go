package msg

import (
	"context"
	"database/sql"
	"encoding/base64"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/oss"
	"gitee.com/wuntsong-auth/backend/src/redis"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"github.com/google/uuid"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type SendWorkOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendWorkOrderLogic {
	return &SendWorkOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendWorkOrderLogic) SendWorkOrder(req *types.SendWorkOrder, r *http.Request) (resp *types.SendMsgResp, err error) {
	user, err := utils2.FindUser(l.ctx, req.UserID, false)
	if errors.Is(err, utils2.UserNotFound) {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
	allow, err := bannedModel.CheckAllow(r.Context(), user.Id, web.ID, db.AllowMsg)
	if err != nil || !allow {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "用户关闭了通信授权许可"),
		}, nil
	}

	workOrderModel := db.NewWorkOrderModel(mysql.MySQLConn)
	workOrderCommunicateModel := db.NewWorkOrderCommunicateModel(mysql.MySQLConn)

	orderIDUUID, success := redis.GenerateUUIDMore(l.ctx, "order", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		orderModel := db.NewWorkOrderModel(mysql.MySQLConn)
		_, err := orderModel.FindOneByUidWithoutDelete(ctx, u.String())
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "创建工单失败"),
			Data: types.SendMsgData{
				Success: false,
				Have:    true,
			},
		}, nil
	}

	orderID := orderIDUUID.String()

	_, err = workOrderModel.FindOneByUidWithoutDelete(l.ctx, orderID)
	if !errors.Is(err, db.ErrNotFound) {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.SendMsgData{
				Success: false,
				Have:    true,
			},
		}, nil
	}

	order := &db.WorkOrder{
		Uid:    orderID,
		UserId: user.Id,
		Title:  req.Title,
		From:   web.Name,
		FromId: web.ID,
		Status: db.WorkOrderStatusWaitReply,
		LastReplyAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}
	res1, err := workOrderModel.Insert(l.ctx, order)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}
	order.Id, err = res1.LastInsertId()
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	filename := make([]string, 0, len(req.File))
	for _, f := range req.File {
		filename = append(filename, f.FileName)
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

	go func(order *db.WorkOrder, communicate *db.WorkOrderCommunicate, file []types.SendWorkOrderFile) {
		for _, f := range file {
			fileByte, err := base64.StdEncoding.DecodeString(f.File)
			if err != nil {
				continue
			}

			err = oss.UploadWorkOrderFile(order, communicate, f.FileName, fileByte, true)
			if err != nil {
				logger.Logger.Error("mysql error: %s", err.Error())
				continue
			}
		}

		// 不需要更新ws事件，因为这是新建的，没人监听
	}(order, communicate, req.File)

	return &types.SendMsgResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SendMsgData{
			Success: true,
			Have:    true,
		},
	}, nil
}
