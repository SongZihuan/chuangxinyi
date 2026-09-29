import request from '/@/utils/request';
import { passwordTypes } from '/@/api/register/types';

export function userApi() {
	return {
		getuserList: (params?: object) => {
			if (params == null) {
				params = {
					page: 1,
					pagesize: 10000,
				};
			}
			return request({
				url: '/admin/user/list',
				method: 'get',
				params,
			});
		},
		//绑定角色(管理员)
		bindRole: (data: any) => {
			return request({
				url: '/admin/role/change',
				method: 'post',
				data,
			});
		},

		userAdd: (data?: object) => {
			return request({
				url: '/system/user/add',
				method: 'post',
				data,
			});
		},

		userUpdate: (data?: object) => {
			return request({
				url: '/system/user/update',
				method: 'post',
				data,
			});
		},
		userChangeStatus: (data?: object) => {
			return request({
				url: '/system/user/ban',
				method: 'post',
				data,
			});
		},
		userDelAction: (data?: object) => {
			return request({
				url: '/system/user/delete',
				method: 'post',
				data,
			});
		},
		// 获取用户主页信息
		getUserHomePage: (params: { id?: number; uid?: string }) => {
			return request({
				url: '/admin/user/homepage',
				method: 'get',
				params,
			});
		},
		// 获取用户主页
		getUserData: (params: { id?: number; uid?: string }) => {
			return request({
				url: '/admin/user/data',
				method: 'get',
				params,
			});
		},
		// 更新用户名
		updateUserName: (data: object) => {
			return request({
				url: '/admin/user/username/update',
				method: 'post',
				data,
			});
		},
		// 更新密码
		updateUserPassword: (data: passwordTypes) => {
			return request({
				url: '/admin/user/password/update',
				method: 'post',
				data,
			});
		},
		// 更新邮箱
		updateUserEmail: (data: object) => {
			return request({
				url: '/admin/user/email/update',
				method: 'post',
				data,
			});
		},
		// 更新手机号
		updateUserPhone: (data: object) => {
			return request({
				url: '/admin/user/phone/update',
				method: 'post',
				data,
			});
		},
		// 解绑微信
		unbindWechat: (data: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/wechat/delete',
				method: 'post',
				data,
			});
		},
		// 解绑2FA
		unbind2FA: (data: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/secondfa/delete',
				method: 'post',
				data,
			});
		},
		// 更新用户状态
		updateUserStatus: (data: { uid?: string; id?: string; status?: string; phoneToken: string }) => {
			return request({
				url: '/admin/user/double-check/status/update',
				method: 'post',
				data,
				headers: {
					'X-Phone-Token': data.phoneToken,
					'X-RunMode': 'release',
				},
			});
		},
		// 获取实名认证信息
		getUserRealName: (params: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/info',
				method: 'get',
				params,
			});
		},
		// 获取homepage
		getUserHomePageData: (params: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/homepage',
				method: 'get',
				params,
			});
		},
		// 获取data
		getUserDataData: (params: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/data',
				method: 'get',
				params,
			});
		},
		// 获取用户资产
		getUserFinance: (params: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/finance',
				method: 'get',
				params,
			});
		},
		// 获取在线token
		getUserOnlineToken: (params: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/token/all',
				method: 'get',
				params,
			});
		},
		// 踢下线
		kickUserOnlineToken: (data: { uid?: string; id?: string; token?: string }) => {
			return request({
				url: '/admin/user/token/delete',
				method: 'post',
				data,
			});
		},
		// 根据外站清除在线token
		deleteUserOnlineTokenBySite: (data: { uid?: string; id?: string; webID?: string }) => {
			return request({
				url: '/admin/user/token/website/delete/all',
				method: 'post',
				data,
			});
		},
		// 根据外站清除在线子用户token
		deleteUserOnlineSubUserBySite: (data: { uid?: string; id?: string; }) => {
			return request({
				url: '/admin/user/token/son/delete/all',
				method: 'post',
				data,
			});
		},
		// 获取在线子用户
		getUserOnlineSubUser: (params: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/token/son/all',
				method: 'get',
				params,
			});
		},
		// 取消授权
		deleteUserOauth2:(data:{token:string})=>{
			return request({
				url: '/admin/user/oauth2/delete',
				method: 'post',
				data,
			});
		},
		// 删除所有用户在线token
		deleteUserOnlineToken: (data: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/token/delete/all',
				method: 'post',
				data,
			});
		},
		// 删除所有用户在线子用户
		deleteUserOnlineSubUser: (data: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/token/son/delete/all',
				method: 'post',
				data,
			});
		},
		// 删除父用户在线子用户
		deleteUserOnlineFatherUser: (data: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/token/father/delete/all',
				method: 'post',
				data,
			})
		},
		// 获取授权记录
		getUserOauth2Record: (params: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/oauth2/all',
				method: 'get',
				params,
			});
		},
		// 获取授权域名列表
		getUserOauth2DomainList: (data: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/oauth2/delete',
				method: 'post',
				data,
			});
		},
		// 清除所有授权
		deleteUserOauth2All: (data: { uid?: string; id?: string }) => {
			return request({
				url: '/admin/user/oauth2/delete/all',
				method: 'post',
				data,
			});
		},
		// 获取授权记录
		getUserOauth2RecordList: (params: any) => {
			return request({
				url: '/admin/user/oauth2/record/list',
				method: 'get',
				params,
			});
		},
		// 获取封禁记录
		getUserBanRecord: (params: any) => {
			return request({
				url: '/admin/user/oauth2/banned/list',
				method: 'get',
				params,
			});
		},
		// 更新用户封禁状态
		updateUserBanStatus: (data: any) => {
			return request({
				url: '/admin/user/oauth2/banned',
				method: 'post',
				data,
			})
		},
		// 更新登录控制
		updateLoginController: (data: any) => {
			return request({
				url: '/admin/user/loginctrl/update',
				method: 'post',
				data,
			});
		}
	};
}
