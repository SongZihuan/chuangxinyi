package center

import (
	"context"
	"fmt"
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

type WechatPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatPayLogic {
	return &WechatPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatPayLogic) WechatPay(req *types.WechatPayReq) (resp *types.WechatPayResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	url, id, timeExpire, err := wechatpay.NewPagePay(l.ctx, user, fmt.Sprintf("%s：%.2f", config.BackendConfig.Coin.Name, float64(req.CNY)/100.00), req.CNY, req.CouponsID)
	switch true {
	case err == nil:
		// pass
	case errors.Is(err, wechatpay.BadCNY):
		return &types.WechatPayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "错误的金额"),
		}, nil
	default:
		return &types.WechatPayResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateTradeFail, errors.WarpQuick(err), "微信支付失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户发起微信支付成功")

	return &types.WechatPayResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.WechatPayData{
			Url:        url,
			ID:         id,
			TimeExpire: timeExpire.Unix(),
		},
	}, nil
}
