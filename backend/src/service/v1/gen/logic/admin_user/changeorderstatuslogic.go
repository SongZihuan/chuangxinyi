package admin_user

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type ChangeOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeOrderStatusLogic {
	return &ChangeOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeOrderStatusLogic) ChangeOrderStatus(req *types.AdminChangeOrderStatusReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	if req.Status != db.WorkOrderStatusWaitReply && req.Status != db.WorkOrderStatusWaitUser && req.Status != db.WorkOrderStatusFinish {
		req.Status = db.WorkOrderStatusFinish
	}

	workOrderModel := db.NewWorkOrderModel(mysql.MySQLConn)
	order, err := workOrderModel.FindOneByUidWithoutDelete(l.ctx, req.OrderID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WorkOrderNotFound, "工单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if web.ID != warp.UserCenterWebsite && web.ID != order.FromId {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WorkOrderNotFound, "工单不属于该站点"),
		}, nil
	}

	if req.Status == db.WorkOrderStatusFinish {
		order.FinishAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
	} else {
		order.FinishAt = sql.NullTime{
			Valid: false,
		}
	}

	order.Status = req.Status
	err = workOrderModel.UpdateCh(l.ctx, order)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员修改工单（%s）状态（%d）成功", order.Uid, order.Status)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
