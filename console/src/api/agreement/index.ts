import request from '/@/utils/request';
import type { fromType } from './types';
export function useAgreementApi() {
	return {
		agreementList(params: pageTypes) {
			return request({
				url: '/admin/agreement/list',
				method: 'get',
				params,
			});
		},
		//添加协议
		createAgreement: (data: fromType) => {
			return request({
				url: '/admin/agreement/create',
				method: 'post',
				data,
			});
		},
		//编辑协议
		editAgreement: (data: fromType) => {
			return request({
				url: '/admin/agreement/update',
				method: 'post',
				data,
			});
		},
		//删除协议
		delrAgreement: (data: { aid: string }) => {
			return request({
				url: '/admin/agreement/delete',
				method: 'post',
				data,
			});
		},
		//获取协议
		getAgreement: (params: { aid: string }) => {
			return request({
				url: '/public/agreement',
				method: 'get',
				params,
				notJson: true,
			} as any);
		},
	};
}
