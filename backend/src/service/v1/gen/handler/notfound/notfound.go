package notfound

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/record"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	"net/http"
)

type NotFound struct{}

func (NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "X-Record", &record.Record{RequestsID: "notfound-404"})
	httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, &types.RespEmpty{
		Resp: respmsg.GetRespByMsg(ctx, respmsg.NotFound, "路由未找到"),
	})
}
