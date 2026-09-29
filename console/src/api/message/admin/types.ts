export interface messageAdminTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
	id?: string;
	uid?: string;
}

/**
 * 手机号
 */
export interface smsMessageAdminTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
	phone?: string;
}

/**
 * email
 */
export interface emailMessageAdminTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
	email?: string;
}

/**
 * fuwuhao
 */
export interface fuwuhaoMessageAdminTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
	openID?: string;
}

/**
 * wxrobot
 */
export interface wxrobotMessageAdminTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
	webhook?: string;
}
export interface sendMsgTypes {
	title: string | null;
	message: string | null;
	uid: number | string | null;
}
export interface sendSmsTypes {
	templateID: string;
	sig?: string;
	data?: dictTypes[];
	id?: number | null;
	uid?: number | null;
}
export interface sendEmailTypes {
	subject: string;
	sender: string;
	content: string;
	id: number | null;
}

export interface sendwxrobotTypes {
	content: string;
	id: number | null;
}
export interface sendFuwuhaoTypes extends sendSmsTypes {
	val?: dictTypes[];
}
