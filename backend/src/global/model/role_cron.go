package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/cron"
)

func RoleCron() {
	roleList, rolesSignList := cron.RoleCron(*AllPermission(), *WebsiteAllPermission(), Websites(), Menus(), UrlPathMap(), PermissionList())

	roleLock.Lock()
	defer roleLock.Unlock()

	roles = roleList
	rolesSign = rolesSignList
}
