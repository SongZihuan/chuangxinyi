import request from '/@/utils/request';
export function rechargeApi() {
	return {
		//获取充值列表
		getRechargeList: (params?: any) => {
			if (params == null) {
				params = {
					name: '',
					pageNum: 1,
					pageSize: 100000,
				};
			}
			return request({
				url: '/user/center/pay/list',
				method: 'get',
				params,
			});
		},
		//自定义充值
		newRecharge: (data: { cny: number; payWay: string; nvc: string | undefined }) => {
			let nvc = data.nvc
			delete data.nvc

			return request({
				url: '/user/center/pay/selfpay',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Nvc': nvc,
				},
			});
		},
	};
}
