package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AlipayWithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlipayWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlipayWithdrawLogic {
	return &AlipayWithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlipayWithdrawLogic) AlipayWithdraw(req *types.AlipayWithdrawReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.Cny < config.BackendConfig.Coin.WithdrawMin {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "提现金额错误"),
		}, nil
	}

	_, err = alipay.NewWithdraw(l.ctx, user, req.Cny, req.Identity, req.Name)
	switch true {
	case err == nil:
	case errors.Is(err, alipay.BadName):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadName, "提现名称错误"),
		}, nil
	case errors.Is(err, alipay.Insufficient):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Insufficient, "提现额度不足"),
		}, nil
	case errors.Is(err, alipay.BadCNY):
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "提现金额错误"),
		}, nil
	default:
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateTradeFail, errors.WarpQuick(err), "提现失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户发起支付宝提现")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
