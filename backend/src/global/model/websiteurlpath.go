package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"sync"
)

var websiteUrlPathLock sync.RWMutex
var websiteUrlPathMap map[int64]warp.WebsiteUrlPath = make(map[int64]warp.WebsiteUrlPath, 0)

func WebsiteUrlPathMap() map[int64]warp.WebsiteUrlPath {
	websiteUrlPathLock.RLock()
	defer websiteUrlPathLock.RUnlock()
	return websiteUrlPathMap
}
