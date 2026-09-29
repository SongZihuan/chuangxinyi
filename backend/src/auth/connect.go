package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

var PriKey *rsa.PrivateKey

func InitAuth() errors.WTError {
	if len(config.BackendConfig.Sign.DefrayPriKey) == 0 {
		return errors.Errorf("defray private key must be given")
	}

	priKeyString, err := base64.StdEncoding.DecodeString(config.BackendConfig.Sign.DefrayPriKey)
	if err != nil {
		return errors.WarpQuick(err)
	}

	PriKey, err = utils.ReadRsaPrivateKey(priKeyString)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
