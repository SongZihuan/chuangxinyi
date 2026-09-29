package cron

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

func AddMenuHandler() errors.WTError {
	MenuHandler() // 先执行一次

	id, err := Cron.AddFunc(config.BackendConfig.Cron.MenuUpdate, MenuHandler)
	if err != nil {
		return errors.WarpQuick(err)
	}

	logger.Logger.Info("menu update job id: %d", id)
	return nil
}

func MenuHandler() {
	defer utils.Recover(logger.Logger, nil, "")
	model.MenuCron()
}
