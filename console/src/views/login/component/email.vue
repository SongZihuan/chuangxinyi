<template>
	<el-form size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="formRef">
		<el-form-item class="login-animation1 mb20" prop="email">
			<el-input text placeholder="请输入邮箱" v-model="ruleForm.email" clearable autocomplete="off">
				<template #prefix>
					<el-icon class="el-input__icon"><ele-Message /></el-icon>
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
			<el-button round type="primary" v-waves class="login-content-submit" @click="onLogin(formRef)">
				<span>登录</span>
			</el-button>
		</el-form-item>
	</el-form>
</template>
<script setup lang="ts" name="registerMobile">
import { reactive, ref, onMounted, onBeforeUnmount } from 'vue';
import sliderSilence from '/@/components/Slider/index.vue';
import { verifyEmail } from '/@/utils/toolsValidate';
import { type FormRules, type FormInstance } from 'element-plus';
import type { codeTypes } from '/@/api/register/types';
import type { sliderHeadersTypes, checkEmailTypes } from '/@/api/base/types';
import { useBaseApi } from '/@/api/base/index';
import { useLoginApi } from '/@/api/login/index';
import { ElMessage } from 'element-plus';
import { useRouter, useRoute } from 'vue-router';
import { useUserInfo } from '/@/stores/userInfo';
const router = useRouter();
const route = useRoute();
const useBaseApiCollect = useBaseApi();
const useLoginrCollect = useLoginApi();
const emit = defineEmits(['signBack', 'totpBack']);
interface Props {
	isAgree: boolean; // 回显图片地址
}
const props = withDefaults(defineProps<Props>(), {
	isAgree: false,
});
const ruleForm = ref<checkEmailTypes>({
	email: '',
	code: '',
	type: 'UserToken',
});
const formRef = ref();
const rules = reactive<FormRules>({
	email: [{ trigger: 'blur', validator: verifyEmail }],
	code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
});
const stores = useUserInfo();
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
const isGetCode = ref<boolean>(false);
const codeState = reactive({
	isSend: false,
	codeName: '获取验证码',
	totalTime: 60, //一般是60
	timer: undefined, //定时器
}) as codeTypes;
const getCode = async () => {
	await useBaseApiCollect
		.sendEmailCode({ email: ruleForm.value.email }, sliderHeaders.value)
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
};
const goPassword = () => {
	router.push({ path: '/forgetPassword', query: route.query});
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
