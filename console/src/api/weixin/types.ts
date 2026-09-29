import { LocationQueryValue } from 'vue-router';
export interface bindWeixinTypes {
	wechatToken?: string;
	isDelete: boolean;
}

export interface weixinAccessTokenTypes {
	code: string | LocationQueryValue[];
	type: string;
}

export interface weixinPayTypes {
	cny: number;
	payType: number;
	couponsID?: string | null;
	payMode?: number;
	nvc: string | undefined;
}
