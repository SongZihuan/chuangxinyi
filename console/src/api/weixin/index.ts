import request from '/@/utils/request';
import type { weixinAccessTokenTypes, bindWeixinTypes, weixinPayTypes } from './types';

export function useWeixinApi() {
	return {
		//获取微信配置
		weixinConfig: () => {
			return request({
				url: '/public/checker-before/allow-website/wechat/get',
				method: 'get',
			});
		},
		//绑定微信
		bindWeixin: (data: bindWeixinTypes) => {
			return request({
				url: '/user/center/wechat/update',
				method: 'post',
				data,
				headers: {
					'x-token': localStorage.getItem('xtoken'),
				},
			});
		},
		//通过微信code获取access_token
		weixinAccessToken: (data: weixinAccessTokenTypes) => {
			return request({
				url: '/public/checker/allow-website/wechat',
				method: 'post',
				data,
			});
		},
		//微信支付PC端
		weixinPay: (data: weixinPayTypes) => {
			let nvc = data.nvc;
			delete data.nvc;

			return request({
				url: '/user/center/pay/wechatpay/native',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Nvc': nvc,
				},
			});
		},
		//微信支付手机端
		weixinPayMobie: (data: weixinPayTypes) => {
			let nvc = data.nvc;
			delete data.nvc;

			return request({
				url: '/user/center/pay/wechatpay/h5',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Nvc': nvc,
				},
			});
		},
		//微信浏览器支付
		weixinPayJsapi: (data: any) => {
			return request({
				url: '/user/center/pay/wechatpay/jsapi',
				method: 'post',
				data,
			});
		},
	};
}
