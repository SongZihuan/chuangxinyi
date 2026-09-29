package pay

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AlipayNotiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlipayNotiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlipayNotiLogic {
	return &AlipayNotiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlipayNotiLogic) AlipayNoti(w http.ResponseWriter, r *http.Request) {
	alipay.Notification(w, r)
}
