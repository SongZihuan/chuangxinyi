package cron

import (
	"gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/peers/wstype"
)

func PathUpdateHandler(update bool) {
	UrlPathHandler()
	PermissionHandler()
	MenuHandler()
	RoleHandler()

	if update {
		websocket.WritePeersMessage(wstype.PeersUrlPathDBUpdate, struct{}{}, websocket.WSMessage{})
	}
}

func PermissionUpdateHandler(update bool) {
	PermissionHandler()
	UrlPathHandler()
	MenuHandler()
	RoleHandler()

	if update {
		websocket.WritePeersMessage(wstype.PeersPermissionDBUpdate, struct{}{}, websocket.WSMessage{})
	}
}

func MenuUpdateHandler(update bool) {
	MenuHandler()
	UrlPathHandler()
	PermissionHandler()
	RoleHandler()

	if update {
		websocket.WritePeersMessage(wstype.PeersMenuDBUpdate, struct{}{}, websocket.WSMessage{})
	}
}

func RoleUpdateHandler(update bool) {
	RoleHandler()
	UrlPathHandler()
	PermissionHandler()
	MenuHandler()

	if update {
		websocket.WritePeersMessage(wstype.PeersRoleDBUpdate, struct{}{}, websocket.WSMessage{})
	}
}

func WebsiteUpdateHandler(update bool) {
	WebsiteHandler()
	WebsitePermissionHandler()
	WebsiteUrlPathHandler()
	ApplicationHandler()

	if update {
		websocket.WritePeersMessage(wstype.PeersWebsiteDBUpdate, struct{}{}, websocket.WSMessage{})
	}
}

func WebsitePathUpdateHandler(update bool) {
	WebsiteUrlPathHandler()
	WebsiteHandler()
	WebsitePermissionHandler()
	ApplicationHandler()

	if update {
		websocket.WritePeersMessage(wstype.PeersWebsiteUrlPathDBUpdate, struct{}{}, websocket.WSMessage{})
	}
}

func WebsitePermissionUpdateHandler(update bool) {
	WebsiteUrlPathHandler()
	WebsiteHandler()
	WebsitePermissionHandler()
	ApplicationHandler()

	if update {
		websocket.WritePeersMessage(wstype.PeersWebsitePermissionDBUpdate, struct{}{}, websocket.WSMessage{})
	}
}

func ApplicationUpdateHandler(update bool) {
	ApplicationHandler()

	if update {
		websocket.WritePeersMessage(wstype.PeersApplicationDBUpdate, struct{}{}, websocket.WSMessage{})
	}
}
