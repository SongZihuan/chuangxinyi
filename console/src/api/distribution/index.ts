import request from '/@/utils/request';
export function useDistributionApi() {
	return {
		//获取分销信息
		distributionList: () => {
			return request({
				url: '/admin/user/distribution/list',
				method: 'post',
			});
		},
		//更新分销
		distributionUpdate: (data: any) => {
			return request({
				url: '/admin/user/distribution/update',
				method: 'post',
				data,
			});
		},
		//删除
		distributionDel: (data: { level: number }) => {
			return request({
				url: '/admin/user/distribution/delete',
				method: 'post',
				data,
			});
		},
	};
}
