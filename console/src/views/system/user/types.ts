import { subUserInfoTypes } from '/@/api/user/subuser/types';

/**
 * 用户信息
 */
export interface userInfosTypes {
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
export interface userPhoneTypes {
	phone: string;
	uid?: string;
	id?: number;
}
export interface userEmailTypes {
	email: string;
	uid?: string;
	id?: number;
}
export interface userStatusTypes {
	status: string;
	uid: string;
}

/**
 * tradeID: '',
 *   subject: '',
 *   cny: '',
 *   companyID: '',
 *   tradeStatus: '',
 *   payWay: '',
 *   createTime: '',
 *   payAt: '',
 */
export interface tradeTypes {
	tradeID: string;
	subject: string;
	cny: string;
	companyID: string;
	tradeStatus: number;
	payWay: string;
	createAt: string;
	payAt: string;
}
/**
 * <!--
 *         UserName string `json:"userName"`  // 用户实名
 *             UserIDCard string `json:"userIDCard"`  // 用户身份证
 *
 *             CompanyName string `json:"companyName"`  // 企业实名
 *             CompanyID string `json:"companyID"`  // 企业统一社会信用代码
 *             LegalPersonName string `json:"legalPersonName"`  // 法人姓名
 *             LegalPersonIDCard string `json:"legalPersonIDCard"`  // 法人身份证
 *         -->
 */
export interface userRealNameTypes {
	userName: string;
	userIDCard: string;
	companyName: string;
	companyID: string;
	legalPersonName: string;
	legalPersonIDCard: string;
}

/**
 *  WalletID int64 `json:"walletID"`
 *     Balance int64 `json:"balance"`  // 余额
 *     NotBilled int64 `json:"notBilled"`  // 未开票金额
 *     Billed int64 `json:"billed"`  // 总共可开票金额
 *     HasBilled int64 `json:"hasBilled"`  // 已开票金额
 *
 *     TitleName string `json:"titleName"`  // 抬头名称
 *     TitleTaxID string `json:"titleTaxID"`  // 抬头税号（个人身份证号码，企业税号）
 *     TitleBankID string `json:"titleBannedID"`  // 抬头银行账户
 *     TitleBank string `json:"titleBank"`  // 抬头开户行
 */
export interface userWalletTypes {
	walletID: number;
	balance: number;
	notBilled: number;
	billed: number;
	hasBilled: number;
	titleName?: string;
	titleTaxID?: string;
	titleBankID?: string;
	titleBank?: string;
}
export interface userTypes {
	userName: string; // 账户名称
	nikeName: string; // 用户昵称
	roleSign: number[]; // 关联角色
	phone: string; // 手机号
	userType: number; //类型
	email: string; // 邮箱
	idcard: string; //身份证号
	vipLevel: number; //会员等级
	socialCode: string; //统一社会信用代码
	sex: number; // 性别
	password: string; // 账户密码
	overdueTime: string; // 账户过期
	userStatus: number; // 用户状态
	describe?: string; // 用户描述
	socialName?: string; //公司名称
}

export interface searchTypes {
	page: number;
	pagesize: number;
	range: any[];
	starttime?: number | null;
	endtime?: number | null;
}

/** 子账户
 * AdminFatherUser {
 *     UserID int64 `json:"userID"`  // 用户数字ID
 *     Phone string `json:"phone"`  // 手机号
 *     Email string `json:"email"`  // 邮箱
 *     UserName string `json:"userName"`  // 用户名
 *     RoleID int64 `json:"roleID"`  // 角色ID
 *     RoleName string `json:"roleName"`  // 角色名称
 *     RoleSign string `json:"roleSign"`  // 角色标识
 *     Lineal bool `json:"lineal"`  // 直系亲属
 *     Son []AdminFatherUser `json:"son"`  // 子用户
 * }
 */
export interface subAccountTypes {
	userID: number;
	phone: string;
	email: string;
	userName: string;
	roleID: number;
	roleName: string;
	roleSign: string;
	lineal: boolean;
	hasChildren?: boolean;
	length?: number;
	son: subAccountTypes[];
}

/**
 * data: [],
 *     total: 0,
 *     loading: false,
 *     param: {
 *       name: '',
 *       page: 1,
 *       pageSize: 10,
 *       total: 0,
 *     },
 */
export interface tableTypes {
	data: subAccountTypes[];
	total: number;
	loading: boolean;
	param: subUserInfoTypes;
}
export interface subAccountStatsTypes {
	tableData: tableTypes;
}
