package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"math/big"
	"sync"
)

var websitePermissionLock sync.RWMutex
var websitePermissionsSign map[string]warp.WebsitePermission = make(map[string]warp.WebsitePermission)
var websitePermissionList []warp.WebsitePermission = make([]warp.WebsitePermission, 0)
var websiteAllPermission *big.Int = big.NewInt(0)

func WebsitePermissionsSign() map[string]warp.WebsitePermission {
	websitePermissionLock.RLock()
	defer websitePermissionLock.RUnlock()
	return websitePermissionsSign
}

func WebsitePermissionList() []warp.WebsitePermission {
	websitePermissionLock.RLock()
	defer websitePermissionLock.RUnlock()
	return websitePermissionList
}

func WebsiteAllPermission() *big.Int {
	websitePermissionLock.RLock()
	defer websitePermissionLock.RUnlock()
	return websiteAllPermission
}
