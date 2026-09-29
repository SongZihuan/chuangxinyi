import request from '/@/utils/request';
export function useCouponApi() {
	return {
		//获取优惠券列表
		couponList: (params?: any) => {
			return request({
				url: '/public/discount/list',
				method: 'get',
				params,
			});
		},
		//分配用户优惠包
		allocationUser: (data: any) => {
			return request({
				url: '/user/center/discount/join',
				method: 'post',
				data,
			});
		},

		couponDict: (params?: any) => {
			return request({
				url: '/user/center/coupons/list',
				method: 'get',
				params,
			});
		},
	};
}
