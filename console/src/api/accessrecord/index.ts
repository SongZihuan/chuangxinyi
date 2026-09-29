import request from '/@/utils/request';
import { accessrecordListTypes } from '/@/api/accessrecord/types';
export function useAccessrecordApi() {
	return {
		accessrecordList(params: accessrecordListTypes) {
			return request({
				url: '/admin/accessrecord/access/list',
				method: 'get',
				params,
			});
		},
		// token访问记录
		accessrecordTokenList(params: accessrecordListTypes) {
			return request({
				url: '/admin/accessrecord/token/list',
				method: 'get',
				params,
			});
		},
		// 根据token获取访问记录
		accessrecordTokenGet(params: any) {
			return request({
				url: '/admin/accessrecord/access/token/list',
				method: 'get',
				params,
			});
		},
		// 根据请求id获取访问记录
		accessrecordIdGet(params: any) {
			return request({
				url: '/admin/accessrecord/access/info',
				method: 'get',
				params,
			});
		},
		// 根据sql获取访问记录
		accessrecordSqlGet(params: any) {
			return request({
				url: '/admin/accessrecord/access/cond/list',
				method: 'get',
				params,
			});
		},
		// 根据token 获取token访问记录
		accessrecordToken(params: any) {
			return request({
				url: '/admin/accessrecord/token/token/list',
				method: 'get',
				params,
			});
		}
	};
}
