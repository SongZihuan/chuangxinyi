import request from '/@/utils/request';
import type { checkFuwuhao } from './types';
export function useFuwuhaoApi() {
	return {
		//获取appid
		appid: () => {
			return request({
				url: '/public/checker-before/fuwuhao/get',
				method: 'get',
			});
		},
		//检验公众号登录
		checkFuwuhao: (data: checkFuwuhao) => {
			return request({
				url: '/public/checker/fuwuhao',
				method: 'post',
				data,
			});
		},
	};
}
