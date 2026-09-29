/**
 * UID string `json:"id"`
 *     RoleID int64 `json:"roleID,omitempty"`
 *     RoleName string `json:"roleName"`
 *     RoleSign string `json:"roleSign"`
 *     Phone string `json:"phone"`
 *     UserName string `json:"userName,omitempty"`
 *     NickName string `json:"nickname,omitempty"`
 *     Header string `json:"header,omitempty"`
 *     Email string `json:"email,omitempty"`
 *     UserRealName string `json:"userRealName,omitempty"`
 *     CompanyName string `json:"companyName,omitempty"`
 *     WeChatNickName string `json:"wechatNickName,omitempty"`
 *     WeChatHeader string `json:"wechatHeader,omitempty"`
 *     UnionID string `json:"unionID,omitempty"`
 *     Signin bool `json:"signin"`
 *     Status string `json:"status"`
 *     InviteCount int64 `json:"inviteCount"`
 *     LastInviteAt int64 `json:"lastInviteAt"`
 *     CreateAt int64 `json:"createAt"`
 */
export interface FatherUser {
	id: string;
	roleID: number;
	roleName: string;
	roleSign: string;
	phone: string;
	userName: string;
	nickname: string;
	header: string;
	email: string;
	userRealName: string;
	companyName: string;
	wechatNickName: string;
	wechatHeader: string;
	unionID: string;
	signin: boolean;
	status: string;
	inviteCount: number;
	lastInviteAt: number;
	createAt: number;
	lineal: boolean;
	hasChildren?: boolean;
	son: FatherUser[];
}
/**
 * UserEasy
 * Lineal bool `json:"lineal"`  // 直系亲属
 * Son []FatherUser `json:"son"` // 子账号
 */
export interface tableTypes {
	data: FatherUser[];
	loading: boolean;
	param: {
		id: string;
		uid: string;
	};
}
export interface subUserAccountStatsTypes {
	tableData: tableTypes;
}

/**
 * roleList
 */
export interface roleListTypes {
	data: any;
	params: {
		page: number;
		pageSize: number;
	};
}
export interface checkPhoneTypes {
	phone: string;
	code: string;
	type: string;
	newWallet: boolean;
}
