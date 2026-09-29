package utils

import (
	"encoding/base64"
	errors "github.com/wuntsong-org/wterrors"
	"strings"
)

func DecodeReqBase64(req string, allowEmpty bool) ([]byte, errors.WTError) {
	if len(req) == 0 {
		if allowEmpty {
			return []byte{}, nil
		}

		return nil, errors.Errorf("req is empty")
	}

	if !strings.HasPrefix(req, "base64:") {
		return nil, errors.Errorf("not prefix")
	}

	data, err := base64.StdEncoding.DecodeString(req[len("base64:"):])
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return data, nil
}
