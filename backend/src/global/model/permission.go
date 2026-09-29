package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"math/big"
	"sync"
)

var permissionLock sync.RWMutex
var permissionsSign map[string]warp.Permission = make(map[string]warp.Permission, 0)
var permissionList []warp.Permission = make([]warp.Permission, 0)
var allPermission *big.Int = big.NewInt(0)
var userPermission *big.Int = big.NewInt(0)
var anonymousPermission *big.Int = big.NewInt(0)

func PermissionsSign() map[string]warp.Permission {
	permissionLock.RLock()
	defer permissionLock.RUnlock()
	return permissionsSign
}

func PermissionList() []warp.Permission {
	permissionLock.RLock()
	defer permissionLock.RUnlock()
	return permissionList
}

func AllPermission() *big.Int {
	permissionLock.RLock()
	defer permissionLock.RUnlock()
	return allPermission
}

func UserPermission() *big.Int {
	permissionLock.RLock()
	defer permissionLock.RUnlock()
	return userPermission
}

func AnonymousPermission() *big.Int {
	permissionLock.RLock()
	defer permissionLock.RUnlock()
	return anonymousPermission
}
