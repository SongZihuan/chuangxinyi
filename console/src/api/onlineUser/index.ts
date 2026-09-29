import request from '/@/utils/request';

export default function useOnlineUser() {
	return {
		// 获取在线用户列表
		getOnlineUserList: () => {
			return request({
				url: 'admin/center/token/all',
				method: 'get',
			});
		},
		// 删除所有在线用户
		deleteAllOnlineUser: () => {
			return request({
				url: 'admin/center/token/all/delete',
				method: 'post',
			});
		},
		// 踢下线
		kickOnlineUser: (data: any) => {
			return request({
				url: 'admin/center/token/other/delete',
				method: 'post',
				data,
			});
		},
		// 退出登录
		deleteOnlineUser: (data: any) => {
			return request({
				url: 'user/center/token/delete',
				method: 'post',
				data,
			});
		},
		// 更新单点登录状态
		updateOnlineUser: (data: { signin: boolean }) => {
			return request({
				url: 'user/center/signin/update',
				method: 'post',
				data,
			});
		},
	};
}
