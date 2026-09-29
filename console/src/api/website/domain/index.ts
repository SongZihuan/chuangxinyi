import request from '/@/utils/request';
import type { createType } from '/@/api/website/ip/types';
export function useDomainnApi() {
	return {
		//新增站点域名
		createdomain: (data: createType) => {
			return request({
				url: '/admin/website/allow-website/domain/create',
				method: 'post',
				data,
			});
		},
		//删除站点域名
		deldomain: (data: { id: number }) => {
			return request({
				url: '/admin/website/allow-website/domain/delete',
				method: 'post',
				data,
			});
		},
	};
}
