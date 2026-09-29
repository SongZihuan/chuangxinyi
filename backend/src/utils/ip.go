package utils

import (
	"net"
	"net/http"
	"strings"
)

func CheckIPInList(ip string, ipList []string, canAll bool, canIP bool, canCIDR bool) bool {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return false
	}

	for _, i := range ipList {
		if i == "0.0.0.0" || i == "0.0.0.0/0" {
			if canAll {
				return true
			}
		} else if canIP && i == ip {
			return true
		} else if canCIDR {
			_, ipNet, err := net.ParseCIDR(i)
			if err == nil && ipNet.Contains(netIP) {
				return true
			}
		}
	}
	return false
}

func GetTargetIP(r *http.Request) (ip string) {
	ip = getTargetIPInXFF(r)
	if len(ip) != 0 {
		return ip
	}

	ip = getTargetIPInRealIP(r)
	if len(ip) != 0 {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return "0.0.0.0"
}

func getTargetIPInXFF(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		return ""
	}

	xffs := strings.Split(xff, ",")
	if len(xffs) == 0 {
		return ""
	}

	ip := net.ParseIP(strings.TrimSpace(xffs[0]))
	if ip != nil {
		return ip.String()
	}
	return ""
}

func getTargetIPInRealIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		return ""
	}

	ip := net.ParseIP(realIP)
	if ip != nil {
		return ip.String()
	}
	return ""
}

var loaclCIRD = []string{
	"127.0.0.0/8",
	"::1/128",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"fe80::/10",
	"fec0::/7",
}

func IsLocalIP(ip string) bool {
	for _, cidr := range loaclCIRD {
		if IsIPInCIDR(ip, cidr) {
			return true
		}
	}

	return false
}

func IsIPInCIDR(ip string, cidr string) bool {
	ipNet := net.ParseIP(ip)
	_, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	return cidrNet.Contains(ipNet)
}
