package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"sync"
)

var menuLock sync.RWMutex
var menus map[int64]warp.Menu = make(map[int64]warp.Menu, 0)

func Menus() map[int64]warp.Menu {
	menuLock.RLock()
	defer menuLock.RUnlock()
	return menus
}
