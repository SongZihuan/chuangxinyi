import request from '/@/utils/request';
export function useAdminWorkOrderApi() {
	return {
		//工单列表
		userWorkOrderList: (params: any) => {
			return request({
				url: '/admin/msg/allow-website/order/list',
				method: 'get',
				params,
			});
		},
		//新增工单
		userWorkOrder: (data: any) => {
			return request({
				url: '/admin/msg/allow-website/order/create',
				method: 'post',
				headers: {
					'Content-Type': 'multipart/form-data',
				},
				data,
			});
		},
		//沟通列表
		communicateList: (params: any) => {
			return request({
				url: '/admin/msg/allow-website/order/communicate/list',
				method: 'get',
				params,
			});
		},
		//获取工单文件
		userOrderFile: (params: any) => {
			return request({
				url: '/admin/msg/allow-website/order/file',
				method: 'get',
				params,
			});
		},
		//回复工单
		userReplyWorkOrder: (data: any) => {
			return request({
				url: '/admin/msg/allow-website/order/reply',
				method: 'post',
				headers: {
					'Content-Type': 'multipart/form-data',
				},
				data,
			});
		},
		userOrderFinsh: (data: { orderID: string ,status:number}) => {
			return request({
				url: '/admin/msg/allow-website/order/finish ',
				method: 'post',
				data,
			});
		},
		//管理员获取工单列表
		adminOrderList: (params: any) => {
			return request({
				url: '/admin/msg/allow-website/order/list',
				method: 'get',
				params,
			});
		},
		//沟通列表(管理员)
		adminCommunicateList: (params: any) => {
			return request({
				url: '/admin/center/order/communicate/list',
				method: 'get',
				params,
			});
		},
		//获取工单文件
		adminOrderFile: (params: any) => {
			return request({
				url: '/admin/msg/allow-website/order/file',
				method: 'get',
				params,
			});
		},
		//完成工单
		adminOrderFinsh: (params: any) => {
			return request({
				url: '/admin/center/order/finish ',
				method: 'get',
				params,
			});
		},
	};
}
