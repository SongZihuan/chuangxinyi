<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="375px">
		<el-form size="large" class="login-content-form" ref="formRef" :rules="rules" :model="data.ruleForm">
			<el-form-item class="login-animation1" prop="code">
				<el-input text placeholder="2FA验证码" v-model="data.ruleForm.code" clearable autocomplete="off">
					<template #prefix>
						<el-icon class="el-input__icon"><ele-Position /></el-icon>
					</template>
				</el-input>
			</el-form-item>
			<el-form-item class="login-animation1 mb20" prop="isRember">
				<el-checkbox v-model="data.ruleForm.isRember">信任该机器</el-checkbox>
			</el-form-item>
			<el-form-item class="login-animation3">
				<el-col :span="24">
					<silenceSlider ref="sliderRef" @nvcValEmit="nvcValEmit" @siderEmit="siderEmit" />
				</el-col>
			</el-form-item>
			<el-form-item class="login-animation4">
				<el-button type="primary" class="login-content-submit" round v-waves @click="onLogin(formRef)" block :loading="data.loading.signIn">
					<span>确认</span>
				</el-button>
			</el-form-item>
		</el-form>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, toRefs, onMounted, ref } from 'vue';
import { type FormRules, type FormInstance, ElMessage } from 'element-plus';
import silenceSlider from '/@/components/Slider/silence.vue';
import { useGoogleAuthApi } from '/@/api/Totp/index';
import { useoauth2Api } from '/@/api/oauth2/index';
import { useUserInfo } from '/@/stores/userInfo';
import { Local, Session } from '/@/utils/storage';
import { useLoginSignIn } from '/@/hooks/useLoginSignIn';
import { useRoute, useRouter } from 'vue-router';
import Cookies from 'js-cookie';
import { NextLoading } from '/@/utils/loading';
import { encodeSearchParams } from '/@/utils/query';
const { onSignIn } = useLoginSignIn();
const useGoogleAuthApiCollect = useGoogleAuthApi();
const stores = useUserInfo();
const router = useRouter();
const route = useRoute();
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '2FA验证',
	},
});
const formRef = ref();
const sliderRef = ref();
const data = reactive({
	isShowPassword: false,
	ruleForm: {
		code: '',
		token: '',
		passToken: '',
		isRember: false,
	},
	loading: {
		signIn: false,
	},
});
const rules = reactive<FormRules>({
	code: [{ required: true, message: '请输入2FA验证码', trigger: 'blur' }],
});
const { dialog } = toRefs(state);
const isSliderCheck = ref(false);
const isFuwuhao = ref(false)


const openDialog = (token: string, isFWH: boolean = false) => {
	data.ruleForm.code = '';
	dialog.value.isShowDialog = true;
	data.ruleForm.token = token;
  isFuwuhao.value = isFWH
};
const directlogin = (token: string, isFWH: boolean = false) => {
	data.ruleForm.code = '';
	data.ruleForm.token = token;
  isFuwuhao.value = isFWH
	checkTotp(data.ruleForm.code, data.ruleForm.token, currentVal.value);
};

const nvcValEmit = (val: string) => {
	checkTotp(data.ruleForm.code, data.ruleForm.token, val);
};
const currentVal = ref();
const isSecond = ref(false);
const siderEmit = (val: string) => {
	currentVal.value = val;
	isSecond.value = true;
	isSliderCheck.value = true;
};
const onLogin = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	if (!isSliderCheck.value && isSecond.value) {
		ElMessage.warning('请滑动滑块移动最右侧!');
		return;
	}
	formEl.validate(async (valid) => {
		if (valid) {
			if (isSecond.value) {
				checkTotp(data.ruleForm.code, data.ruleForm.token, currentVal.value);
			} else {
				sliderRef.value.registerClick();
			}
		}
	});
};

const checkTotp = (code: string, token: string, nvc: string) => {
	useGoogleAuthApiCollect
		.checkTotp({
			code: code,
			token: token,
			nvc: nvc || '',
			rememberHour: data.ruleForm.isRember ? 7 * 24 : 0,
			passToken: code ? '' : (Cookies.get('passToken') as string), // 没有code的时候才使用passToken
		})
		.then(async (res: any) => {
			sliderRef.value?.resetSider();
			await stores.setUserType({ type: res.data.type, subType: res.data.subType });
			if (res.code === 'LOGIC_DENY' && res.subCode === 'CAPTCHA_SECOND_CHECK') {
				sliderRef.value.siderCheck();
				data.ruleForm.code = '';
				return;
			} else if (res.code === 'LOGIC_DENY' && res.subCode === 'BAD_PASS_TOKEN') {
				Cookies.set('passToken', '', { sameSite: 'Strict' });
				data.ruleForm.code = '';
				dialog.value.isShowDialog = true;
			} else if (res.code === 'SUCCESS') {
				let domain_uid = route.query?.domain;
				let redirect_uri = route.query?.redirect_uri;
				let next_redirect_uri = route.query?.next_redirect_uri;
				let params = route.query?.params;
				let isLoginToken = Number(route.query?.isLoginToken) == 1;
				dialog.value.isShowDialog = false;
				//如果勾选了信任该机器 下次2fa登录使用passToken
				if (data.ruleForm.isRember && res.data.subToken) {
					//七天有效期
					Cookies.set('passToken', res.data.subToken, { expires: 7, sameSite: 'Strict' });
				}
				let signres = await onSignIn(res.data.token, res.data.type, res.data.subType);
				if (signres) {
          if (isFuwuhao.value) {
            sessionStorage.setItem("fuwuhao-login", "true")
          }
					if (route.query?.type === 'fromOauth2' || route.query?.type === 'fromApplication') {
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
							.then(async (res: any) => {
								if (res.code === 'SUCCESS') {
									let queryString = encodeSearchParams({
										token: res.data.token,
										subToken: res.data.subToken,
										params: params,
										redirect: next_redirect_uri,
									});

									if (route.query?.type === 'fromApplication') {
										window.open(`${redirect_uri}?${queryString}`);
										await router.replace('/');
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
						await router.replace({
							path: <string>route.query?.redirect,
							query: Object.keys(<string>route.query?.params).length > 0 ? JSON.parse(<string>route.query?.params) : '',
						});
					} else {
						await router.replace('/');
					}
				} else {
					ElMessage.error('2FA错误');
				}
				NextLoading.done();
				state.dialog.isShowDialog = false;
			}
		});
};

onMounted(() => {});
defineExpose({
	openDialog,
	directlogin,
});
</script>

<style lang="scss" scoped>
.forgetPassword {
	margin-top: 10px;
	text-align: right;
	color: var(--el-color-primary);
}
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

	.login-content-password {
		display: inline-block;
		width: 20px;
		cursor: pointer;
		&:hover {
			color: #909399;
		}
	}
	.login-content-code {
		width: 100%;
		padding: 0;
		font-weight: bold;
		letter-spacing: 5px;
	}
	.login-content-submit {
		width: 100%;
		letter-spacing: 2px;
		font-weight: 300;
		margin-top: 15px;
	}
}
</style>
