package fuwuhao

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/fuwuhao"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
)

func NotiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := fuwuhao.NewNotiLogic(r.Context(), svcCtx)
		l.Noti(w, r)
	}
}
