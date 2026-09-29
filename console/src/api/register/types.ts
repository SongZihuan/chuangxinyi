export interface startRegistrants {
	phoneToken: string;
	inviteID?: number;
}
export interface checkPhoneCodeTypes {
	phone: string;
	email: string;
	code: string;
	type: string;
	inviteID?: number;
}

/**
 * *phoneToken 实名的手机token，不一定是注册的手机
 * * id 注册的手机号
 * *user-name 注册用户名字
 * * user-idcard 注册用户身份证号码
 * *is-company 是否企业(true或者false)
 * *company-name 企业名称
 * *company-id 企业统一社会信用代码
 * *legal-person-name 法人名字
 * *idcard 使用人身份证人像面照片
 */
export interface personInfoFormTypes {
	id: string;
	phoneToken: string;
	faceToken: string;
	userName: string;
	userIdCard: string;
	isCompany: boolean;
	companyName: string;
	companyID: string;
	legalPersonName: string;
}
/**
 *phoneToken 实名的手机token，不一定是注册的手机
 phone 注册的手机号
 * company-name 企业名称
 *company-id 企业统社会信用代码
 * legal-person-name 法人名字
 *legal-person-idcard 法人身份证号码
 *idcard 法人身份证人像面照片
 *license 营业执照照片
 * **/
export interface businessInfoFormTypes {
	id: string;
	phoneToken: string;
	faceToken: string;
	companyName: string;
	companyID: string;
	legalPersonName: string;
	legalPersonID: string;
}
export type CombinedFormTypes = personInfoFormTypes | businessInfoFormTypes;

export interface codeTypes {
	isSend: boolean;
	codeName: string;
	totalTime: number; //一般是60
	timer: any; //定时器
}
export interface uploadUserInfoJsontTypes {
	authenticationMethod?: number;
	userName: string;
	userIDCard: string;
}
export interface updatePassword {}

export interface uploadCompanyInfoJsonTypes {
	legalPersonName: string;
	legalPersonIDCard: string;
	companyName: string;
	companyID: string;
	authenticationMethod: number;
}

export interface passwordTypes {
	password?: string;
	checkPass?: string;
	uid?: string;
	newPasswordHash?: string;
	passwordHash?: string;
	isDelete?: boolean;
}

export interface updatePasswordTypes {
	newPasswordHash: string;
	isDelete: boolean;
}

/**
 * UserName string `json:"username,optional"`
 * IsDelete bool `json:"isDelete,optional"`
 * ID int64 `json:"id,optional"`
 * UID string `json:"uid,optional"`
 */
export interface updateUserNameTypes {
	userName?: string;
	username?: string;
	isDelete?: boolean;
	id?: number;
	uid?: string;
	nickname?: string;
}

export interface tokenExpirationType {
	tokenExpiration: string;
}
