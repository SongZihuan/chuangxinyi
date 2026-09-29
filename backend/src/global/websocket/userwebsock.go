package websocket

import "sync"

var UserConnMap = make(map[int64][]chan WSMessage, 10)
var UserConnMapMutex = new(sync.Mutex)

var TokenConnMap = make(map[string][]chan WSMessage, 10)
var TokenConnMapMutex = new(sync.Mutex)

var RoleConnMap = make(map[int64][]chan WSMessage, 10)
var RoleConnMapMutex = new(sync.Mutex)

var WalletConnMap = make(map[int64][]chan WSMessage, 10)
var WalletConnMapMutex = new(sync.Mutex)

var AnnouncementConnList = make([]chan WSMessage, 0, 20)
var AnnouncementConnListMutex = new(sync.Mutex)

var MessageConnMap = make(map[int64][]chan WSMessage, 10)
var MessageConnMapMutex = new(sync.Mutex)

var OrderConnMap = make(map[string][]chan WSMessage, 10)
var OrderConnMapMutex = new(sync.Mutex)
