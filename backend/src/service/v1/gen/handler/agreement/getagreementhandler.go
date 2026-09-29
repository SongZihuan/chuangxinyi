package agreement

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/agreement"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
)

func GetAgreementHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAgreementReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := agreement.NewGetAgreementLogic(r.Context(), svcCtx)
		err := l.GetAgreement(&req, w, r)
		if err != nil {
			resp := respmsg.GetRespByError(r.Context(), respmsg.NotFound, errors.WarpQuick(err), "获取协议错误")
			utils.NotFound(w, r, err, config.BackendConfig.GetModeFromRequests(r) == config.RunModeDevelop, resp.RequestsID)
		}
	}
}
