package svc

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/middleware"
	"github.com/wuntsong-org/go-zero-plus/rest"
)

type ServiceContext struct {
	Config      config.UserConfig
	Develop     rest.Middleware
	PolicyCheck rest.Middleware
	IPCheck     rest.Middleware
	WebSocket   rest.Middleware
}

func NewServiceContext(c config.UserConfig) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		Develop:     middleware.NewDevelopMiddleware().Handle,
		PolicyCheck: middleware.NewPolicyCheckMiddleware().Handle,
		IPCheck:     middleware.NewIPCheckMiddleware().Handle,
		WebSocket:   middleware.NewWebSocketMiddleware().Handle,
	}
}
