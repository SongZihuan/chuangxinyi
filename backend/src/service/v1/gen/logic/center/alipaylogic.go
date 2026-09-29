package center

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
)

type AlipayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlipayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlipayLogic {
	return &AlipayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlipayLogic) Alipay(req *types.AlipayReq) (resp *types.AlipayResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	url, id, timeExpire, err := alipay.NewPagePay(l.ctx, user, fmt.Sprintf("%s：%.2f", config.BackendConfig.Coin.Name, float64(req.CNY)/100.00), req.CNY, req.PayMode, req.CouponsID, req.ReturnURL)
	switch true {
	case err == nil:
	case errors.Is(err, alipay.BadCNY):
		return &types.AlipayResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "支付金额错误"),
		}, nil
	default:
		return &types.AlipayResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateTradeFail, errors.WarpQuick(err), "支付失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户发起支付宝支付")

	return &types.AlipayResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AlipayData{
			Url:        url,
			ID:         id,
			TimeExpire: timeExpire.Unix(),
		},
	}, nil
}
