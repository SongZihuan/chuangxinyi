package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"sync"
)

var websiteLock sync.RWMutex
var websiteList []warp.Website = make([]warp.Website, 0)
var websites map[int64]warp.Website = make(map[int64]warp.Website, 0)
var websitesUID map[string]warp.Website = make(map[string]warp.Website) // 登录站点列表

func WebsiteList() []warp.Website {
	websiteLock.RLock()
	defer websiteLock.RUnlock()
	return websiteList
}

func Websites() map[int64]warp.Website {
	websiteLock.RLock()
	defer websiteLock.RUnlock()
	return websites
}

func WebsitesUID() map[string]warp.Website {
	websiteLock.RLock()
	defer websiteLock.RUnlock()
	return websitesUID
}
