package admin_role

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/admin_role"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
)

func DBUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin_role.NewDBUpdateLogic(r.Context(), svcCtx)
		resp, err := l.DBUpdate()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
