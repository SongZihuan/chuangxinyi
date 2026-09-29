package middleware

import (
	"context"
	"net/http"
)

type WebSocketMiddleware struct {
}

func NewWebSocketMiddleware() *WebSocketMiddleware {
	return &WebSocketMiddleware{}
}

func (m *WebSocketMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 不用检测是否Upgrade请求
		next(w, r.WithContext(context.WithValue(r.Context(), "X-Websocket", "True")))
	}
}
