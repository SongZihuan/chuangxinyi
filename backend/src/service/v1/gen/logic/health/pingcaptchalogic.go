package health

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type PingCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingCaptchaLogic {
	return &PingCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingCaptchaLogic) PingCaptcha(r *http.Request) (resp *types.PingResp, err error) {
	return &types.PingResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: GetPingData(r),
	}, nil
}
