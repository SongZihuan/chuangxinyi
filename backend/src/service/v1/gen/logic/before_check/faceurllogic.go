package before_check

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"net/http"
	"net/url"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type FaceUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFaceUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FaceUrlLogic {
	return &FaceUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FaceUrlLogic) FaceUrl(req *types.FaceUrlReq, w http.ResponseWriter, r *http.Request) error {
	u, err := alipay.GetFaceUrl(l.ctx, req.CertifyID)
	if err != nil {
		return respmsg.AlipayError.WarpQuick(err)
	}

	p := &url.Values{}
	p.Set("appId", "20000067") // 固定appID
	p.Set("url", u)

	http.Redirect(w, r, fmt.Sprintf("%s?%s", "alipays://platformapi/startapp", p.Encode()), http.StatusFound)
	return nil
}
