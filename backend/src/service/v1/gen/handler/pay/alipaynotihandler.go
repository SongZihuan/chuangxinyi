package pay

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/pay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
)

func AlipayNotiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := pay.NewAlipayNotiLogic(r.Context(), svcCtx)
		l.AlipayNoti(w, r)
	}
}
