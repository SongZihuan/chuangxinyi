package cron

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/dbinit"
	"github.com/robfig/cron/v3"
	errors "github.com/wuntsong-org/wterrors"
)

var Cron *cron.Cron

func InitCron() errors.WTError {
	if len(config.BackendConfig.Cron.MenuUpdate) == 0 {
		return errors.Errorf("cron menu update time must be given")
	}

	if len(config.BackendConfig.Cron.PermissionUpdate) == 0 {
		return errors.Errorf("cron permission update time must be given")
	}

	if len(config.BackendConfig.Cron.UrlPathUpdate) == 0 {
		return errors.Errorf("cron url path update time must be given")
	}

	if len(config.BackendConfig.Cron.RoleUpdate) == 0 {
		return errors.Errorf("cron role update time must be given")
	}

	if len(config.BackendConfig.Cron.WebsiteUpdate) == 0 {
		return errors.Errorf("cron website update time must be given")
	}

	if len(config.BackendConfig.Cron.WebsitePermissionUpdate) == 0 {
		return errors.Errorf("cron website permission update time must be given")
	}

	if len(config.BackendConfig.Cron.WebsiteUrlPathUpdate) == 0 {
		return errors.Errorf("cron website url path update time must be given")
	}

	Cron = cron.New()

	err := AddPermissionHandler()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = AddMenuHandler()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = AddRoleHandler()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = AddUrlPathHandler()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = dbinit.ResetRoleAdmin()
	if err != nil {
		return errors.WarpQuick(err)
	}

	//UrlPathHandler()  // 不需要再次刷新，因为UrlPathHandler本来就是最后一个

	MenuHandler() // 再次执行刷新Menu的role列表

	RoleHandler() // 再次刷新role的menu列表

	PermissionHandler() // 再次执行刷新permission列表

	err = AddWebsitePermissionHandler()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = AddWebsiteHandler()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = AddWebsiteUrlPathHandler()
	if err != nil {
		return errors.WarpQuick(err)
	}

	//WebsiteUrlPathHandler()  // 不需要再次刷新，因为是最后一个

	WebsiteHandler() // 再次刷新website列表

	WebsitePermissionHandler() // 再次执行刷新permission列表

	err = AddApplicationHandler() // 需要在website后面执行
	if err != nil {
		return errors.WarpQuick(err)
	}

	Cron.Start()

	return nil
}
