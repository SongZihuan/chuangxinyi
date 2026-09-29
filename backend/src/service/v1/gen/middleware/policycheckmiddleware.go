package middleware

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/policycheck"
	"net/http"
)

type PolicyCheckMiddleware struct {
}

func NewPolicyCheckMiddleware() *PolicyCheckMiddleware {
	return &PolicyCheckMiddleware{}
}

func (m *PolicyCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		policycheck.UserPolicyCheck(w, r, next)
	}
}
