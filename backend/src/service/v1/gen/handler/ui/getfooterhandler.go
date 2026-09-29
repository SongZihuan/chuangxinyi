package ui

import (
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/ui"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
)

func GetFooterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ui.NewGetFooterLogic(r.Context(), svcCtx)
		resp, err := l.GetFooter()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
