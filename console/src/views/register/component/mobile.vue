<template>
	<el-form size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="mobileFormRef">
		<el-form-item class="login-animation1 mb20 " prop="phone">
			<el-input text placeholder="请输入手机号" v-model="ruleForm.phone" clearable autocomplete="off" >
				<template #prefix>
					<el-icon><ele-Phone /></el-icon>
				</template>
			</el-input>
		</el-form-item>
		<el-form-item class="login-animation1" prop="inviteID">
			<el-input text placeholder="（非必填）请输入邀请人ID或手机号" v-model="ruleForm.inviteID" clearable autocomplete="off" >
				<template #prefix>
					<el-icon class="el-input__icon"><ele-User /></el-icon>
				</template>
			</el-input>
		</el-form-item>
		<el-form-item class="login-animation3">
			<el-col :span="24">
				<slider-silence @siderEmit="siderEmit" ref="sliderRef" />
			</el-col>
		</el-form-item>
		<el-form-item class="login-animation2" prop="code" v-if="isSliderCheck">
			<el-col :span="15">
				<el-input text maxlength="6" placeholder="请输入验证码" v-model="ruleForm.code" clearable autocomplete="off" >
					<template #prefix>
						<el-icon class="el-input__icon"><ele-Position /></el-icon>
					</template>
				</el-input>
			</el-col>
			<el-col :span="1"></el-col>
			<el-col :span="8">
				<el-button v-waves class="login-content-code" @click="getCode" :disabled="codeState.isSend"> {{ codeState.codeName }}</el-button>
			</el-col>
		</el-form-item>
		<el-form-item class="login-animation3">
			<el-button round type="primary" v-waves class="login-content-submit" @click="onRegister(mobileFormRef)" v-if="props.showRegisterBtn">
				<span>{{ props.isBindWeixin ?'绑定' :'注册' }}</span>
			</el-button>
		</el-form-item>
		<div class="login-action">
			<el-checkbox v-model="isAgree">我已阅读并同意 </el-checkbox>
			<div @click="clickAgreement" class="admin-agreement">用户协议</div>
		</div>
	</el-form>
	<!-- 管理员协议 -->
	<agreementDialog ref="dialogRef" />
</template>
<script setup lang="ts" name="registerMobile">
import { defineAsyncComponent, reactive, ref } from 'vue';
import sliderSilence from '/@/components/Slider/index.vue';
import { useRegisterApi } from '/@/api/register/index';
import { useBaseApi } from '/@/api/base/index';
import { verifyPhone } from '/@/utils/toolsValidate';
import { type FormRules, type FormInstance, ElMessageBox } from 'element-plus';
import type { checkPhoneCodeTypes, codeTypes } from '/@/api/register/types';
import type { sliderHeadersTypes } from '/@/api/base/types';
import { useUserInfo } from '/@/stores/userInfo';
import { ElMessage } from 'element-plus';

const useRegisterCollect = useRegisterApi();
import { useLoginSignIn } from '/@/hooks/useLoginSignIn';
import { useRoute, useRouter } from 'vue-router';
import { encodeSearchParams } from '/@/utils/query';
import { useoauth2Api } from '/@/api/oauth2';
import { Local, Session } from '/@/utils/storage';
const useBaseApiCollect = useBaseApi();

const agreementDialog = defineAsyncComponent(() => import('/@/views/login/component/agreementDialog.vue')); //管理员协议组件

const route = useRoute();
const router = useRouter();
const isAgree = ref<boolean>(false);
const dialogRef = ref();

const emit = defineEmits(['success']);

const clickAgreement = () => {
	dialogRef.value.openDialog();
};

interface Props {
	showRegisterBtn: boolean; // 显示注册按钮
	isBindWeixin: boolean;
}
const props = withDefaults(defineProps<Props>(), {
	showRegisterBtn: true,
	isBindWeixin: false,
});
// 定义变量内容
const ruleForm = ref<checkPhoneCodeTypes>({
	phone: '',
	code: '',
	inviteID: '',
	type: 'Auto',
});
const mobileFormRef = ref();
const stores = useUserInfo();
const rules = reactive<FormRules>({
	phone: [{ trigger: 'blur', validator: verifyPhone }],
	code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
});
const sliderRef = ref();

//阿里云滑块返回的参数
let sliderHeaders = ref<sliderHeadersTypes>({
	sig: '',
	sessionId: '',
	token: '',
});
const isSliderCheck = ref<boolean>(false);
const siderEmit = (data: any) => {
	sliderHeaders.value = data;
	isSliderCheck.value = true;
};
const { onSignIn } = useLoginSignIn();
const isGetCode = ref<boolean>(false);
const codeState = reactive({
	isSend: false,
	codeName: '获取验证码',
	totalTime: 60, //一般是60
	timer: undefined, //定时器
}) as codeTypes;
const getCode = async () => {
	await useBaseApiCollect
		.sendPhoneCode({ phone: ruleForm.value.phone }, sliderHeaders.value)
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				isGetCode.value = true;
				if (codeState.isSend) return;
				// getCode() // 获取验证码的接口
				codeState.isSend = true;
				codeState.codeName = codeState.totalTime + 's后重新发送';
				codeState.timer = setInterval(() => {
					codeState.totalTime--;
					codeState.codeName = codeState.totalTime + 's后重新发送';
					if (codeState.totalTime < 0) {
						clearInterval(codeState.timer);
						codeState.codeName = '重新发送';
						codeState.totalTime = 60;
						codeState.isSend = false;
						sliderRef.value.resetSider();
					}
				}, 1000);
	
				ElMessage.success('发送手机验证码成功');
			} else {
				sliderRef.value.resetSider();
			}
			sliderRef.value.resetSider();
		})
		.catch(() => {
			sliderRef.value.resetSider();
		});
};
const onRegister = (formEl: FormInstance | undefined) => {
	if (!isAgree.value) {
		ElMessage.warning('请勾选用户协议');
		return;
	}

	if (!isSliderCheck.value) {
		ElMessage.warning('请滑动滑块移动最右侧!');
		return;
	}

	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			await useRegisterCollect.checkPhoneCode(ruleForm.value, sliderHeaders.value).then(async (res: any) => {
				isSliderCheck.value = false;
				sliderRef.value.resetSider();
				if (res.code === 'SUCCESS') {
					stores.setUserType({ type: res.data.type, subType: res.data.subType });
					if (res.data.type === 'UserToken') {
						ElMessageBox.confirm('您的手机号已注册，是否直接登录?', '提示', {
							distinguishCancelAndClose: true,
							confirmButtonText: '直接登录',
							cancelButtonText: '取消',
						})
							.then(async () => {
								emit('success');
								stores.setxtoken(res.data.token);
								localStorage.setItem('xtoken', res.data.token);
								localStorage.setItem('xtokentype', res.data.type);
								localStorage.setItem('xtokensubtype', res.data.subType);
								await isGoLogin();
							})
							.catch(() => {});
					} else if (res.data.type === 'Login2FA') {
						return;
					} else {
						emit('success');
						stores.setxtoken(res.data.token);
						localStorage.setItem('xtoken', res.data.token);
						await startRegistrants(res.data.token);
					}
				}
			});
		}
	});
	const startRegistrants = async (token: string) => {
		await useRegisterCollect.startRegistrants({ phoneToken: token, inviteID: ruleForm.value.inviteID }).then((res: any) => {
			if (res.code === 'SUCCESS') {
				stores.setxtoken(res.data.token);
				stores.setPhone(ruleForm.value.phone);
				localStorage.setItem('xtoken', res.data.token);
				localStorage.setItem('xtokentype', res.data.type);
				localStorage.setItem('xtokensubtype', res.data.subType);
				isGoLogin();
			}
		});
	};
	const isGoLogin = async () => {
		let res = await onSignIn(
			localStorage.getItem('xtoken') as string,
			localStorage.getItem('xtokentype') as string,
			localStorage.getItem('xtokensubtype') as string,
      false,
      true
		);
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
};

</script>

<style scoped lang="scss">
.login-content-form {
	margin-top: 20px;
	@for $i from 1 through 4 {
		.login-animation#{$i} {
			opacity: 0;
			animation-name: error-num;
			animation-duration: 0.5s;
			animation-fill-mode: forwards;
			animation-delay: calc($i/10) + s;
		}
	}
	.login-content-code {
		width: 100%;
		padding: 0;
	}
	.login-content-submit {
		width: 100%;
		letter-spacing: 2px;
		font-weight: 300;
		margin-top: 15px;
	}
	.login-msg {
		color: var(--el-text-color-placeholder);
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
</style>
