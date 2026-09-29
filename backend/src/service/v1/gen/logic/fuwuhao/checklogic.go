package fuwuhao

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"net/http"
	"sort"
	"strings"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckLogic {
	return &CheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckLogic) Check(w http.ResponseWriter, r *http.Request) {
	signature := r.URL.Query().Get("signature")
	timestamp := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	echostr := r.URL.Query().Get("echostr")

	tmpArr := []string{config.BackendConfig.FuWuHao.Token, timestamp, nonce}
	sort.Strings(tmpArr)
	tmpStr := utils.HashSHA1(strings.Join(tmpArr, ""))

	if tmpStr == signature {
		_, _ = fmt.Fprintf(w, echostr)
		w.WriteHeader(http.StatusOK)
	}
}
