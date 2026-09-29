package pay

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AlipayWangguanNotiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlipayWangguanNotiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlipayWangguanNotiLogic {
	return &AlipayWangguanNotiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlipayWangguanNotiLogic) AlipayWangguanNoti(w http.ResponseWriter, r *http.Request) {
	alipay.NotificationWangguan(w, r)
}
