package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/cron"
)

func UrlPathCron() {
	tmp := cron.UrlPathCron(*AllPermission(), Roles(), PermissionList())

	urlPathLock.Lock()
	defer urlPathLock.Unlock()

	urlPathMap = tmp
}
