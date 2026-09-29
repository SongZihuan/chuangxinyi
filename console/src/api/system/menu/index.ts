import request from '/@/utils/request';
import type { menuTypes, menuDeltypes, searchTypes, menuMoveTypes, roleActionTypes } from './types';
import { MenuRauthTypes } from '/@/views/system/menu/types';

export function useMenuApi() {
	return {
		menuList: (params?: searchTypes) => {
			return request({
				url: '/admin/menu/list',
				method: 'get',
				params,
			});
		},
		subAll: (params?: searchTypes) => {
			return request({
				url: '/admin/path/sub/all',
				method: 'get',
				params,
			});
		},
		addAdminMenu: (data: menuTypes) => {
			return request({
				url: '/admin/menu/create',
				method: 'post',
				data,
			});
		},
		menuUpdate: (data: menuTypes) => {
			return request({
				url: '/admin/menu/update',
				method: 'post',
				data,
			});
		},
		menuDelete: (data: menuDeltypes) => {
			return request({
				url: '/admin/menu/delete',
				method: 'post',
				data,
			});
		},
		//菜单添加角色
		menuAddRole: (data: roleActionTypes) => {
			return request({
				url: '/admin/menu//role/add',
				method: 'post',
				data,
			});
		},
		//菜单删除角色
		menuDeleteRole: (data: roleActionTypes) => {
			return request({
				url: '/admin/menu//role/delete',
				method: 'post',
				data,
			});
		},
		//移动菜单
		menuMove: (data: menuMoveTypes) => {
			return request({
				url: '/admin/menu/move',
				method: 'post',
				data,
			});
		},
		// 从数据库更新菜单
		updateMenuFromDB: () => {
			return request({
				url: '/admin/menu/db/update',
				method: 'post',
			});
		},
		// 更新角色权限 /policy/update
		updateRolePolicy: (data: MenuRauthTypes) => {
			return request({
				url: '/admin/menu/policy/update',
				method: 'post',
				data,
			});
		},
	};
}
