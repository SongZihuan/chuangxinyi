package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/cron"
)

func MenuCron() {
	tmp := cron.MenuCron(*AllPermission(), Roles(), PermissionList())

	menuLock.Lock()
	defer menuLock.Unlock()

	menus = tmp
}
