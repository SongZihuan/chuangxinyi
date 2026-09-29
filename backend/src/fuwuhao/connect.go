package fuwuhao

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"github.com/fastwego/offiaccount"
	errors "github.com/wuntsong-org/wterrors"
)

var OffiAccount *offiaccount.OffiAccount

func InitFuWuHao() errors.WTError {
	if len(config.BackendConfig.FuWuHao.AppID) == 0 {
		return errors.Errorf("fuwuhao appid must be given")
	}

	if len(config.BackendConfig.FuWuHao.Secret) == 0 {
		return errors.Errorf("fuwuhao secret must be given")
	}

	if len(config.BackendConfig.FuWuHao.Token) == 0 {
		return errors.Errorf("fuwuhao token must be given")
	}

	if len(config.BackendConfig.FuWuHao.EncodingAESKey) == 0 {
		return errors.Errorf("fuwuhao encoding aes key must be given")
	}

	OffiAccount = offiaccount.New(offiaccount.Config{
		Appid:          config.BackendConfig.FuWuHao.AppID,
		Secret:         config.BackendConfig.FuWuHao.Secret,
		Token:          config.BackendConfig.FuWuHao.Token,
		EncodingAESKey: config.BackendConfig.FuWuHao.EncodingAESKey,
	})

	OffiAccount.AccessToken.GetAccessTokenHandler = func(ctx *offiaccount.OffiAccount) (string, error) {
		accessToken, _, err := getAccessToken(context.Background(), ctx.Config.Appid, ctx.Config.Secret)
		if err != nil {
			logger.Logger.Error("can not get access token: %s", err.Error())
			return "", err
		}
		return accessToken, nil
	}

	OffiAccount.AccessToken.NoticeAccessTokenExpireHandler = func(ctx *offiaccount.OffiAccount) error {
		err := delAccessToken(context.Background(), ctx.Config.Appid)
		if err != nil {
			logger.Logger.Error("delete access token: %s", err.Error())
			return err
		}
		return nil
	}

	go func() {
		if config.BackendConfig.FuWuHao.Menu.UpdateMenu {
			_ = CreateMenu()
		}
	}()

	return nil
}
