package wstype

const (
	PeersLogoutToken               = "LOGOUT"
	PeersUpdateUserInfo            = "UPDATE_USER_INFO"
	PeersUpdateWalletInfo          = "UPDATE_WALLET_INFO"
	PeersUpdateRoleInfo            = "UPDATE_ROLE_INFO"
	PeersUpdateAnnouncement        = "UPDATE_ANNOUNCEMENT"
	PeersUpdateMessage             = "UPDATE_MESSAGE"
	PeersUpdateOrder               = "UPDATE_ORDER"
	PeersRoleChange                = "ROLE_CHANGE"
	PeersRoleChangeByDelete        = "ROLE_CHANGE_DELETE"
	PeersRoleDBUpdate              = "DB_ROLE"
	PeersMenuDBUpdate              = "DB_MENU"
	PeersPermissionDBUpdate        = "DB_PERMISSION"
	PeersUrlPathDBUpdate           = "DB_URL_PATH"
	PeersWebsiteDBUpdate           = "DB_WEBSITE"
	PeersWebsiteUrlPathDBUpdate    = "DB_WEBSITE_URL_PATH"
	PeersWebsitePermissionDBUpdate = "DB_WEBSITE_PERMISSION"
	PeersApplicationDBUpdate       = "DB_APPLICATION"
	PeersClose                     = "CLOSE"    // 关闭通道（前端用不到）
	PeersBadCode                   = "BAD_CODE" // 错误的Code
	PeersBadData                   = "BAD_DATA" // 错误的Data
	PeersPing                      = "PING"
)

var PeersForward = []string{
	PeersLogoutToken,
	PeersUpdateUserInfo,
	PeersUpdateWalletInfo,
	PeersUpdateRoleInfo,
	PeersUpdateAnnouncement,
	PeersUpdateMessage,
	PeersUpdateOrder,
	PeersRoleChange,
	PeersRoleChangeByDelete,
	PeersRoleDBUpdate,
	PeersMenuDBUpdate,
	PeersWebsiteDBUpdate,
}
