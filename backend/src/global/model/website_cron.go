package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/cron"
)

func WebsiteCron() {
	byID, byUID, byLst := cron.WebsiteCron(*WebsiteAllPermission(), WebsiteUrlPathMap(), WebsitePermissionList())

	websiteLock.Lock()
	defer websiteLock.Unlock()

	websiteList = byLst
	websites = byID
	websitesUID = byUID
}
