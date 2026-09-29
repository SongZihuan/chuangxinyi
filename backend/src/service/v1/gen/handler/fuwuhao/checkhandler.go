package fuwuhao

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/fuwuhao"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
)

func CheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := fuwuhao.NewCheckLogic(r.Context(), svcCtx)
		l.Check(w, r)
	}
}
