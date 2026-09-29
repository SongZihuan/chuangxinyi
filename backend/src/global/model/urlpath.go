package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"sync"
)

var urlPathLock sync.RWMutex
var urlPathMap map[int64]warp.UrlPath = make(map[int64]warp.UrlPath, 0)

func UrlPathMap() map[int64]warp.UrlPath {
	urlPathLock.RLock()
	defer urlPathLock.RUnlock()
	return urlPathMap
}
