import request from '/@/utils/request';
export function withdrawApi() {
	return {
		//提现列表
		getwithdrawList: (params?: any) => {
			return request({
				url: '/user/center/withdraw/list',
				method: 'get',
				params,
			});
		},
	};
}
