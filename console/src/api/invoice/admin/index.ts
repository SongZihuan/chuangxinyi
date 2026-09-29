import request from '/@/utils/request';
import { adminProcessInvoiceTypes, invoiceAdminTypes } from './types';

export function useInvoiceAdminApi() {
	return {
		// 获取发票列表
		getInvoiceList(params: invoiceAdminTypes) {
			return request({
				url: '/admin/finance/invoice/list',
				method: 'get',
				params,
			});
		},
		processInvoice(data: adminProcessInvoiceTypes) {
			return request({
				url: '/admin/finance/invoice/process',
				method: 'post',
				data,
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
			});
		},
		getInvoiceInfo(params: { id: string }) {
			return request({
				url: '/admin/finance/invoice/info',
				method: 'get',
				params,
			});
		},
	};
}
