package middleware

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/policycheck"
	"net/http"
)

type IPCheckMiddleware struct {
}

func NewIPCheckMiddleware() *IPCheckMiddleware {
	return &IPCheckMiddleware{}
}

func (m *IPCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		policycheck.WebsitePolicyCheck(w, r, next, false)
	}
}
