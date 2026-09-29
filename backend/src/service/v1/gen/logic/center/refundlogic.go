package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/selfpay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
)

type RefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundLogic {
	return &RefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundLogic) Refund(req *types.RefundReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	pay, err := payModel.FindByPayID(l.ctx, req.TradeID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TradeNotFound, "订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if pay.WalletId != user.WalletId {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TradeNotFound, "订单并非此用户"),
		}, nil
	}

	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap {
		err = alipay.NewRefund(l.ctx, user, pay)
	} else if pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		err = wechatpay.NewRefund(l.ctx, user, pay)
	} else {
		err = selfpay.NewRefund(l.ctx, user, pay)
	}
	if errors.Is(err, alipay.Insufficient) || errors.Is(err, wechatpay.Insufficient) || errors.Is(err, selfpay.Insufficient) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Insufficient, "额度不足"),
		}, nil
	} else if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.RefundFail, errors.WarpQuick(err), "退款失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户发起充值退款成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
