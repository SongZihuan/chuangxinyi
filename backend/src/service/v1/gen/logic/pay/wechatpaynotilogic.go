package pay

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type WechatPayNotiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatPayNotiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatPayNotiLogic {
	return &WechatPayNotiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatPayNotiLogic) WechatPayNoti(w http.ResponseWriter, r *http.Request) {
	wechatpay.Notification(w, r)
}
