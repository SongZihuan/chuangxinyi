package policycheck

import "net/http"

func IsWebsitePolicyCheck(r *http.Request) bool {
	domainUID := r.Header.Get("X-Domain")
	if len(domainUID) == 0 {
		domainUID = r.URL.Query().Get("xdomain")
		if len(domainUID) == 0 {
			return false
		}
	}

	return true
}
