package before_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetWeChatAppIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWeChatAppIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWeChatAppIDLogic {
	return &GetWeChatAppIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWeChatAppIDLogic) GetWeChatAppID() (resp *types.AppIDResp, err error) {
	if len(config.BackendConfig.WeChat.AppID) == 0 {
		return &types.AppIDResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadAppID, "微信APPID未配置"),
		}, nil
	}

	return &types.AppIDResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AppIDData{
			AppID: config.BackendConfig.WeChat.AppID,
		},
	}, nil
}
