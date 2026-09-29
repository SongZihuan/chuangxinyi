package userwstype

const (
	LogoutToken        = "LOGOUT"
	UpdateUserInfo     = "UPDATE_USER_INFO"
	UpdateWalletInfo   = "UPDATE_WALLET_INFO"
	UpdateRoleInfo     = "UPDATE_ROLE_INFO"
	NewAnnouncement    = "NEW_ANNOUNCEMENT"
	UpdateAnnouncement = "UPDATE_ANNOUNCEMENT"
	DeleteAnnouncement = "DELETE_ANNOUNCEMENT"
	UpdateMessage      = "NEW_MESSAGE"
	ReadMessage        = "READ_MESSAGE"
	NewOrderReply      = "NEW_ORDER_REPLY"
	UpdateOrder        = "UPDATE_ORDER"
	RoleChange         = "ROLE_CHANGE"
	Close              = "CLOSE"             // 关闭通道（前端用不到）
	BadCode            = "BAD_CODE"          // 错误的Code
	BadData            = "BAD_DATA"          // 错误的Data
	WithoutToken       = "WITHOUT_TOKEN"     // 没有激活
	WebsiteNotAllow    = "WEBSITE_NOT_ALLOW" // 外站不允许
	BadToken           = "BAD_TOKEN"
	Pong               = "PONG"
)

const (
	GetUserInfo     = "GET_USER_INFO"    // 订阅 UpdateUserInfo RoleChange UpdateRoleInfo
	GetWalletInfo   = "GET_WALLET_INFO"  // 订阅 UpdateWalletInfo
	GetAnnouncement = "GET_ANNOUNCEMENT" // 订阅 NewAnnouncement, UpdateAnnouncement, DeleteAnnouncement
	GetMessage      = "GET_MESSAGE"      // 订阅 UpdateMessage
	GetOrder        = "GET_ORDER"        // 订阅 NewOrderReply UpdateOrder
	CloseOrder      = "CLOSE_ORDER"      // 取消 NewOrderReply UpdateOrder
	Ping            = "PING"
	Token           = "TOKEN" // 发送Token
	Bye             = "BYE"
)
