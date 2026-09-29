import request from '/@/utils/request';
import {bannedTypes, bannedUpdateTypes} from '/@/api/audit/types';
export function useBannedApi() {
	return {
		// 用户获取审计列表
		getBannedList(params: bannedTypes) {
			return request({
				url: '/admin/center/oauth2/banned/list',
				method: 'get',
				params,
			});
		},
		// 更新
		updateBanned(data: bannedUpdateTypes) {
			return request({
				url: '/admin/center/oauth2/double-check/banned',
				method: 'post',
				data,
			});
		}
	};
}
