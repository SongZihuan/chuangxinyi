package center

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/center"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
)

func WSGetInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := center.NewWSGetInfoLogic(r.Context(), svcCtx)
		l.WSGetInfo(w, r)
	}
}
