package before_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAfsAppKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAfsAppKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAfsAppKeyLogic {
	return &GetAfsAppKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAfsAppKeyLogic) GetAfsAppKey() (resp *types.AFSResp, err error) {
	return &types.AFSResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AFSData{
			HAppKey: config.BackendConfig.Aliyun.AFS.CAPTCHAAppKey,
			HScene:  config.BackendConfig.Aliyun.AFS.CAPTCHAScene,
			SAppKey: config.BackendConfig.Aliyun.AFS.SilenceCAPTCHAAppKey,
			SScene:  config.BackendConfig.Aliyun.AFS.SilenceCAPTCHAScene,
		},
	}, nil
}
