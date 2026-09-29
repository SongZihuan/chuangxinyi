package cron

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

func AddWebsiteUrlPathHandler() errors.WTError {
	WebsiteUrlPathHandler() // 先执行一次

	id, err := Cron.AddFunc(config.BackendConfig.Cron.WebsiteUrlPathUpdate, WebsiteUrlPathHandler)
	if err != nil {
		return errors.WarpQuick(err)
	}

	logger.Logger.Info("website url path update job id: %d", id)
	return nil
}

func WebsiteUrlPathHandler() {
	defer utils.Recover(logger.Logger, nil, "")
	model.WebsiteUrlPathCron()
}
