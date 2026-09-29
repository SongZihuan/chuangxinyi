package before_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"net/http"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type FaceInternalUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFaceInternalUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FaceInternalUrlLogic {
	return &FaceInternalUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FaceInternalUrlLogic) FaceInternalUrl(req *types.FaceUrlReq, w http.ResponseWriter, r *http.Request) error {
	u, err := alipay.GetFaceUrl(l.ctx, req.CertifyID)
	if err != nil {
		return respmsg.AlipayError.WarpQuick(err)
	}

	http.Redirect(w, r, u, http.StatusFound)
	return nil
}
