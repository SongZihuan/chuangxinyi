import request from '/@/utils/request';
import type { listParamType } from './types';

export function useDefraytApi() {
	return {
		defrayList(params: pageTypes) {
			return request({
				url: '/user/center/allow-website/defray/list',
				method: 'get',
				params,
			});
		},
		defrayInfo(params: { token: string }) {
			return request({
				url: '/user/center/allow-website/defray/info',
				method: 'get',
				params,
			});
		},
		// 购买人 owner
		defrayOwner(params:any) {
			return request({
				url: '/user/center/allow-website/defray/owner/list',
				method: 'get',
				params,
			});
		},
		admindefrayList(params: pageTypes) {
			return request({
				url: '/admin/finance/defray/list',
				method: 'get',
				params,
			});
		},
		//管理员获取用户线下充值请求
		adminPayList: (params: listParamType) => {
			return request({
				url: '/admin/finance/pay/list',
				method: 'get',
				params,
			});
		},
		//管理员处理退款
		adminRefund: (data: { tradeID: string; success: boolean; phoneToken: string }) => {
			return request({
				url: '/admin/finance/pay/refund/process',
				method: 'post',
				data,
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
			});
		},
		adminRefundInside: (data: { tradeID: string; phoneToken: string }) => {
			return request({
				url: '/admin/finance/pay/refund/inside/process',
				method: 'post',
				data,
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
			});
		},
		//管理员处理消费退款
		adminDefrayRefund: (data: { tradeID: string; phoneToken: string }) => {
			return request({
				url: '/admin/finance/defray/return/process',
				method: 'post',
				data: { tradeID: data.tradeID },
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
			});
		},
		//管理员获取用户订单信息
		adminDefrayInfo: (params: { id: string }) => {
			return request({
				url: '/admin/finance/defray/info',
				method: 'get',
				params,
			});
		},
		//管理员处理自充值
		adminpayProcess: (data: { id: string; get: number; phoneToken: string, success: boolean }) => {
			return request({
				url: '/admin/finance/pay/self/process',
				method: 'post',
				data,
				headers: {
					'X-Phone-Token': data.phoneToken,
				},
			});
		},
		//管理员获取用户充值记录
		adminPayHistory: (params: any) => {
			return request({
				url: '/admin/finance/pay/list',
				method: 'get',
				params,
			});
		},
		//管理员获取充值信息
		adminPayInfo: (params: { id: string }) => {
			return request({
				url: '/admin/finance/pay/info',
				method: 'get',
				params,
			});
		},
		//增加余额
		addblance: (data: any) => {
			return request({
				url: '/admin/finance/pay/admin/add',
				method: 'post',
				data,
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
			});
		},
		//减少余额
		adddefary: (data: any) => {
			return request({
				url: '/admin/finance/defray/admin/add',
				method: 'post',
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
				data,
			});
		},
		// 获取用户发票信息
		addbilled: (data: any) => {
			return request({
				url: '/admin/finance/invoice/billed/add',
				method: 'post',
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
				data,
			});
		},
		// 获取用户发票信息
		subBilled: (data: any) => {
			return request({
				url: '/admin/finance/invoice/billed/sub',
				method: 'post',
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
				data,
			});
		},
	};
}
