/**
 * homePage类型
 * company
 * introduction // 用户简介
 * address      // 联系地址
 * phone        // 联系电话
 * email        // 联系邮件
 * wechat       // 微信号
 * qq           // QQ号
 * sex          // 是否男性
 * link         // 外部连接
 * industry     // 行业
 * position     // 职位
 */
export interface homePageTypes {
	company: string;
	introduction: string;
	address: string;
	phone: string;
	email: string;
	wechat: string;
	qq: string;
	sex: string;
	link: string;
	industry: string;
	position: string;
	isDelete?: boolean;
}

/**
 * userID
 */
export interface pageHomeParamsTypes {
	userID?: string;
	isDelete?: boolean;
	uid?: string;
}

/**
 * 详细信息
 * ID int64  `json:"id"`  // 用户数字ID
 *     UID string `json:"uid"`  // 用户ID
 *     Status string `json:"status"`  // 状态
 *     Signin bool `json:"signin"`  // 是否单点登录
 *     Father int64 `json:"father"`  // 父亲数字ID
 *     tokenExpiration int64 `json:"tokenExpiration"`  // 登录有效时长
 *     RoleID int64 `json:"roleID"`  // 角色ID
 *     RoleName string `json:"roleName"`  // 角色名称
 *     RoleSign string `json:"roleSign"`  // 角色标识
 *     IsAdmin bool `json:"isAdmin"`  // 是否根管理员
 *     CreateAt int64 `json:"createAt"`  // 注册时间
 *
 *     Phone string `json:"phone"`  // 手机号
 *     Email string `json:"email"`  // 邮箱
 *
 *     Nickname string `json:"nickname"`  // 昵称
 *     Header string `json:"header"`  // 头像
 *
 *     WxOpenID string `json:"wxOpenID"`  // 微信OpenID
 *     WxUnionID string `json:"wxUnionID"`  // 微信UnionID
 *     FuwuhaoOpenID string `json:"fuwuhaoOpenID"`  // 服务号OpenID
 *     WxNickName string `json:"wxNickName"`  // 微信昵称
 *     WxHeader string `json:"wxHeader"`  // 微信头像
 *
 *     WxWebHook string `json:"wxWebHook"`  // 企业微信webhook
 *
 *     HasPassword bool `json:"hasPassword"`  // 是否有密码
 *
 *     UserName string `json:"userName"`  // 用户名
 *
 *     Has2FA bool `json:"has2FA"`  // 是否有2FA
 *
 *     AddressName     string `json:"addressName"`  // 地址收件人
 *     AddressPhone    string `json:"addressPhone"`  // 地址手机手机号
 *     AddressEmail    string `json:"addressEmail"`  // 地址手机邮箱
 *     AddressProvince string `json:"addressProvince"`  // 省份
 *     AddressCity     string `json:"addressCity"`  // 城市
 *     AddressDistrict string `json:"addressDistrict"`  // 区县
 *     AddressAddress  string `json:"addressAddress"`  // 详细地址
 *     AddressArea []string `json:"addressArea"`
 */
export interface userDataTypes {
	id: number;
	uid: string;
	status: string;
	signin: boolean;
	father: number;
	tokenExpiration: number;
	roleID: number;
	roleName: string;
	roleSign: string;
	isAdmin: boolean;
	createAt: number;
	phone: string;
	email: string;
	nickname: string;
	header: string;
	wxOpenID: string;
	wxUnionID: string;
	fuwuhaoOpenID: string;
	wxNickName: string;
	wxHeader: string;
	wxWebHook: string;
	hasPassword: boolean;
	userName: string;
	has2FA: boolean;
	addressName: string;
	addressPhone: string;
	addressEmail: string;
	addressProvince: string;
	addressCity: string;
	addressDistrict: string;
	addressAddress: string;
	addressArea: string[];
}