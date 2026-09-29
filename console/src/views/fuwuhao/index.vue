<template>
	<div class="container" v-if="show">
		<div class="logo-warp">
			<img :src="logo" class="logo" />
			<div class="logo-warp-title">{{ fuwuhaoName }}</div>
		</div>
		<div class="btn">
			<el-button type="primary" size="large" round @click="handleLogin">立即登录</el-button>
		</div>
		<div class="login-action">
			<el-checkbox v-model="isAgree">我已阅读并同意 </el-checkbox>
			<div @click="clickAgreement" class="admin-agreement">用户协议</div>
		</div>
	</div>
	<agreementDialog ref="dialogRef" />
	<silenceSliderDialog ref="silenceSliderRef" />
	<el-dialog v-model="showWeixin" title="绑定手机号" width="500px" :modal="false">
		<Mobile :isBindWeixin="true" :withLogin="true" @success="onSuccess"></Mobile>
	</el-dialog>
</template>

<script setup lang="ts">
import { defineAsyncComponent, onMounted, ref, shallowRef } from 'vue';
import { useFuwuhaoApi } from '/@/api/fuwuhao';
import { NextLoading } from '/@/utils/loading';
import { useLoginSignIn } from '/@/hooks/useLoginSignIn';
import { useRoute, useRouter } from 'vue-router';
import useFile from '/@/hooks/useFile';
import { ElMessage } from 'element-plus';
import { useoauth2Api } from '/@/api/oauth2';
import { encodeSearchParams } from '/@/utils/query';
import { isWeiXin } from '/@/utils/weixin';
import { Local, Session } from '/@/utils/storage';
import Cookies from 'js-cookie';
import { useWeixinApi } from '/@/api/weixin/index';
const agreementDialog = defineAsyncComponent(() => import('/@/views/login/component/agreementDialog.vue')); //管理员协议组件
const silenceSliderDialog = defineAsyncComponent(() => import('/@/views/login/component/silenceSliderDialog.vue'));
const Mobile = defineAsyncComponent(() => import('/@/views/register/component/mobile.vue'));
const show = ref(false);
const fuwuhaoName = import.meta.env.VITE_FUWUHAO_NAME;
const logo = useFile().getFile('logo');
const showWeixin = ref(false);
const onSuccess = () => {
	showWeixin.value = false;
};
const { onSignIn } = useLoginSignIn();
const appID = ref<string>('');

const route = useRoute();
const router = useRouter();

const isAgree = ref<boolean>(false);
const dialogRef = ref();
const clickAgreement = () => {
	dialogRef.value.openDialog();
};

const getAppid = async () => {
	await useFuwuhaoApi()
		.appid()
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				appID.value = res.data.appID;
			}
		});
};
const handleLogin = async () => {
	if (!isAgree.value) {
		ElMessage.warning('请先同意用户协议');
		return;
	}

	if (!appID.value) {
		ElMessage.error('系统繁忙，请稍后');
		return;
	}
	redirect();
};

const silenceSliderRef = shallowRef();

const handleLoginSuccess = async () => {
	if (!route.query.code) {
		show.value = true;
		return;
	}
	sessionStorage.setItem('weixincode', route.query.code as string);
	useFuwuhaoApi()
		.checkFuwuhao({ code: route.query.code as string, type: 'Auto' })
		.then(async (res: any) => {
			if (res.code === 'SUCCESS') {
				if (res.data.type === 'Login2FA') {
					if (Cookies.get('passToken')) {
						silenceSliderRef.value.directlogin(res.data.token);
						return;
					}
					silenceSliderRef.value.openDialog(res.data.token);
				} else if (res.data.type === 'Fuwuhao') {
					sessionStorage.setItem('fuwuhaoToken', res.data.token);
					showWeixin.value = true;
				} else {
					let signRes = await onSignIn(res.data.token, res.data.type, res.data.subType, false);
					if (signRes) {
            sessionStorage.setItem("fuwuhao-login", "true")
						if (route.query?.type === 'fromOauth2' || route.query?.type === 'fromApplication') {
							let domain_uid = route.query?.domain;
							let redirect_uri = route.query?.redirect_uri;
							let next_redirect_uri = route.query?.next_redirect_uri;
							let params = route.query?.params;
							let isLoginToken = Number(route.query?.isLoginToken) == 1;

							let res = await useoauth2Api()
								.oauth2(
									{ domainUID: domain_uid },
									async () => {
										ElMessage.error('授权登录失败');
										Session.clear();
										Local.clear();
									},
									isLoginToken
								)
								.then(async (res: any): Promise<boolean> => {
									if (res.code === 'SUCCESS') {
										let queryString = encodeSearchParams({
											token: res.data.token,
											subToken: res.data.subToken,
											params: params,
											redirect: next_redirect_uri,
										});

										// 不用判断fromApplication，均直接打开
										window.location.href = `${redirect_uri}?${queryString}`;
										return true;
									} else if (res.subCode === 'NOT_OPEN_WEBSITE') {
										// 不用判断fromApplication，均直接打开
										await router.push({
											path: '/oauth2/open',
											query: {
												type: route.query?.type,
												params: params,
												next_redirect_uri: next_redirect_uri,
												redirect_uri: redirect_uri,
												domain: domain_uid,
												isLoginToken: isLoginToken ? 1 : 0,
											},
										}); // 去登录页
										return true;
									} else {
										return false;
									}
								});
							if (!res) {
								ElMessage.error('授权登录失败');
								Session.clear();
								Local.clear();
							}
						} else if (route.query?.redirect) {
							let params = route.query?.params as string;
							let querys = {};
							if (params) {
								querys = JSON.parse(params);
							}

							await router.replace({
								path: <string>route.query?.redirect,
								query: querys,
							});
						} else {
							await router.replace('/home');
						}
					} else {
						ElMessage.error('登录失败');
					}
				}
				NextLoading.done();
			} else {
				show.value = true;
			}
		});
};
const useWeixinApiCollect = useWeixinApi();
const setWeixinCode = async () => {
	sessionStorage.setItem('weixinStatus', '1');
	const code = sessionStorage.getItem('weixincode');
	if (code) {
		sessionStorage.setItem('weixincode', route.query.code as string);
		await useWeixinApiCollect.weixinAccessToken({ code: code, type: 'Wechat' }).then(async (res: any) => {
			if (res.code === 'SUCCESS') {
				await userBindWeixin(res.data.token);
			}
		});
	}
};
const userBindWeixin = async (token: string) => {
	await useWeixinApiCollect.bindWeixin({ wechatToken: token, isDelete: false }).then((res: any) => {
		if (res.code === 'SUCCESS') {
			ElMessage.success('绑定微信成功');
			sessionStorage.setItem('weixinStatus', '2');
		}
	});
};
const redirect = () => {
	const url = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appID.value}&redirect_uri=${encodeURIComponent(
		window.location.href
	)}&response_type=code&scope=snsapi_userinfo&state=1&connect_redirect=1#wechat_redirect`;
	window.location.replace(url);
};

onMounted(async () => {
	if (!isWeiXin()) {
		await router.replace({
			path: '/login',
			query: route.query,
		});
	} else {
		await getAppid();
		await handleLoginSuccess();
	}

	NextLoading.done();
});
</script>

<style lang="scss" scoped>
.container {
	width: 100%;
	height: 100vh;

	.logo-warp {
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;

		img {
			margin-top: 80px;
			width: 100px;
			height: auto;
		}
		&-title {
			font-weight: bold;
			font-size: 20px;
		}
	}
	.btn {
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-top: 90px;
		font-size: 22px;
	}
}
.login-action {
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	position: relative;
	z-index: 2002;
	.admin-agreement {
		color: var(--el-color-primary);
		margin-left: 4px;
		margin-right: 30px;
		text-decoration: underline;
		cursor: pointer;
		margin-top: -2px;
	}
}
</style>
