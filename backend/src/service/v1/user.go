package user

import (
	"context"
	"flag"
	"gitee.com/wuntsong-auth/backend/src/config"
	initall "gitee.com/wuntsong-auth/backend/src/init"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/peers"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/handler"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/handler/notallow"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/handler/notfound"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/accessrecord"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/checker"
	config2 "gitee.com/wuntsong-auth/backend/src/service/v1/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/policycheck"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/signalexit"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/wuntsong-org/go-zero-plus/rest"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	"net/http"
	"os"
)

var configFile = flag.String("f", "etc", "the config path")

func GlobalMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessrecord.AccessRecordHandle(w, r, next.ServeHTTP, []string{"/api/v1/ping"})
	})
}

func CmdMain() {
	flag.Parse()
	var err error

	respmsg.InitErrorHandler()
	restConfig, err := config2.InitConfig(*configFile)
	utils.MustNotError(err)

	err = initall.InitUserCenter()
	utils.MustNotError(err)

	v, err := checker.GetValidator()
	utils.MustNotError(err)

	httpx.SetValidator(v)

	server, err := rest.NewServer(restConfig,
		rest.WithNotFoundHandler(notfound.NotFound{}),
		rest.WithNotAllowedHandler(notallow.NotAllow{}),
		rest.WithOptionsHandler(policycheck.Options{}),
		rest.WithGlobalMiddleware(GlobalMiddleWare),
	)
	utils.MustNotError(err)
	defer server.Stop()

	ctx := svc.NewServiceContext(config.BackendConfig.User)
	handler.RegisterHandlers(server, ctx)

	logger.Logger.WXInfo("启动服务 %s 在端口 %s:%d...", config.BackendConfig.User.ReadableName, restConfig.Host, restConfig.Port)
	srv := server.StartAsGoRoutine()

	signalexit.AddExitByFunc(func(ctx context.Context, signal os.Signal) context.Context {
		_ = srv.Shutdown(ctx)
		return context.WithValue(ctx, "Server-Shutdown", true)
	})

	err = peers.ConnectPeers()
	if err != nil {
		logger.Logger.Error("连接节点失败 %s", err.Error())
		return
	}

	peers.SendPeersPing()

	select {} // 阻塞
}
