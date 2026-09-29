package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/cron"
)

func WebsitePermissionCron() {
	pm, pl, allP := cron.WebsitePermissionCron(Websites())

	websitePermissionLock.Lock()
	defer websitePermissionLock.Unlock()

	websitePermissionsSign = pm
	websitePermissionList = pl
	websiteAllPermission = &allP
}
