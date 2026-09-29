import request from '/@/utils/request';
import type { alipayTypes } from './types';
export function useAlipayApi() {
	return {
		//支付宝支付PC端
		alipay: (data: alipayTypes) => {
			let nvc = data.nvc;
			delete data.nvc;

			return request({
				url: '/user/center/pay/alipay/pc',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Nvc': nvc,
				},
			});
		},
		//支付宝支付手机端
		wapAlipay: (data: alipayTypes) => {
			let nvc = data.nvc;
			delete data.nvc;

			return request({
				url: '/user/center/pay/alipay/wap',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Nvc': nvc,
				},
			});
		},
	};
}
