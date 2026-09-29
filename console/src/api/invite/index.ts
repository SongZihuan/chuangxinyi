import request from '/@/utils/request';

export function useInviteApi() {
	return {
		//获取邀请人信息
		inviterInfo: () => {
			return request({
				url: '/user/center/invite',
				method: 'get',
			});
		},
		//获取父账号信息
		fatherInfo: () => {
			return request({
				url: '/user/center/father',
				method: 'get',
			});
		},
		//获取当前用户获取邀请列表
		inviterList: () => {
			return request({
				url: '/user/center/invite/list',
				method: 'get',
			});
		},
	};
}
