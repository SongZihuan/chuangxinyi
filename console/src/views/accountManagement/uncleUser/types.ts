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
export interface uncleUserTypes {
	id: string;
	roleID: number;
	roleName: string;
	roleSign: string;
	phone: string;
	userName: string;
	nickName: string;
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
}

export interface uncleUserListTypes {
	id: string;
	roleID: number;
	roleName: string;
	roleSign: string;
	phone: string;
	header: string;
	userRealName: string;
	signin: boolean;
	status: string;
	inviteCount: number;
	lastInviteAt: number;
	createAt: number;
}
interface tableTypes {
	data: uncleUserTypes[];
	loading: boolean;
}
export interface uncleAccountStatsTypes {
	tableData: tableTypes;
}

export interface addUncleTypes {
	uncleID: string;
}
