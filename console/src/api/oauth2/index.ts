import request from '/@/utils/request';
import { Local, Session } from '/@/utils/storage';
export function useoauth2Api() {
	return {
		//oauth2登录获取token
		oauth2: (data: { domainUID: any }, logout: Function, isLoginToken: boolean = true) => {
			return request({
				url: '/user/oauth2/',
				method: 'post',
				data,
				isLoginToken: isLoginToken,
				logout: () => {
					Session.clear();
					Local.clear();
					logout();
				},
			});
		},
		//获取域名信息
		oauth2Domain: (params: { domainUID: string }) => {
			return request({
				url: '/public/oauth2',
				method: 'get',
				params,
			});
		},
		// 授权列表
		oauth2List: (params: any) => {
			return request({
				url: '/admin/center/oauth2/list',
				method: 'get',
				params,
			});
		},
		// 取消授权
		oauth2Cancel: (data: { token: string }) => {
			return request({
				url: '/admin/center/oauth2/delete',
				method: 'post',
				data,
			});
		},
		// 清除所有授权
		oauth2Clear: (data: { webID: number }) => {
			return request({
				url: '/admin/center/oauth2/delete/all',
				method: 'post',
				data,
			});
		},
		// 获取授权域名列表
		oauth2DomainList: (params: any) => {
			return request({
				url: '/admin/center/oauth2/website',
				method: 'get',
				params,
			});
		},
		// 获取授权记录
		oauth2Record: (params: any) => {
			return request({
				url: '/admin/center/oauth2/record/list',
				method: 'get',
				params,
			});
		},
		//二次验证
		oauth2Banned: (data: any) => {
			return request({
				url: '/admin/center/oauth2/double-check/banned',
				method: 'post',
				data,
			});
		},

		oauth2Open: (data: { webID: string }) => {
			return request({
				url: '/user/oauth2/open',
				method: 'post',
				data,
			});
		},
	};
}
