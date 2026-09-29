import request from '/@/utils/request';

export function usePathApi() {
	return {
		pathList: (params?: any) => {
			return request({
				url: '/admin/path/list',
				method: 'get',
				params,
			});
		},
		addPathMenu: (data: any) => {
			return request({
				url: '/admin/path/create',
				method: 'post',
				data,
			});
		},
		pathUpdate: (data: any) => {
			return request({
				url: '/admin/path/update',
				method: 'post',
				data,
			});
		},
		pathDelete: (data: any) => {
			return request({
				url: '/admin/path/delete',
				method: 'post',
				data,
			});
		},
		updatePathFromDB: () => {
			return request({
				url: '/admin/path/db/update',
				method: 'post',
			});
		},
		pathAddRole: (data: any) => {
			return request({
				url: '/admin/path/role/add',
				method: 'post',
				data,
			});
		},
		pathDeleteRole: (data: any) => {
			return request({
				url: '/admin/path/role/delete',
				method: 'post',
				data,
			});
		},
		pathMove: (data: any) => {
			return request({
				url: '/admin/path/move',
				method: 'post',
				data,
			});
		},
		updateMenuFromDB: () => {
			return request({
				url: '/admin/path/db/update',
				method: 'post',
			});
		},
		updateRolePolicy: (data: any) => {
			return request({
				url: '/admin/path/policy/update',
				method: 'post',
				data,
			});
		},
		//获取所有权限值
		allPermissions: () => {
			return request({
				url: '/admin/path/all',
				method: 'get',
			});
		},
	};
}
