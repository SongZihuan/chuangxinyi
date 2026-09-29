package websocket

import (
	"gitee.com/wuntsong-auth/backend/src/utils"
	"sync"
)

var PeersConnMap = make(map[string]chan WSPeersMessage, 10)
var PeersConnMapMutex = new(sync.Mutex)
var PeersCounter = utils.NewCounter()
