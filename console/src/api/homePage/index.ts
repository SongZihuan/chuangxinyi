import request from '/@/utils/request';
import { pageHomeParamsTypes } from '/@/api/homePage/types';
export function useHomePageApi() {
	return {
		//获取首页信息
		getHomePage: (params: pageHomeParamsTypes) => {
			return request({
				url: '/public/homepage',
				method: 'get',
				params,
			});
		},
		//修改首页信息
		updateHomePage: (data: pageHomeParamsTypes) => {
			return request({
				url: '/user/center/homepage/update',
				method: 'post',
				data,
			});
		},
		// 访客首页信息接口
		getVisitorHomePage: (params: pageHomeParamsTypes) => {
			return request({
				url: '/admin/user/homepage',
				method: 'get',
				params,
			});
		}
	};
}
