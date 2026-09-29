package pay

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/pay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
)

func WechatPayNotiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := pay.NewWechatPayNotiLogic(r.Context(), svcCtx)
		l.WechatPayNoti(w, r)
	}
}
