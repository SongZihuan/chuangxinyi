package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/cron"
)

func PermissionCron() {
	pm, pl, allP, userP, anonymousP := cron.PermissionCron(Roles())

	permissionLock.Lock()
	defer permissionLock.Unlock()

	permissionsSign = pm
	permissionList = pl
	allPermission = &allP
	userPermission = &userP
	anonymousPermission = &anonymousP
}
