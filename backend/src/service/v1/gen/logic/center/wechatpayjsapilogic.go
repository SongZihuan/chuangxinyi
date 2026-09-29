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

type WechatPayJsAPILogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatPayJsAPILogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatPayJsAPILogic {
	return &WechatPayJsAPILogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatPayJsAPILogic) WechatPayJsAPI(req *types.WechatPayJsAPIReq) (resp *types.WechatPayJsAPIResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	prePayID, sign, id, timeExpire, err := wechatpay.NewPageJsAPI(l.ctx, user, fmt.Sprintf("%s：%.2f", config.BackendConfig.Coin.Name, float64(req.CNY)/100.00), req.CNY, req.CouponsID)
	switch true {
	case err == nil:
		// pass
	case errors.Is(err, wechatpay.BadCNY):
		return &types.WechatPayJsAPIResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "错误的金额"),
		}, nil
	default:
		return &types.WechatPayJsAPIResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateTradeFail, errors.WarpQuick(err), "微信支付失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户发起微信支付成功")

	return &types.WechatPayJsAPIResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.WechatPayJsAPIData{
			PrePayID:   prePayID,
			ID:         id,
			AppId:      sign.AppId,
			TimeStamp:  sign.TimeStamp,
			NonceStr:   sign.NonceStr,
			Package:    sign.Package,
			SignType:   sign.SignType,
			PaySign:    sign.PaySign,
			TimeExpire: timeExpire.Unix(),
		},
	}, nil
}
