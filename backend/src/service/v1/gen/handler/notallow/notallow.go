package notallow

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/record"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	"net/http"
)

type NotAllow struct{}

func (NotAllow) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "X-Record", &record.Record{RequestsID: "notallow-403"})
	httpx.WriteJsonCtx(r.Context(), w, http.StatusMethodNotAllowed, &types.RespEmpty{
		Resp: respmsg.GetRespByMsg(ctx, respmsg.MethodNotAllow, "路由方法不允许"),
	})
}
