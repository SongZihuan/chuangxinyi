package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"sync"
)

var roleLock sync.RWMutex
var roles map[int64]warp.Role = make(map[int64]warp.Role, 0)
var rolesSign map[string]warp.Role = make(map[string]warp.Role, 0) // role类型是types.UserRole

func Roles() map[int64]warp.Role {
	roleLock.RLock()
	defer roleLock.RUnlock()
	return roles
}

func RolesSign() map[string]warp.Role {
	roleLock.RLock()
	defer roleLock.RUnlock()
	return rolesSign
}
