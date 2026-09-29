import request from '/@/utils/request';
export function useFooterApi() {
	return {
		//获取底部信息(备案号等)
		footer: () => {
			return request({
				url: '/public/ui/footer',
				method: 'get',
			});
		},
		editFooter: (data: any) => {
			return request({
				url: '/admin/ui/footer/update',
				method: 'post',
				data,
			});
		},
	};
}
