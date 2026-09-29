package before_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAlipayAppIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAlipayAppIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlipayAppIDLogic {
	return &GetAlipayAppIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAlipayAppIDLogic) GetAlipayAppID() (resp *types.AppIDResp, err error) {
	if len(config.BackendConfig.Alipay.AppID) == 0 {
		return &types.AppIDResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadAppID, "支付宝APPID未配置"),
		}, nil
	}

	return &types.AppIDResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AppIDData{
			AppID: config.BackendConfig.Alipay.AppID,
		},
	}, nil
}
