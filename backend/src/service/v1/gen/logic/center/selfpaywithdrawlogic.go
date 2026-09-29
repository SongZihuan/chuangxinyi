package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/selfpay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type SelfpayWithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelfpayWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelfpayWithdrawLogic {
	return &SelfpayWithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelfpayWithdrawLogic) SelfpayWithdraw(req *types.SelfpayWithdrawReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.Cny < config.BackendConfig.Coin.WithdrawMin {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "提现金额错误"),
		}, nil
	}

	_, err = selfpay.NewWithdraw(l.ctx, req.Cny, user, req.WithdrawWay, req.Name)
	switch true {
	case err == nil:
		// pass
	case errors.Is(err, selfpay.BadName):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadName, "错误的提现人"),
		}, nil
	case errors.Is(err, selfpay.BadPayWay):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayWay, "提现方式错误"),
		}, nil
	case errors.Is(err, selfpay.BadCNY):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "金额错误"),
		}, nil
	case errors.Is(err, selfpay.Insufficient):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Insufficient, "额度不足"),
		}, nil
	default:
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateTradeFail, errors.WarpQuick(err), "人工提现失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户申请线下人工提现")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
