import request from '/@/utils/request';
import type { addressUpdateypes } from './types';
export function useUserApi() {
	return {
		//获取用户信息
		userInfo: () => {
			return request({
				url: '/user/center/info',
				method: 'get',
			});
		},
		//更新用户头像
		avatarUpdate: (data: { header: any; isDelete: boolean }) => {
			return request({
				url: '/user/center/header/update',
				method: 'post',
				data,
				headers: {
					'Content-Type': 'multipart/form-data',
				},
			});
		},
		//根据用户id获取用户头像
		userAvatar: (data: { id: string }) => {
			return request({
				url: '/user/center/header/user',
				method: 'get',
				data,
			});
		},
		//更新手机号
		updatePhone: (data: { phoneToken: string }) => {
			return request({
				url: '/user/center/phone/update',
				method: 'post',
				data,
			});
		},
		//更新邮箱
		updateEmail: (data: { emailToken: string; isDelete: boolean }) => {
			return request({
				url: '/user/center/email/update',
				method: 'post',
				data,
			});
		},
		//更新2fa
		updateTopt: (data: { emailToken: string; isDelete: boolean }) => {
			return request({
				url: '/admin/auth/2fa/update',
				method: 'post',
				data,
			});
		},
		//注销账号
		userLogOff: (data: { phoneToken: string }) => {
			return request({
				url: '/user/center/double-check/delete',
				method: 'post',
				headers: {
					'X-Phone-Token': data.phoneToken,
				},
			});
		},
		//更新个人信息
		addressUpdate: (data: addressUpdateypes) => {
			return request({
				url: '/user/center/address/update ',
				method: 'post',
				data,
			});
		},
		// 判断token是否过期
		checkToken: () => {
			return request({
				url: '/ping/user',
				method: 'get',
			});
		},
	};
}
