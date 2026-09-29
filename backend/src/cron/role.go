package cron

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

func AddRoleHandler() errors.WTError {
	RoleHandler() // 先执行一次

	id, err := Cron.AddFunc(config.BackendConfig.Cron.RoleUpdate, RoleHandler)
	if err != nil {
		return errors.WarpQuick(err)
	}

	logger.Logger.Info("role update job id: %d", id)
	return nil
}

func RoleHandler() {
	defer utils.Recover(logger.Logger, nil, "")
	model.RoleCron()
}
