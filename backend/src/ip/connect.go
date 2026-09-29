package ip

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	errors "github.com/wuntsong-org/wterrors"
)

func InitYunIP() errors.WTError {
	if len(config.BackendConfig.Aliyun.IP.AppCode) == 0 {
		return errors.Errorf("aliyun ip app code must be given")
	}

	if len(config.BackendConfig.Aliyun.IP.AppSecret) == 0 {
		return errors.Errorf("aliyun ip app secret must be given")
	}

	if len(config.BackendConfig.Aliyun.IP.AppKey) == 0 {
		return errors.Errorf("aliyun ip app key must be given")
	}
	return nil
}
