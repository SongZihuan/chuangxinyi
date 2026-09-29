import request from '/@/utils/request';
import { allocationUserType, adminUserType } from './types';
export function useCouponApi() {
	return {
		//获取优惠券列表
		couponList: (params: any) => {
			return request({
				url: '/admin/discount/list',
				method: 'get',
				params,
			});
		},
		//新增优惠券
		createCoupon: (data: { domain: string }) => {
			return request({
				url: '/admin/discount/create',
				method: 'post',
				data,
			});
		},
		//更新优惠券
		updateCoupon: (data: { domain: string }) => {
			return request({
				url: '/admin/discount/update',
				method: 'post',
				data,
			});
		},
		//删除优惠券
		delCoupon: (data: { id: number }) => {
			return request({
				url: '/admin/discount/delete',
				method: 'post',
				data,
			});
		},
		//分配用户优惠包
		allocationUser: (data: allocationUserType) => {
			return request({
				url: '/admin/discount/join',
				method: 'post',
				data,
			});
		},
		//管理员获取用户优惠列表
		adminGetUserDiscount: (params: adminUserType) => {
			return request({
				url: '/admin/discount/coupons/list',
				method: 'get',
				params,
			});
		},
		// 获取优惠券列表
		getCouponList: (params: any) => {
			return request({
				url: '/admin/discount/coupons/list',
				method: 'get',
				params,
			});
		}
	};
}
