package admin_user

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/admin_user"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
)

func UploadDefaultHeaderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadDefaultHeader
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin_user.NewUploadDefaultHeaderLogic(r.Context(), svcCtx)
		resp, err := l.UploadDefaultHeader(&req, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
