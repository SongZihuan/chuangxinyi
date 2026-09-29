package center

import (
	"context"
	"fmt"
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

type AlipayWapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlipayWapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlipayWapLogic {
	return &AlipayWapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlipayWapLogic) AlipayWap(req *types.AlipayWapReq) (resp *types.AlipayWapResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	url, id, timeExpire, err := alipay.NewPageWap(l.ctx, user, fmt.Sprintf("%s：%.2f", config.BackendConfig.Coin.Name, float64(req.CNY)/100.00), req.CNY, req.CouponsID, req.ReturnURL, req.QuiteUrl)
	switch true {
	case err == nil:
		// pass
	case errors.Is(err, alipay.BadCNY):
		return &types.AlipayWapResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPayCNY, "支付金额错误"),
		}, nil
	default:
		return &types.AlipayWapResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.CreateTradeFail, errors.WarpQuick(err), "支付失败"),
		}, nil
	}

	audit.NewUserAudit(user.Id, "用户发起支付宝支付")

	return &types.AlipayWapResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AlipayWapData{
			PayUrl:     url,
			ID:         id,
			TimeExpire: timeExpire.Unix(),
		},
	}, nil
}
