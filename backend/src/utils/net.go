package utils

import (
	"github.com/wuntsong-org/wterrors"
	"net"
)

func IsNetClose(err error) bool {
	if err == nil {
		return false
	}

	var netErr *net.OpError
	if !errors.As(err, &netErr) {
		return false
	}

	return errors.Is(netErr, net.ErrClosed)
}
