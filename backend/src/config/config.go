package config

import (
	"fmt"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
	"os"
	"strconv"
)

var BackendConfig Config
var BackendConfigViper *viper.Viper
var EnvPrefix string

func InitBackendConfigViper(config string, env string) errors.WTError {
	EnvPrefix = env
	BackendConfigViper = viper.New()

	BackendConfigViper.SetConfigType("yaml")
	BackendConfigViper.SetConfigName("config")
	BackendConfigViper.SetEnvPrefix(env)
	BackendConfigViper.AddConfigPath(config)
	setDefaultValue(BackendConfigViper)

	err := BackendConfigViper.ReadInConfig()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = BackendConfigViper.Unmarshal(&BackendConfig)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func setDefaultValue(v *viper.Viper) {
	v.SetDefault("mode", "develop")
	v.SetDefault("mysql.sqlFilePath", "sql")
	v.SetDefault("smtp.templateFilePath", "template")
	v.SetDefault("aliyun.ocrEndpoint", "ocr-api.cn-hangzhou.aliyuncs.com")
	v.SetDefault("aliyun.afs.captchaStatus", true)
	v.SetDefault("aliyun.afs.silenceCAPTCHAStatus", true)
	v.SetDefault("admin.rootRoleName", "root")
	v.SetDefault("admin.rootRoleDescribe", "create by system")
	v.SetDefault("admin.rootAdminName", "admin")
	v.SetDefault("admin.rootAdminPassword", "123456")
	v.SetDefault("alipay.sandbox", false)
	v.SetDefault("fuWuHao.menu.updateMenu", true)

	portString := os.Getenv(fmt.Sprintf("%sPORT", EnvPrefix))
	port, err := strconv.ParseInt(portString, 10, 64)
	if err != nil || port == 0 {
		port = 3351
	}

	devPortString := os.Getenv(fmt.Sprintf("%sDEV_PORT", EnvPrefix))
	devPort, err := strconv.ParseInt(devPortString, 10, 64)
	if err != nil || devPort == 0 {
		devPort = 4351
	}

	GoZeroHttpSetDefaultValue(v, "user", "user-server", port, devPort)
	JwtSetDefaultValue(v)
}
