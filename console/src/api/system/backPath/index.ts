import request from '/@/utils/request';

export function usePathApi() {
	return {
		pathList: (params?: any) => {
			return request({
				url: '/admin/website-path/list',
				method: 'get',
				params,
			});
		},
		addAdminpath: (data: any) => {
			return request({
				url: '/admin/website-path/create',
				method: 'post',
				data,
			});
		},
		pathUpdate: (data: any) => {
			return request({
				url: '/admin/website-path/update',
				method: 'post',
				data,
			});
		},
		pathDelete: (data: any) => {
			return request({
				url: '/admin/website-path/delete',
				method: 'post',
				data,
			});
		},

		// 从数据库更新菜单
		updatepathFromDB: () => {
			return request({
				url: '/admin/website-path/db/update',
				method: 'post',
			});
		},
		//获取所有权限值
		allPermissions: () => {
			return request({
				url: '/admin/website-path/all',
				method: 'get',
			});
		},
	};
}
