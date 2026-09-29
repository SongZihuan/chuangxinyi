package model

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"sync"
)

var applicationLock sync.RWMutex
var applicationName map[string]warp.Application = make(map[string]warp.Application)
var applicationList []warp.Application = make([]warp.Application, 0)

func ApplicationName() map[string]warp.Application {
	applicationLock.RLock()
	defer applicationLock.RUnlock()
	return applicationName
}

func ApplicationList() []warp.Application {
	applicationLock.RLock()
	defer applicationLock.RUnlock()
	return applicationList
}
