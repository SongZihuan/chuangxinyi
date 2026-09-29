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

type QueryWithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryWithdrawLogic {
	return &QueryWithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryWithdrawLogic) QueryWithdraw(req *types.QueryWithdrawReq) (resp *types.QueryWithdrawResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)

	withdraw, err := withdrawModel.FindByWithdrawID(l.ctx, req.TradeID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.QueryWithdrawResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WithdrawNotFound, "提现订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if user.WalletId != withdraw.WalletId {
		return &types.QueryWithdrawResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WithdrawNotFound, "提现订单并非此用户"),
		}, nil
	}

	var status int64
	if withdraw.WithdrawWay == alipay.WithdrawAlipay {
		withdraw.Status, err = alipay.QueryWithdraw(l.ctx, user, withdraw)
		if err != nil {
			return &types.QueryWithdrawResp{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.QueryTradeFail, errors.WarpQuick(err), "支付宝支付查询失败"),
			}, nil
		}
	} else if withdraw.WithdrawWay == wechatpay.WithdrawWechatpay {
		withdraw.Status, err = wechatpay.QueryWithdraw(l.ctx, user, withdraw)
		if err != nil {
			return &types.QueryWithdrawResp{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.QueryTradeFail, errors.WarpQuick(err), "微信支付查询失败"),
			}, nil
		}
	}

	status = withdraw.Status

	if status == db.WithdrawWait || status == db.WithdrawFail {
		return &types.QueryWithdrawResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryWithdrawData{
				Withdraw: false,
			},
		}, nil
	} else {
		return &types.QueryWithdrawResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.QueryWithdrawData{
				Withdraw: true,
			},
		}, nil
	}
}
