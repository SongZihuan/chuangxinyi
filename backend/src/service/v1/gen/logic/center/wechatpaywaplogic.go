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

type WechatPayWapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatPayWapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatPayWapLogic {
	return &WechatPayWapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatPayWapLogic) WechatPayWap(req *types.WechatPayWapReq) (resp *types.WechatPayWapResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	ip, ok := l.ctx.Value("X-Real-IP").(string)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Real-IP")
	}

	url, id, timeExpire, err := wechatpay.NewPageH5(l.ctx, user, fmt.Sprintf("%s：%.2f", config.BackendConfig.Coin.Name, float64(req.CNY)/100.00), req.CNY, req.CouponsID, req.H5Type, ip)
	switch true {
	case err == nil:
		// pass
	case errors.Is(err, wechatpay.BadCNY):
		return &types.WechatPayWapResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "错误的金额"),
		}, nil
	default:
		return &types.WechatPayWapResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateTradeFail, errors.WarpQuick(err), "微信支付失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户发起微信支付成功")

	return &types.WechatPayWapResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.WechatPayWapData{
			H5Url:      url,
			ID:         id,
			TimeExpire: timeExpire.Unix(),
		},
	}, nil
}
