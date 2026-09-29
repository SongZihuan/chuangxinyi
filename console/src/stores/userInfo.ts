import { defineStore } from 'pinia';
import { Session, Local } from '/@/utils/storage';
import { useSocketStore } from '/@/stores/webstocket';
import { MessageTypeEnum } from '/@/data/enum';
import {initBackEndControlRoutes} from "/@/router/backEnd";

/**
 * 用户信息
 * @methods setUserInfos 设置用户信息
 */
export const useUserInfo = defineStore('userInfo', {
	state: (): UserInfosState => ({
		userInfos: {
			role: {}, //用户角色列表
			authBtnList: [], //用户按钮权限列表
			XToken: '', //用户验证token
			bindweixinStatus: false, //用户是否绑定微信成功
			userData: {},
		},
		userType: {
			subType: '',
			type: '',
		},
	}),
	actions: {
		async setUserInfos() {
			// 存储用户信息到浏览器缓存
			if (Session.get('userInfo')) {
				this.userInfos = Session.get('userInfo');
			} else if (Local.get('userInfo')) {
				this.userInfos = Local.get('userInfo');
				this.userType = Local.get('userType').subType;
			} else {
				this.userInfos = {
					role: {}, //用户角色信息
					authBtnList: [], //用户按钮权限列表
					XToken: '', //用户验证token
					bindweixinStatus: false, //用户是否绑定微信成功
					userData: {},
				};
			}
		},
		//设置用户信息
		async setUserData(data: object) {
			this.userInfos.userData = data;
		},
		async setxtoken(data: string) {
			if (data) {
				if (!useSocketStore().sendTokenStatus) {
					useSocketStore().sendMessage(MessageTypeEnum.TOKEN, {
						data: JSON.stringify({
							Token: Session.get('token'),
						}),
					});
				}
			}
			this.userInfos.XToken = data;
			// checkToken();
		},
		async setUserType(data: { subType: string; type: string }) {
			this.userType.subType = data.subType;
			Local.set('userType', data);
			Session.set('userType', data);
		},
		async setRoles(data: string[]) {
			this.userInfos.role = data;
			Local.set('userInfo', this.userInfos);
			Session.set('userInfo', this.userInfos);
			await initBackEndControlRoutes();
		},
		// 钱包
		async setWalletInfos(data: object) {
			this.userInfos.balance = data;
			Local.set('userInfo', this.userInfos);
			Session.set('userInfo', this.userInfos);
		},
		async setPhone(data: string) {
			this.userInfos.phone = data;
		},
		async setWeixinStatus(data: boolean) {
			this.userInfos.bindweixinStatus = data;
		},
	},
});
