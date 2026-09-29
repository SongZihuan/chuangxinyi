import request from '/@/utils/request';
import type { phoneCodeTypes, sliderHeadersTypes, emailCodeTypes } from './types';
export function useBaseApi() {
	return {
		//发送手机验证码
		sendPhoneCode: (data: phoneCodeTypes, headers: sliderHeadersTypes) => {
			return request({
				url: '/public/checker-before/captcha/allow-website/phone/send',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Token': headers.token,
					'X-CAPTCHA-Sig': headers.sig,
					'X-CAPTCHA-SessionId': headers.sessionId,
					'X-CAPTCHA-Scene': 'nc_other',
				},
			});
		},
		//发送邮箱验证码
		sendEmailCode: (data: emailCodeTypes, headers: sliderHeadersTypes) => {
			return request({
				url: '/public/checker-before/captcha/allow-website/email/send',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Token': headers.token,
					'X-CAPTCHA-Sig': headers.sig,
					'X-CAPTCHA-SessionId': headers.sessionId,
					'X-CAPTCHA-Scene': 'nc_other',
				},
			});
		},
		//验证支付是否成功
		checkOrder: (params: { tradeid: string }) => {
			return request({
				url: '/user/center/pay/query',
				method: 'get',
				params,
			});
		},
		checkRefund: (params: { tradeid: string }) => {
			return request({
				url: '/user/center/refund/query',
				method: 'get',
				params,
			});
		},
		//用户申请退款
		refundOrder: (data: { tradeID: string; phoneToken: string }) => {
			return request({
				url: '/user/center/double-check/pay/refund',
				method: 'post',
				data,
				headers: {
					'X-Phone-Token': data.phoneToken,
				},
			});
		},
		//获取前端加密盐值
		salt: () => {
			return request({
				url: '/public/checker-before/salt/get',
				method: 'get',
			});
		},
		// 上传文件
		uploadFile: (data: any) => {
			return request({
				url: '/admin/ui/file/update',
				method: 'post',
				data,
				headers: {
					'Content-Type': 'multipart/form-data',
				},
			});
		},
		// 获取文件
		getFile: (params: any) => {
			return request({
				url: '/public/ui/file',
				method: 'get',
				params,
			});
		},
		// 上传文件默认头像
		uploadHeader: (data: any) => {
			return request({
				url: '/admin/user/header/upload',
				method: 'post',
				data,
				headers: {
					'Content-Type': 'multipart/form-data',
				},
			});
		},
		//获取阿里云appkey
		afsGet: () => {
			return request({
				url: '/public/afs',
				method: 'get',
			});
		},
		//验证手机验证码

		checkPhoneCode: (data: { phone: string,code:string,type:string }) => {
			// @ts-ignore
			return request({
				url: '/public/checker/allow-website/phone',
				method: 'post',
				data,
				isCenter: true,
			});
		},
		// 检验邮箱是否存在
		checkEmailCode: (data: { email: string,code:string,type:string }) => {
			return request({
				url: '/public/checker/allow-website/email',
				method: 'post',
				data,
			});
		},
		// 2fa验证
		check2faCode: (data: { code: string; token: string }) => {
			// @ts-ignore
			return request({
				url: '/public/checker/captcha/allow-website/secondfa',
				method: 'post',
				data: data,
				isCenter: true,
			});
		},
	};
}
