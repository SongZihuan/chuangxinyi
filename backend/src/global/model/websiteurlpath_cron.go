package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/cron"
)

func WebsiteUrlPathCron() {
	tmp := cron.WebsiteUrlPathCron(*WebsiteAllPermission(), Websites(), WebsitePermissionList())

	websiteUrlPathLock.Lock()
	defer websiteUrlPathLock.Unlock()

	websiteUrlPathMap = tmp
}
