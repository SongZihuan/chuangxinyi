package wechat

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	errors "github.com/wuntsong-org/wterrors"
)

func InitWeChat() errors.WTError {
	if len(config.BackendConfig.WeChat.AppID) == 0 {
		return errors.Errorf("wechat appid must be given")
	}

	if len(config.BackendConfig.WeChat.AppSecret) == 0 {
		return errors.Errorf("wechat app secret must be given")
	}

	return nil
}
