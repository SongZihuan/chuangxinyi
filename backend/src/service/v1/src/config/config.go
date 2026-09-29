package config

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"github.com/wuntsong-org/go-zero-plus/rest"
	errors "github.com/wuntsong-org/wterrors"
)

func InitConfig(configPath string) (rest.RestConf, errors.WTError) {
	var err error
	err = config.InitBackendConfigViper(configPath, "AUTH_")
	if err != nil {
		return rest.RestConf{}, errors.WarpQuick(err)
	}

	return config.BackendConfig.User.GetRestConfig(), nil
}
