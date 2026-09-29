package application

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/application"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
)

func GetApplicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := application.NewGetApplicationLogic(r.Context(), svcCtx)
		resp, err := l.GetApplication()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
