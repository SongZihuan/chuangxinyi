<template>
	<el-form size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="mobileFormRef">
		<el-form-item class="login-animation1 mb20" prop="phone">
			<el-input text placeholder="请输入手机号或邮箱" v-model="ruleForm.phone" clearable autocomplete="off">
				<template #prefix>
					<el-icon><ele-Phone /></el-icon>
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
				<el-input text maxlength="6" placeholder="请输入验证码" v-model="ruleForm.code" clearable autocomplete="off">
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
		<div class="forgetPassword login-animation2" @click="goPassword">
			<div>忘记密码?</div>
		</div>
		<el-form-item class="login-animation3">
			<el-button round type="primary" v-waves class="login-content-submit" @click.prevent="onLogin(mobileFormRef)">
				<span>登录</span>
			</el-button>
		</el-form-item>
	</el-form>
</template>
<script setup lang="ts">
import { reactive, ref, computed} from 'vue';
import sliderSilence from '/@/components/Slider/index.vue';
import { useLoginApi } from '/@/api/login/index';
import { verifyPhoneOrEmail } from '/@/utils/toolsValidate';
import { useBaseApi } from '/@/api/base/index';
import { type FormRules, type FormInstance } from 'element-plus';
import type { checkPhoneCodeTypes, codeTypes } from '/@/api/register/types';
import type { sliderHeadersTypes } from '/@/api/base/types';
import { useUserInfo } from '/@/stores/userInfo';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRegisterApi } from '/@/api/register/index';

const useRegisterCollect = useRegisterApi();
const router = useRouter();
const route = useRoute();
const emit = defineEmits(['signBack', 'totpBack', 'openForgetPassword']);
const stores = useUserInfo();
interface Props {
	isAgree: boolean;
	isBindWeixin: boolean;
}
const props = withDefaults(defineProps<Props>(), {
	isAgree: false,
	isBindWeixin: false,
});
const useBaseApiCollect = useBaseApi();
const useLoginrCollect = useLoginApi();

const isPhone = computed(() => {
	if (ruleForm.value.phone) {
		return ruleForm.value.phone.indexOf('@') === -1;
	}
	return true;
});

// 定义变量内容
const ruleForm = ref<checkPhoneCodeTypes>({
	phone: '',
	code: '',
	type: 'Auto',
});
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	phone: [{ trigger: 'blur', validator: verifyPhoneOrEmail }],
	code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
});
const sliderRef = ref();
//阿里云滑块返回的参数
let sliderHeaders = ref<sliderHeadersTypes>({
	sig: '',
	sessionId: '',
	token: '',
});
const goPassword = () => {
	router.push({ path: '/forgetPassword', query: route.query });
};
const isSliderCheck = ref<boolean>(false);
const siderEmit = (data: any) => {
	sliderHeaders.value = data;
	isSliderCheck.value = true;
};
const isGetCode = ref<boolean>(false);
const codeState = reactive({
	isSend: false,
	codeName: '获取验证码',
	totalTime: 60, //一般是60
	timer: undefined, //定时器
}) as codeTypes;
const getCode = async () => {
	if (!ruleForm.value.phone) {
		ElMessage.warning('请输入手机号或邮箱');
	}

	if (isPhone.value) {
		await useBaseApiCollect
			.sendPhoneCode({ phone: ruleForm.value.phone }, sliderHeaders.value)
			.then((res: any) => {
				if (res.code === 'SUCCESS') {
					isGetCode.value = true;
					sliderRef.value.resetSider();
					ElMessage.success('发送手机验证码成功');
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
						}
					}, 1000);
				} else {
					sliderRef.value.resetSider();
				}
			})
			.catch(() => {
				sliderRef.value.resetSider();
			});
	} else {
		await useBaseApiCollect
			.sendEmailCode({ email: ruleForm.value.phone }, sliderHeaders.value)
			.then((res: any) => {
				if (res.code === 'SUCCESS') {
					isGetCode.value = true;
					if (codeState.isSend) return;
					// getCode() // 获取验证码的接口
					codeState.isSend = true;
					sliderRef.value.resetSider();
					codeState.codeName = codeState.totalTime + 's后重新发送';
					codeState.timer = setInterval(() => {
						codeState.totalTime--;
						codeState.codeName = codeState.totalTime + 's后重新发送';
						if (codeState.totalTime < 0) {
							clearInterval(codeState.timer);
							codeState.codeName = '重新发送';
							codeState.totalTime = 60;
							codeState.isSend = false;
						}
					}, 1000);
					ElMessage.success('发送邮箱验证码成功');
				} else {
					sliderRef.value.resetSider();
				}
			})
			.catch(() => {
				sliderRef.value.resetSider();
			});
	}
};
const onLogin = (formEl: FormInstance | undefined) => {
	if (!isSliderCheck.value) {
		ElMessage.warning('请滑动滑块移动最右侧!');
		return;
	}
	if (!props.isAgree) {
		ElMessage.warning('请勾选用户协议');
		return;
	}
	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			if (isPhone.value) {
				await useLoginrCollect.phoneLogin(ruleForm.value, sliderHeaders.value).then(async (res: any) => {
					isSliderCheck.value = false;
					sliderRef.value.resetSider();
					if (res.code === 'SUCCESS') {
						await stores.setUserType({ type: res.data.type, subType: res.data.subType });
						if (res.data.type == 'UserToken') {
							emit('signBack', res.data.token, res.data.type, res.data.subType);
						} else if (res.data.type === 'Login2FA') {
							emit('totpBack', res.data.token);
						} else {
							ElMessageBox.confirm('您的手机号未注册，请问是否直接注册 ?', '提示', {
								distinguishCancelAndClose: true,
								confirmButtonText: '直接注册',
								cancelButtonText: '取消',
							})
								.then(() => {
									useRegisterCollect.startRegistrants({ phoneToken: res.data.token }).then(async (res: any) => {
										if (res.code === 'SUCCESS') {
											await stores.setxtoken(res.data.token);
											await stores.setPhone(ruleForm.value.phone);
											localStorage.setItem('xtoken', res.data.token);
											await stores.setUserType({ type: res.data.type, subType: res.data.subType });

											emit('signBack', res.data.token, res.data.type, res.data.subType);
										}
									});
								})
								.catch(() => {});
						}
					}
			
				});
			} else {
				ruleForm.value.email = ruleForm.value.phone;
				ruleForm.value.type = 'UserToken';
				await useLoginrCollect.emailLogin(ruleForm.value, sliderHeaders.value).then(async (res: any) => {
					isSliderCheck.value = false;
					sliderRef.value.resetSider();
					if (res.code === 'SUCCESS') {
						await stores.setUserType({ type: res.data.type, subType: res.data.subType });
						if (res.data.type === 'Login2FA') {
							emit('totpBack', res.data.token);
						} else {
							emit('signBack', res.data.token, res.data.type, res.data.subType);
						}
					} else {
						sliderRef.value.resetSider();
					}
				});
			}
		}
	});
};

</script>

<style scoped lang="scss">
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
</style>
