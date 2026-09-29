import request from '/@/utils/request';
import type { websiteListTypes } from '/@/api/website/site/types';
export function useSiteApi() {
	return {
		//获取登录站点列表
		siteList: (params: any) => {
			return request({
				url: '/admin/website/allow-website/list',
				method: 'get',
				params,
			});
		},
		//新增登录站点
		createSite: (data: websiteListTypes) => {
			return request({
				url: '/admin/website/create',
				method: 'post',
				data,
			});
		},
		//更新登录站点
		updateSite: (data: websiteListTypes) => {
			return request({
				url: '/admin/website/allow-website/update',
				method: 'post',
				data,
			});
		},
		//删除登录站点
		delSite: (data: { id: string }) => {
			return request({
				url: '/admin/website/delete',
				method: 'post',
				data,
			});
		},
		//更新站点公钥
		updateSiteSecret: (data: { id: string; pubkey: string }) => {
			return request({
				url: '/admin/website/allow-website/pubkey/change',
				method: 'post',
				data,
			});
		},
		//获取所有权限值
		allPermissions: () => {
			return request({
				url: '/admin/website/allow-website/all',
				method: 'get',
			});
		},
		// 从数据库更新站点
		updateSiteFromDB: () => {
			return request({
				url: '/admin/website/db/update',
				method: 'post',
			});
		},
	};
}
