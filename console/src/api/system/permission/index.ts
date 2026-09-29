import request from '/@/utils/request';

export function usePermissionApi() {
	return {
		permissionList: (params?: any) => {
			return request({
				url: '/admin/permission/list',
				method: 'get',
				params,
			});
		},
		addPermissionMenu: (data: any) => {
			return request({
				url: '/admin/permission/create',
				method: 'post',
				data,
			});
		},
		permissionUpdate: (data: any) => {
			return request({
				url: '/admin/permission/update',
				method: 'post',
				data,
			});
		},
		permissionDelete: (data: any) => {
			return request({
				url: '/admin/permission/delete',
				method: 'post',
				data,
			});
		},
		updatePermissionFromDB: () => {
			return request({
				url: '/admin/permission/db/update',
				method: 'post',
			});
		},
		//菜单添加角色
		permissionAddRole: (data: any) => {
			return request({
				url: '/admin/permission/role/add',
				method: 'post',
				data,
			});
		},
		//菜单删除角色
		permissionDeleteRole: (data: any) => {
			return request({
				url: '/admin/permission/role/delete',
				method: 'post',
				data,
			});
		},
		//移动菜单
		permissionMove: (data: any) => {
			return request({
				url: '/admin/permission/move',
				method: 'post',
				data,
			});
		},
		// 从数据库更新菜单
		updateMenuFromDB: () => {
			return request({
				url: '/admin/permission/db/update',
				method: 'post',
			});
		},
		// 更新角色权限 /policy/update
		updateRolePolicy: (data: any) => {
			return request({
				url: '/admin/permission/policy/update',
				method: 'post',
				data,
			});
		},
	};
}
