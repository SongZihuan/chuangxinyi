package before_check

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/before_check"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
)

func GetAlipayAppIDHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := before_check.NewGetAlipayAppIDLogic(r.Context(), svcCtx)
		resp, err := l.GetAlipayAppID()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
