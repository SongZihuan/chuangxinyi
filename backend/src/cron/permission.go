package cron

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

func AddPermissionHandler() errors.WTError {
	PermissionHandler() // 先执行一次

	id, err := Cron.AddFunc(config.BackendConfig.Cron.PermissionUpdate, PermissionHandler)
	if err != nil {
		return errors.WarpQuick(err)
	}

	logger.Logger.Info("permission update job id: %d", id)
	return nil
}

func PermissionHandler() {
	defer utils.Recover(logger.Logger, nil, "")
	model.PermissionCron()
}
