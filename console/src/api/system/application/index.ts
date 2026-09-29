import request from '/@/utils/request';
import type { applicationTypes, applicationDeltypes, searchTypes, applicationMoveTypes } from './types';

export function useApplicationApi() {
	return {
		applicationList: (params?: searchTypes) => {
			return request({
				url: '/admin/application/list',
				method: 'get',
				params,
			});
		},
		websiteAll: () => {
			return request({
				url: '/admin/application/website',
				method: 'get',
			});
		},
		applicationAll: () => {
			return request({
				url: '/public/application',
				method: 'get',
			});
		},
		addAdminapplication: (data: applicationTypes) => {
			return request({
				url: '/admin/application/create',
				method: 'post',
				data,
			});
		},
		applicationUpdate: (data: applicationTypes) => {
			return request({
				url: '/admin/application/update',
				method: 'post',
				data,
			});
		},
		applicationDelete: (data: applicationDeltypes) => {
			return request({
				url: '/admin/application/delete',
				method: 'post',
				data,
			});
		},

		applicationMove: (data: applicationMoveTypes) => {
			return request({
				url: '/admin/application/move',
				method: 'post',
				data,
			});
		},
		// 从数据库更新菜单
		updateapplicationFromDB: () => {
			return request({
				url: '/admin/application/db/update',
				method: 'post',
			});
		},
	};
}
