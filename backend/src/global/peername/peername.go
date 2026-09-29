package peername

import (
	"fmt"
	errors "github.com/wuntsong-org/wterrors"

	"net"
	"os"
	"strings"
)

const Single = "single"

var PeerName = Single
var SelfIP []string
var Eth0IP string

func InitPeerName(envPrefix string) errors.WTError {
	SelfIP = GetSelfIP()

	if len(Eth0IP) == 0 {
		return errors.Errorf("can not get eth0")
	}

	PeerName = GetPeerName(envPrefix)

	return nil
}

func GetPeerName(envPrefix string) string {
	peerName := os.Getenv(fmt.Sprintf("%sPEER_NAME", envPrefix))
	if len(peerName) == 0 {
		hostname, err := os.Hostname()
		if err == nil && len(hostname) != 0 {
			peerName = hostname
		} else {
			hostnameEnv := os.Getenv("HOSTNAME")
			if len(hostnameEnv) != 0 {
				peerName = hostnameEnv
			} else {
				if len(SelfIP) != 0 {
					peerName = fmt.Sprintf("%s", SelfIP[0])
				} else {
					peerName = Single
				}
			}
		}
	}

	return peerName
}

func GetSelfIP() []string {
	interFace, err := net.Interfaces()
	if err != nil {
		return []string{}
	}

	ipList := make([]string, 0, 10)
	for _, face := range interFace {
		if !strings.HasPrefix(face.Name, "eth") && !strings.HasPrefix(face.Name, "以太网") && !strings.HasPrefix(face.Name, "WLAN") {
			continue
		}

		address, err := face.Addrs()
		if err != nil {
			continue
		}

		for _, address := range address {
			ipnet, ok := address.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}

			if face.Name == "eth0" || face.Name == "以太网" {
				if len(Eth0IP) == 0 || strings.Contains(Eth0IP, ":") && !strings.Contains(ipnet.IP.String(), ":") { // ipv4优先
					Eth0IP = ipnet.IP.String()
				}
			}
			ipList = append(ipList, ipnet.IP.String())
		}
	}

	if len(Eth0IP) == 0 && len(ipList) != 0 {
		Eth0IP = ipList[0]
	}

	return ipList
}
