//发送验证码
export interface phoneCodeTypes {
	phone: string;
}
//阿里云滑块参数
export interface sliderHeadersTypes {
	sig: string;
	sessionId: string;
	token: string;
}

export interface emailCodeTypes {
	email: string;
}

export interface checkEmailTypes {
	email: string;
	code: string;
	type: string;
}
export interface checkPhoneTypes {
	phone: string;
	code: string;
	type: string;
}
export interface checkWxRobotTypes {
	webhook: string;
	isDelete: boolean;
}
