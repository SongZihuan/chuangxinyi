import request from '/@/utils/request';
import type { roleTypes, roleMoveTypes } from './types';
export function roleApi() {
	return {
		//获取角色列表
		getRoleList: (params?: any) => {
			if (params == null) {
				params = {
					name: '',
					pageNum: 1,
					pageSize: 100000,
				};
			}
			return request({
				url: '/admin/role/list',
				method: 'get',
				params,
			});
		},
		//添加角色
		roleAdd: (data?: roleTypes) => {
			return request({
				url: '/admin/role/create',
				method: 'post',
				data,
			});
		},
		//更新角色
		roleUpdate: (data?: roleTypes) => {
			return request({
				url: '/admin/role/update',
				method: 'post',
				data,
			});
		},
		//删除角色
		roleDelAction: (data: { id: number }) => {
			return request({
				url: '/admin/role/delete',
				method: 'post',
				data,
			});
		},
		//移动角色
		roleMove: (data: roleMoveTypes) => {
			return request({
				url: '/admin/role/move',
				method: 'post',
				data,
			});
		},
		//获取所有内置权限
		allPermissions: () => {
			return request({
				url: '/admin/role/all',
				method: 'get',
			});
		},
		// 从数据库更新角色
		updateRoleFromDB: () => {
			return request({
				url: '/admin/role/db/update',
				method: 'post',
			});
		},
	};
}
