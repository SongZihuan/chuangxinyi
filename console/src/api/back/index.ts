import request from '/@/utils/request';
export function backApi() {
	return {
		//提现列表
		getBackList: (params?: any) => {
			return request({
				url: '/user/center/allow-website/back/list',
				method: 'get',
				params,
			});
		},
	};
}
