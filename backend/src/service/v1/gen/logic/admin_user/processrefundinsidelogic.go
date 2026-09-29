package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/selfpay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type ProcessRefundInsideLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessRefundInsideLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessRefundInsideLogic {
	return &ProcessRefundInsideLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessRefundInsideLogic) ProcessRefundInside(req *types.AdminProcessRefundReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	payModel := db.NewPayModel(mysql.MySQLConn)
	pay, err := payModel.FindByPayID(l.ctx, req.TradeID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TradeNotFound, "订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	srcUser, err := userModel.FindOneByIDWithoutDelete(l.ctx, pay.UserId)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if srcUser.WalletId != pay.WalletId {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TradeNotFound, "订单和用户的钱包ID不匹配"),
		}, nil
	}

	if req.Success {
		err = selfpay.RefundInside(l.ctx, srcUser, pay)
	} else if pay.TradeStatus == db.PayWait || pay.TradeStatus == db.PayClose || pay.TradeStatus == db.PaySuccessRefund || pay.TradeStatus == db.PaySuccessRefundInside {
		err = nil // 退款不成功
	} else {
		err = selfpay.RefundFail(l.ctx, pay)
	}
	if errors.Is(err, selfpay.Insufficient) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Insufficient, "额度不足"),
		}, nil
	} else if errors.Is(err, selfpay.BadStatus) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.RefundFail, "错误的支付状态"),
		}, nil
	} else if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.RefundFail, errors.WarpQuick(err), "退款失败"),
		}, nil
	}

	audit.NewAdminAudit(user.Id, "管理员处理单边退款（%s）", pay.PayId)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
