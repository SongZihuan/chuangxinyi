package respmsg

type Code struct {
	NumCode int64  `json:"_code"`
	Code    string `json:"code"`
}

var (
	SuccessCode         = Code{NumCode: 0, Code: "SUCCESS"}
	PolicyDenyCode      = Code{NumCode: 1, Code: "POLICY_DENY"}
	CORSDenyCode        = Code{NumCode: 2, Code: "CORS_DENY"}
	TokenDenyCode       = Code{NumCode: 3, Code: "TOKEN_DENY"}
	PeerDenyCode        = Code{NumCode: 4, Code: "PEER_DENY"}
	WebsiteDenyCode     = Code{NumCode: 5, Code: "WEBSITE_DENY"}
	LogicCode           = Code{NumCode: 6, Code: "LOGIC_DENY"}
	RobotDenyCode       = Code{NumCode: 7, Code: "ROBOT_DENY"}
	DoubleCheckDenyCode = Code{NumCode: 8, Code: "DOUBLE_CHECK_DENY"}
	NotTokenDenyCode    = Code{NumCode: 9, Code: "NOT_TOKEN_DENY"}
)
