import request from '/@/utils/request';
import type { createType } from '/@/api/website/ip/types';
export function useIpApi() {
	return {
		//新增登录站点IP
		createIp: (data: createType) => {
			return request({
				url: '/admin/website/allow-website/ip/create',
				method: 'post',
				data,
			});
		},
		//删除登录站点IP
		delIp: (data: { id: number }) => {
			return request({
				url: '/admin/website/allow-website/ip/delete',
				method: 'post',
				data,
			});
		},
	};
}
