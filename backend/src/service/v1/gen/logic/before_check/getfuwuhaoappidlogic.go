package before_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetFuwuhaoAppIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFuwuhaoAppIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFuwuhaoAppIDLogic {
	return &GetFuwuhaoAppIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFuwuhaoAppIDLogic) GetFuwuhaoAppID() (resp *types.AppIDResp, err error) {
	if len(config.BackendConfig.FuWuHao.AppID) == 0 {
		return &types.AppIDResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadAppID, "服务号APPID未配置"),
		}, nil
	}

	return &types.AppIDResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AppIDData{
			AppID: config.BackendConfig.FuWuHao.AppID,
		},
	}, nil
}
