import request from '/@/utils/request';

export function usePermissionApi() {
	return {
		permissionList: (params?: any) => {
			return request({
				url: '/admin/website-permission/list',
				method: 'get',
				params,
			});
		},
		addPermissionMenu: (data: any) => {
			return request({
				url: '/admin/website-permission/create',
				method: 'post',
				data,
			});
		},
		permissionUpdate: (data: any) => {
			return request({
				url: '/admin/website-permission/update',
				method: 'post',
				data,
			});
		},
		permissionDelete: (data: any) => {
			return request({
				url: '/admin/website-permission/delete',
				method: 'post',
				data,
			});
		},
		updatePermissionFromDB: () => {
			return request({
				url: '/admin/website-permission/db/update',
				method: 'post',
			});
		},
		updateMenuFromDB: () => {
			return request({
				url: '/admin/website-permission/db/update',
				method: 'post',
			});
		},
	};
}
