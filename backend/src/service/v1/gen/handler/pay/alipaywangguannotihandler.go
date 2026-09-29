package pay

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/pay"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
)

func AlipayWangguanNotiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := pay.NewAlipayWangguanNotiLogic(r.Context(), svcCtx)
		l.AlipayWangguanNoti(w, r)
	}
}
