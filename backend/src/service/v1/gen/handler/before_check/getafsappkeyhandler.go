package before_check

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/before_check"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
)

func GetAfsAppKeyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := before_check.NewGetAfsAppKeyLogic(r.Context(), svcCtx)
		resp, err := l.GetAfsAppKey()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
