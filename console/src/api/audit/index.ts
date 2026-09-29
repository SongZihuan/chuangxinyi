import request from '/@/utils/request';
import { auditTypes } from '/@/api/audit/types';
export function useAuditApi() {
	return {
		// 用户获取审计列表
		getAuditList(params: auditTypes) {
			return request({
				url: '/user/center/allow-website/audit/list',
				method: 'get',
				params,
			});
		},
		// 管理员获取审计列表
		adminAuditList(params: auditTypes) {
			return request({
				url: '/admin/msg/allow-website/audit/list',
				method: 'get',
				params,
			});
		},
	};
}
