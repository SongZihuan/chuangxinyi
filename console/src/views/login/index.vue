<template>
	<div class="flex login-container">
		<div class="login-left">
			<div class="login-left-img">
				<img :src="loginMain" />
			</div>
		</div>
		<div class="flex login-right">
			<div class="login-right-warp flex-margin">
				<div class="flex login-right-wx" v-if="!isMobile()">
					<div class="code">
						<Scan @signBack="signIn" @totpBack="totpBack" />
					</div>
				</div>
				<div class="w-full">
					<span class="login-right-warp-one"></span>
					<span class="login-right-warp-two"></span>
					<div class="login-right-warp-mian">
						<div class="logo-warp">
							<img :src="logo" class="logo" />
						</div>
						<div class="login-right-warp-main-title" v-if="!state.isScan">{{ !state.isScan ? '登录' : null }}</div>
						<div class="login-right-warp-main-form">
							<div v-if="!state.isScan">
								<el-tabs v-model="state.tabsActiveName">
									<el-tab-pane label="手机和邮箱登录" name="mobile">
										<Mobile v-if="state.tabsActiveName == 'mobile'" @signBack="signIn" :isAgree="isAgree" @totpBack="totpBack" />
									</el-tab-pane>
									<el-tab-pane label="账号密码登录" name="account">
										<Account v-if="state.tabsActiveName == 'account'" :isAgree="isAgree" @signBack="signIn" @totpBack="totpBack" ref="countRef" />
									</el-tab-pane>
								</el-tabs>
								<div class="login-action">
									<el-checkbox v-model="isAgree">我已阅读并同意 </el-checkbox>
									<div @click="clickAgreement" class="admin-agreement">用户协议</div>
								</div>
							</div>
							<Scan v-if="state.isScan" @signBack="signIn" />
							<div class="login-content-main-sacn" v-if="isMobile()" @click="state.isScan = !state.isScan">
								<SvgIcon name="my-pc" :size="50" v-if="state.isScan"></SvgIcon>
								<SvgIcon name="my-qrcode " :size="50" v-else color="#056de8"></SvgIcon>

								<!-- <i class="iconfont" :class="state.isScan ? 'icon-diannao1' : 'icon-barcode-qr'"></i> -->
								<div class="login-content-main-sacn-delta"></div>
							</div>
							<div class="register" v-if="!state.isScan">
								<div>还没有账号?</div>
								<div @click="goRegister">立即注册</div>
							</div>
						</div>
						<div class="page-foot">
							<div v-if="!state.isScan && footInfo.copyright" class="footer-text">
								版权所有 © {{ dayjs().format('YYYY') }} {{ footInfo.copyright }}
							</div>
							<div class="footer-text">
								<a href="https://beian.miit.gov.cn/" target="_blank" v-if="footInfo.icp1">{{ footInfo.icp1 }}</a>
								<a :href="icp2Url" target="_blank" v-if="footInfo.icp2">{{ footInfo.icp2 }}</a>
							</div>
							<div class="footer-text" v-if="footInfo.gongan">
								<img :src="gongan" alt="gongan" class="" />
								<a href="https://beian.mps.gov.cn/#/query/webSearch" target="_blank">{{ footInfo.gongan }}</a>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<!-- 管理员协议 -->
		<agreementDialog ref="dialogRef" />
		<!-- 二次验证 -->
		<silenceSliderDialog ref="silenceSliderRef" />
	</div>
</template>

<script setup lang="ts" name="loginIndex">
import { defineAsyncComponent, onMounted, reactive, ref, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import dayjs from 'dayjs';
import { NextLoading } from '/@/utils/loading';
import { Local, Session } from '/@/utils/storage';
import { useFooterApi } from '/@/api/footer/index';
import Scan from '/@/views/login/component/scan.vue';
import { useUserInfo } from '/@/stores/userInfo';
import { useLoginSignIn } from '/@/hooks/useLoginSignIn';
import useFile from '/@/hooks/useFile';
import gongan from '/@/assets/gongan.png';
import commonFunction from '/@/utils/commonFunction';
import Cookies from 'js-cookie';
import { ElMessage } from 'element-plus';
import { useoauth2Api } from '/@/api/oauth2';
import { encodeSearchParams } from '/@/utils/query';
import { isWeiXin } from '/@/utils/weixin';
const { isMobile } = commonFunction();
const logo = useFile().getFile('logo');
const loginMain = useFile().getFile('login-bg');
const useFooterApiCollect = useFooterApi();
const stores = useUserInfo();
const countRef = ref();
const footInfo = ref({
	copyright: '',
	icp1: '',
	gongan: '',
	icp2: '',
});
const icp2Url = ref<string>();
const silenceSliderRef = ref();
// 引入组件
const Account = defineAsyncComponent(() => import('/@/views/login/component/account.vue'));
const Mobile = defineAsyncComponent(() => import('/@/views/login/component/mobile.vue'));
const Email = defineAsyncComponent(() => import('/@/views/login/component/email.vue'));
const agreementDialog = defineAsyncComponent(() => import('/@/views/login/component/agreementDialog.vue')); //管理员协议组件
const silenceSliderDialog = defineAsyncComponent(() => import('/@/views/login/component/silenceSliderDialog.vue'));
const route = useRoute();

const state = reactive({
	tabsActiveName: 'mobile',
	isScan: false,
});
const router = useRouter();
const isAgree = ref<boolean>(false);
const dialogRef = ref();
const clickAgreement = () => {
	dialogRef.value.openDialog();
};
const { onSignIn } = useLoginSignIn();
const goRegister = () => {
	router.push({ path: '/register', query: route.query });
};

//2FA第一次次验证回调
const totpBack = (token: string) => {
	if (Cookies.get('passToken')) {
		silenceSliderRef.value.directlogin(token);
		return;
	}
	silenceSliderRef.value.openDialog(token);
};

const getFoot = () => {
	useFooterApiCollect.footer().then((res: any) => {
		if (res.code === 'SUCCESS') {
			footInfo.value = res.data;
			icp2Url.value = 'http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=' + res.data.icp2;
		}
	});
};

const signIn = async (token: string, type: string, subType: string) => {
	let res = await onSignIn(token, type, subType);
	if (res) {
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

						if (route.query?.type === 'fromApplication') {
							window.open(`${redirect_uri}?${queryString}`);
							await router.replace('/home');
						} else {
							window.location.href = `${redirect_uri}?${queryString}`;
						}
						return true;
					} else if (res.subCode === 'NOT_OPEN_WEBSITE') {
						if (route.query?.type === 'fromApplication') {
							window.open(
								router.resolve({
									path: '/oauth2/open',
									query: {
										type: route.query?.type,
										params: params,
										next_redirect_uri: next_redirect_uri,
										redirect_uri: redirect_uri,
										domain: domain_uid,
										isLoginToken: isLoginToken ? 1 : 0,
									},
								}).href
							); // 新窗口打开
						} else {
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
						}
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
};

onUnmounted(async () => {
	NextLoading.done();
});

// 页面加载时
onMounted(async () => {
	if (isWeiXin()) {
		await router.replace({
			path: '/fuwuhao',
			query: route.query,
		});
	} else {
		getFoot();
		//注册手机号后点立即登录
		if (route.query?.type == 'fromRegister') {
			let res = await onSignIn(stores.userInfos.XToken, stores.userType.type, stores.userType.subType, false);
			if (res) {
				if (route.query?.redirect) {
					let params = route.query?.params as string;
					let querys = {};
					if (params) {
						querys = JSON.parse(params);
					}

					await router.replace({
						path: route.query?.redirect as string,
						query: querys,
					});
				} else {
					await router.push({ path: '/login', query: route.query });
				}
			}
		}
	}
	NextLoading.done();
});
</script>

<style scoped lang="scss">
.content-box {
	width: 340px !important;
}
.logo-warp {
	width: 100%;
	display: flex;
	justify-content: center;
	.logo {
		width: 100px;
		height: auto;
		margin-top: 30px;
	}
}

.login-action {
	display: flex;
	flex-direction: row;
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
:deep(.el-tabs__active-bar) {
	height: 3px;
}

:deep(input::-webkit-input-placeholder) {
	color: #9ca5ba;
	font-size: 12px;
}

:deep(.el-input__inner) {
	background-color: transparent !important;
}
:deep(.el-tabs__item.is-active) {
	color: #121212;
	font-weight: bolder;
}

.login-container {
	height: 100%;
	background: var(--el-color-white);
	.login-left {
		flex: 1;
		position: relative;
		display: flex;
		background-color: rgba(211, 239, 255, 1);
		align-items: center;
		justify-content: center;
		.content-box {
			padding: 15px;
			border: 1px solid var(--el-color-primary);
			background: #fff;
			border-radius: 4px;
			width: 540px;
			margin-top: 20px;
			.content-title {
				font-size: 18px;
				color: #40485b;
				margin-bottom: 10px;
				text-align: center;
				font-weight: bold;
			}
			.content-item {
				padding: 8px 0px;
				color: #40485b;
				margin-left: 15px;
				line-height: 20px;
			}
			.btn {
				display: flex;
				flex-direction: row-reverse;
			}
		}
		.login-left-logo {
			display: flex;
			align-items: center;
			position: absolute;
			top: 50px;
			left: 80px;
			z-index: 1;
			animation: logoAnimation 0.3s ease;
			img {
				width: 52px;
				height: 52px;
			}
			.login-left-logo-text {
				display: flex;
				flex-direction: column;
				span {
					margin-left: 10px;
					font-size: 28px;
					color: #26a59a;
				}
				.login-left-logo-text-msg {
					font-size: 12px;
					color: #32a99e;
				}
			}
		}
		.login-left-img {
			position: absolute;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
			width: 50%;
			height: 40%;
			img {
				width: 100%;
				height: 100%;
				animation: error-num 0.6s ease;
				object-fit: contain;
			}
		}
		.login-left-waves {
			position: absolute;
			top: 0;
			right: -100px;
		}
	}
	.login-right {
		width: 1000px;
		.login-right-warp {
			border: 1px solid var(--el-color-primary-light-3);
			border-radius: 3px;
			width: 900px;
			min-height: 580px;
			position: relative;
			overflow: hidden;
			background-color: var(--el-color-white);
			display: flex;
			.login-right-wx {
				width: 400px;
				height: 100%;
				min-height: 700px;
				display: flex;
				align-items: center;
				justify-content: center;
				.code {
					width: 400px;
					height: 500px;
					border-right: 1px solid #ebeef5;
					display: flex;
					align-items: center;
					justify-content: center;
				}
			}
			.login-right-warp-one,
			.login-right-warp-two {
				position: absolute;
				display: block;
				width: inherit;
				height: inherit;
				&::before,
				&::after {
					content: '';
					position: absolute;
					z-index: 1;
				}
			}
			.login-right-warp-one {
				&::before {
					filter: hue-rotate(0deg);
					top: 0px;
					left: 0;
					width: 100%;
					height: 3px;
					background: linear-gradient(90deg, transparent, var(--el-color-primary));
					animation: loginLeft 3s linear infinite;
				}
				&::after {
					filter: hue-rotate(60deg);
					top: -100%;
					right: 2px;
					width: 3px;
					height: 100%;
					background: linear-gradient(180deg, transparent, var(--el-color-primary));
					animation: loginTop 3s linear infinite;
					animation-delay: 0.7s;
				}
			}
			.login-right-warp-two {
				&::before {
					filter: hue-rotate(120deg);
					bottom: 2px;
					right: -100%;
					width: 100%;
					height: 3px;
					background: linear-gradient(270deg, transparent, var(--el-color-primary));
					animation: loginRight 3s linear infinite;
					animation-delay: 1.4s;
				}
				&::after {
					filter: hue-rotate(300deg);
					bottom: -100%;
					left: 0px;
					width: 3px;
					height: 100%;
					background: linear-gradient(360deg, transparent, var(--el-color-primary));
					animation: loginBottom 3s linear infinite;
					animation-delay: 2.1s;
				}
			}
			.login-right-warp-mian {
				display: flex;
				flex-direction: column;
				justify-content: center;

				height: 100%;
				.login-right-warp-main-title {
					height: 80px;
					line-height: 80px;
					font-size: 27px;
					text-align: center;
					letter-spacing: 3px;
					animation: logoAnimation 0.3s ease;
					animation-delay: 0.3s;
					color: var(--el-text-color-primary);
				}
				.login-right-warp-main-form {
					flex: 1;
					padding: 0 50px 50px;
					.login-content-main-sacn {
						position: absolute;
						top: 0;
						right: 0;
						width: 50px;
						height: 50px;
						overflow: hidden;
						cursor: pointer;
						transition: all ease 0.3s;
						color: var(--el-color-primary);
						&-delta {
							position: absolute;
							width: 35px;
							height: 70px;
							z-index: 2;
							top: 2px;
							right: 21px;
							background: var(--el-color-white);
							transform: rotate(-45deg);
						}
						&:hover {
							opacity: 1;
							transition: all ease 0.3s;
							color: var(--el-color-primary) !important;
						}
						i {
							width: 47px;
							height: 50px;
							display: inline-block;
							font-size: 48px;
							position: absolute;
							right: 1px;
							top: 0px;
						}
					}
					.register {
						display: flex;
						justify-content: center;
						font-size: 14px;
						margin-top: 20px;
						cursor: pointer;
						position: relative;
						z-index: 2004;
						:nth-child(1) {
							color: #9ca5ba;
						}
						:nth-child(2) {
							color: var(--el-color-primary);
							margin-left: 6px;
							text-decoration: underline;
						}
					}
				}
				.page-foot {
					width: calc(100% - 30px);
					margin-left: 15px;
					display: flex;
					padding-bottom: 40px;
					text-align: center;
					color: #9ca5ba;
					display: flex;
					flex-direction: row;
					justify-content: center;
					cursor: pointer;
					position: relative;
					z-index: 2004;
					display: flex;
					flex-direction: column;
					align-items: center;
					a {
						margin-right: 8px;
						color: #9ca5ba;
					}
					img {
						width: 16px;
						height: 16px;
						margin-right: 4px;
					}
				}
			}
		}
	}
}

.footer-text {
	margin-bottom: 5px;
	margin-top: 5px;
}
</style>
