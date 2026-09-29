import request from '/@/utils/request';
import type { totpConfigTypes, bindTotpTypes } from './types';
export function useGoogleAuthApi() {
	return {
		//获取双因素验证器二维码
		totpUrl: (params: totpConfigTypes) => {
			return request({
				url: '/public/checker-before/totp/get',
				method: 'get',
				params: params,
			});
		},
		//绑定2fa
		bindTotp: (data: bindTotpTypes) => {
			return request({
				url: '/user/center/2fa/bind',
				method: 'post',
				data: data,
			});
		},
		//删除2fa
		delTotp: (data: { code?: string }) => {
			return request({
				url: '/user/center/2fa/delete/secret',
				method: 'post',
				data: data,
			});
		},
		rootDelTotp: (data: {}) => {
			return request({
				url: '/user/center/root-only/2fa/delete',
				method: 'post',
				data: data || {},
			});
		},
		//二次验证
		checkTotp: (data: { code: string; token: string; nvc: string | undefined; rememberHour: number; passToken?: string }) => {
			let nvc = data.nvc;
			delete data.nvc;

			return request({
				url: '/public/checker-after/2fa',
				method: 'post',
				data: data,
				headers: {
					'X-CAPTCHA-Nvc': nvc,
				},
			});
		},
	};
}
