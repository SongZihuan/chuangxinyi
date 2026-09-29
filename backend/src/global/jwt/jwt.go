package jwt

const (
	UserNotSubTokenPermission = 1 << iota
	UserHighAuthoritySubTokenPermission
	UserRootFatherSubTokenPermission
	UserRootSubTokenPermission
	UserFatherSubTokenPermission
	UserUncleSubTokenPermission
	UserSonSubTokenPermission
	UserWebsiteSubTokenPermission
)

var UserSubTokenStringList = []string{
	UserNotTokenString,
	UserRootTokenString,
	UserSonTokenString,
	UserFatherTokenString,
	UserRootFatherTokenString,
	UserUncleTokenString,
	UserHighAuthorityRootTokenString,
	UserWebsiteTokenString,
}

var UserSubTokenStringChineseMap = map[string]string{
	UserNotTokenString:               "无",
	UserRootTokenString:              "非子用户",
	UserSonTokenString:               "子用户",
	UserFatherTokenString:            "非根父用户",
	UserRootFatherTokenString:        "根父用户",
	UserUncleTokenString:             "协作用户",
	UserHighAuthorityRootTokenString: "非子用户高权限",
	UserWebsiteTokenString:           "外站授权用户",
}

var UserSubTokenStringMap = map[string]int64{
	UserNotTokenString:               UserNotSubTokenPermission,
	UserRootTokenString:              UserRootSubTokenPermission,
	UserSonTokenString:               UserSonSubTokenPermission,
	UserFatherTokenString:            UserFatherSubTokenPermission,
	UserRootFatherTokenString:        UserRootFatherSubTokenPermission,
	UserUncleTokenString:             UserUncleSubTokenPermission,
	UserHighAuthorityRootTokenString: UserHighAuthoritySubTokenPermission,
	UserWebsiteTokenString:           UserWebsiteSubTokenPermission,
}

var UserSubTokenPermissionMap = map[int]int64{
	UserNotToken:               UserNotSubTokenPermission,
	UserRootToken:              UserRootSubTokenPermission,
	UserSonToken:               UserSonSubTokenPermission,
	UserFatherToken:            UserFatherSubTokenPermission,
	UserRootFatherToken:        UserRootFatherSubTokenPermission,
	UserUncleToken:             UserUncleSubTokenPermission,
	UserHighAuthorityRootToken: UserHighAuthoritySubTokenPermission,
	UserWebsiteToken:           UserWebsiteSubTokenPermission,
}

const (
	UserNotToken               = 0
	UserRootToken              = 1
	UserSonToken               = 2
	UserFatherToken            = 3
	UserRootFatherToken        = 4
	UserUncleToken             = 5
	UserHighAuthorityRootToken = 7
	UserWebsiteToken           = 8

	UserNotTokenString               = "UserNotToken"
	UserRootTokenString              = "UserRootToken"              // 非子用户
	UserSonTokenString               = "UserSonToken"               // 子用户
	UserFatherTokenString            = "UserFatherToken"            // 父亲用户登录的子用户（非根父亲）
	UserRootFatherTokenString        = "UserRootFatherToken"        // 根父亲登录的子用户
	UserUncleTokenString             = "UserUncleToken"             // 叔账号登录的子用户
	UserHighAuthorityRootTokenString = "UserHighAuthorityRootToken" // 非子用户（高权限）
	UserWebsiteTokenString           = "UserWebsiteToken"           // 外站
)

var SubTypeMap = map[int]string{
	UserRootToken:              UserRootTokenString,
	UserSonToken:               UserSonTokenString,
	UserFatherToken:            UserFatherTokenString,
	UserRootFatherToken:        UserRootFatherTokenString,
	UserUncleToken:             UserUncleTokenString,
	UserHighAuthorityRootToken: UserHighAuthorityRootTokenString,
	UserWebsiteToken:           UserWebsiteTokenString,
}
