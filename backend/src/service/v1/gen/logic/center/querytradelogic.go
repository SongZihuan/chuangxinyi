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

type QueryTradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTradeLogic {
	return &QueryTradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTradeLogic) QueryTrade(req *types.QueryTradeReq) (resp *types.QueryTradeResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	pay, err := payModel.FindByPayID(l.ctx, req.TradeID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.QueryTradeResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TradeNotFound, "订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if user.WalletId != pay.WalletId {
		return &types.QueryTradeResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TradeNotFound, "订单并非此用户"),
		}, nil
	}

	var status int64
	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap {
		pay.TradeStatus, err = alipay.QueryTrade(l.ctx, user, pay)
		if err != nil {
			return &types.QueryTradeResp{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.QueryTradeFail, errors.WarpQuick(err), "支付宝查询失败"),
			}, nil
		}
	} else if pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		pay.TradeStatus, err = wechatpay.QueryTrade(l.ctx, user, pay)
		if err != nil {
			return &types.QueryTradeResp{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.QueryTradeFail, errors.WarpQuick(err), "微信支付查询失败"),
			}, nil
		}
	}

	status = pay.TradeStatus

	if status == db.PayWait || status == db.PayClose {
		return &types.QueryTradeResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryTradeData{
				Success: false,
			},
		}, nil
	} else if status == db.PaySuccess || status == db.PayFinish || status == db.PayCloseRefund {
		return &types.QueryTradeResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryTradeData{
				Success: true,
			},
		}, nil
	} else {
		return &types.QueryTradeResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryTradeData{
				Success: true,
			},
		}, nil
	}
}
