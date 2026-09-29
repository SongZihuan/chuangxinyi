package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type WechatpayWithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatpayWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatpayWithdrawLogic {
	return &WechatpayWithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatpayWithdrawLogic) WechatpayWithdraw(req *types.WechatpayWithdrawReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.Cny < config.BackendConfig.Coin.WithdrawMin {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "提现金额错误"),
		}, nil
	}

	_, err = wechatpay.NewWithdraw(l.ctx, user, req.Cny, req.Name)
	switch true {
	case err == nil:
		// pass
	case errors.Is(err, wechatpay.BadName):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadName, "错误的提现人"),
		}, nil
	case errors.Is(err, wechatpay.BadCNY):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "错误的金额"),
		}, nil
	case errors.Is(err, wechatpay.Insufficient):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Insufficient, "额度不足"),
		}, nil
	default:
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateTradeFail, errors.WarpQuick(err), "微信支付提现失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户发起微信提现成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
