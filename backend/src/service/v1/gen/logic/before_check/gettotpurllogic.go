package before_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetTotpUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTotpUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTotpUrlLogic {
	return &GetTotpUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTotpUrlLogic) GetTotpUrl(req *types.GetTotpUrlReq) (resp *types.GetTotpUrlResp, err error) {
	if !utils.IsPhoneNumber(req.Phone) {
		return &types.GetTotpUrlResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPhone, "错误的手机号"),
		}, nil
	}

	secret := utils.GetRandomSecret()

	url := utils.GenerateTotpURL(secret, req.Phone, config.BackendConfig.Totp.IssuerName)
	return &types.GetTotpUrlResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetTotpUrlData{
			Url:    url,
			Secret: secret,
		},
	}, nil
}
