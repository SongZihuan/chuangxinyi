package policycheck

import "net/http"

type Options struct {
}

func (Options) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	UserPolicyCheckOptions(w, r)
}
