import { Local, Session } from '/@/utils/storage';
import { initBackEndControlRoutes } from '/@/router/backEnd';
import { NextLoading } from '/@/utils/loading';
import { useUserApi } from '/@/api/user/user';
import { useRouter } from 'vue-router';
import { computed } from 'vue';
import { formatAxis } from '/@/utils/formatTime';
import { useUserInfo } from '/@/stores/userInfo';
import { useSubUserApi } from '/@/api/user/subuser';
import { ElMessage } from 'element-plus';
import { useWeixinApi } from '/@/api/weixin/index';
export function useLoginSignSub() {
	const useUserApiCollect = useUserApi();
	const stores = useUserInfo();
	const currentTime = computed(() => {
		return formatAxis(new Date());
	});

	const onSignSub = async (userID: string, quite: boolean = false): Promise<boolean> => {
		let logintoken = Session.get('login-token');
		if (!logintoken) {
			return false;
		}
		return await useSubUserApi()
			.getToken({ id: userID })
			.then(async (res: any): Promise<boolean> => {
				if (res.code === 'SUCCESS') {
					Session.clear();
					Session.set('token', res.data.token);
					Session.set('login-token', logintoken);
					Session.set('token-id', userID);

					await stores.setUserType({ type: res.data.type, subType: res.data.subType });
					await stores.setxtoken(res.data.token);

					return await afterSignSub(quite);
				} else {
					return false;
				}
			});
	};

	const regetToken = async (): Promise<boolean> => {
		let tokenID = Session.get('token-id');
		if (!tokenID) {
			return false;
		}

		return onSignSub(tokenID, true);
	};

	const afterSignSub = async (quite: boolean = false): Promise<boolean> => {
		let ok = await useUserApiCollect.userInfo().then(async (res: any): Promise<boolean> => {
			if (res.code === 'SUCCESS') {
				await stores.setUserInfos();

				Session.set('userData', res.data);
				Session.set('userInfo', res.data);
				Local.set('userInfo', res.data);

				if (!res.data.role.menus || res.data.role.menus.length === 0) {
					if (!quite) {
						ElMessage.warning('抱歉，您没有登录权限');
					}
					return false;
				}
				return true;
			}
			return false;
		});
		if (!ok) {
			return false;
		}

		const hasPower = await initBackEndControlRoutes();
		return signInSubSuccess(!hasPower, quite);
	};

	const signInSubSuccess = (isNoPower: boolean | undefined, quite: boolean = false): boolean => {
		if (isNoPower) {
			if (!quite) {
				ElMessage.warning('抱歉，您没有登录权限');
			}
			Session.clear();
			Local.clear();
			return false;
		} else {
			if (!quite) {
				let currentTimeInfo = currentTime.value;
				const signInText = '欢迎回来！';
				ElMessage.success(`${currentTimeInfo}，${signInText}`);
			}
			return true;
		}
	};

	return { onSignSub, afterSignSub, regetToken };
}

export function useLoginSignIn() {
	const useUserApiCollect = useUserApi();
	const stores = useUserInfo();
	const router = useRouter();
	const currentTime = computed(() => {
		return formatAxis(new Date());
	});

	const resetSignIn = async (isReload: boolean = false): Promise<boolean> => {
		let logintoken = Session.get('login-token');
		let nowtoken = Session.get('token');
		if (!logintoken || nowtoken === logintoken.token) {
			return false;
		}

		Session.clear();
		Session.set('token', logintoken.token);
		Session.set('login-token', logintoken);
		Session.set('token-id', '');

		await stores.setUserType({ type: logintoken.tokenType, subType: logintoken.tokenSubType });
		await stores.setxtoken(logintoken.token);

		return await afterSignIn(isReload);
	};

	const onSignIn = async (
		token: string,
		tokenType: string,
		tokenSubType: string,
		isReload: boolean = false,
		isbindWeixin: boolean = false
	): Promise<boolean> => {
		// Session.clear();
		Session.set('token', token);
		Session.set('login-token', {
			token: token,
			tokenType: tokenType,
			tokenSubType: tokenSubType,
		});
		Session.set('token-id', '');

		await stores.setUserType({ type: tokenType, subType: tokenSubType });
		await stores.setxtoken(token);
		if (isbindWeixin) {
			// await setWeixinCode();
			const weixinToken = sessionStorage.getItem('weixinToken') as string ||  sessionStorage.getItem('fuwuhaoToken') as string
			userBindWeixin(weixinToken)
		}
		return await afterSignIn(isReload);
	};
	const useWeixinApiCollect = useWeixinApi();
	const userBindWeixin = async (token: string) => {
		await useWeixinApiCollect.bindWeixin({ wechatToken: token, isDelete: false }).then((res: any) => {
			if (res.code === 'SUCCESS') {
				ElMessage.success('绑定微信成功');
				sessionStorage.setItem('weixinStatus', '2');
			}
		});
	};

	const afterSignIn = async (isReload: boolean = false): Promise<boolean> => {
		await useUserApiCollect.userInfo().then(async (res: any) => {
			if (res.code === 'SUCCESS') {
				Session.set('userData', res.data);
				Session.set('userInfo', res.data);
				Local.set('userInfo', res.data);

				await stores.setUserInfos();

				if (!res.data.role.menus || res.data.role.menus.length === 0) {
					ElMessage.warning('抱歉，您没有登录权限');
					return false;
				}
				if (isReload) {
					// 重定向到首页
					await router.push({ path: '/home' });
					location.reload();
				}
			}
		});

		const hasPower = await initBackEndControlRoutes();
		return signInSuccess(!hasPower);
	};

	// 登录成功后的跳转
	const signInSuccess = (isNoPower: boolean | undefined): boolean => {
		if (isNoPower) {
			ElMessage.warning('抱歉，您没有登录权限');
			Session.clear();
			Local.clear();
			return false;
		} else {
			let currentTimeInfo = currentTime.value;
			const signInText = '欢迎回来！';
			ElMessage.success(`${currentTimeInfo}，${signInText}`);
			NextLoading.start();
			return true;
		}
	};

	const isLoginUser = function (): boolean {
		let tokenID = Session.get('token-id');
		return !tokenID;
	};

	return { onSignIn, afterSignIn, resetSignIn, isLoginUser };
}
