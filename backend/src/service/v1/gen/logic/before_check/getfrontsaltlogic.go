package before_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetFrontSaltLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFrontSaltLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFrontSaltLogic {
	return &GetFrontSaltLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFrontSaltLogic) GetFrontSalt() (resp *types.SaltResp, err error) {
	return &types.SaltResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SaltData{
			Salt: config.BackendConfig.Password.FrontSalt,
		},
	}, nil
}
