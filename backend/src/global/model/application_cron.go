package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/cron"
)

func ApplicationCron() {
	name, lst := cron.ApplicationCron(*WebsiteAllPermission(), Websites())

	applicationLock.Lock()
	defer applicationLock.Unlock()

	applicationName = name
	applicationList = lst
}
