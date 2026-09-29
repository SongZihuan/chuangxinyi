package middleware

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/ip"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"net/http"
)

type DevelopMiddleware struct {
}

func NewDevelopMiddleware() *DevelopMiddleware {
	return &DevelopMiddleware{}
}

func (m *DevelopMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		geo, ok := r.Context().Value("X-Real-IP-Geo").(string)
		if config.BackendConfig.GetMode() == config.RunModeDevelop || (ok && geo == ip.LocalGeo) {
			next(w, r)
		} else {
			resp := respmsg.GetRespByMsg(r.Context(), respmsg.OnlyDevelop, "接口不支持非develop模式下的非内网调用")
			utils.NotFound(w, r, nil, false, resp.RequestsID)
		}
	}
}
