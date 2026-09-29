package check

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/check"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
)

func CheckUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckUser
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := check.NewCheckUserLogic(r.Context(), svcCtx)
		resp, err := l.CheckUser(&req, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
