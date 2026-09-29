package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type QueryRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRefundLogic {
	return &QueryRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRefundLogic) QueryRefund(req *types.QueryRefundReq) (resp *types.QueryRefundResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	pay, err := payModel.FindByPayID(l.ctx, req.TradeID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.QueryRefundResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TradeNotFound, "订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if user.WalletId != pay.WalletId {
		return &types.QueryRefundResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TradeNotFound, "订单并非此用户"),
		}, nil
	}

	var status int64
	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap {
		pay.TradeStatus, err = alipay.QueryRefund(l.ctx, pay)
		if err != nil {
			return &types.QueryRefundResp{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.QueryTradeFail, errors.WarpQuick(err), "支付宝查询失败"),
			}, nil
		}
	} else if pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		pay.TradeStatus, err = wechatpay.QueryRefund(l.ctx, user, pay)
		if err != nil {
			return &types.QueryRefundResp{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.QueryTradeFail, errors.WarpQuick(err), "微信支付查询失败"),
			}, nil
		}
	}

	status = pay.TradeStatus

	if status == db.PayWait || status == db.PayClose {
		return &types.QueryRefundResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryRefundData{
				Refund: false,
			},
		}, nil
	} else if status == db.PaySuccess || status == db.PayFinish || status == db.PayCloseRefund {
		return &types.QueryRefundResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryRefundData{
				Refund: false,
			},
		}, nil
	} else if status == db.PayWaitRefund {
		return &types.QueryRefundResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryRefundData{
				Refund: false,
			},
		}, nil
	} else {
		return &types.QueryRefundResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryRefundData{
				Refund: true,
			},
		}, nil
	}
}
