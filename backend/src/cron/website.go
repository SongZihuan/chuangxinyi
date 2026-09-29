package cron

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

func AddWebsiteHandler() errors.WTError {
	WebsiteHandler() // 先执行一次

	id, err := Cron.AddFunc(config.BackendConfig.Cron.WebsiteUpdate, WebsiteHandler)
	if err != nil {
		return errors.WarpQuick(err)
	}

	logger.Logger.Info("Website update job id: %d", id)
	return nil
}

func WebsiteHandler() {
	defer utils.Recover(logger.Logger, nil, "")
	model.WebsiteCron()
}
