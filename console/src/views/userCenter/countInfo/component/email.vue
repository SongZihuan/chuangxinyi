<template>
	<!-- 更新邮箱 -->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
		<el-form size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="mobileFormRef">
			<el-col :span="24" class="mb20">
				<el-form-item prop="email">
					<el-input text placeholder="请输入邮箱" v-model="ruleForm.email" clearable autocomplete="off">
						<template #prefix>
							<el-icon class="el-input__icon"><ele-Message /></el-icon>
						</template>
					</el-input>
				</el-form-item>
			</el-col>
			<el-form-item>
				<el-col :span="24">
					<slider-silence @siderEmit="siderEmit" ref="sliderRef" />
				</el-col>
			</el-form-item>
			<el-form-item prop="code" v-if="isSliderCheck">
				<el-col :span="15">
					<el-input text maxlength="6" placeholder="请输入验证码" v-model="ruleForm.code" clearable autocomplete="off">
						<template #prefix>
							<el-icon class="el-input__icon">
								<ele-Position />
							</el-icon>
						</template>
					</el-input>
				</el-col>
				<el-col :span="1"></el-col>
				<el-col :span="8">
					<el-button v-waves class="login-content-code" @click="getCode" :disabled="codeState.isSend">
						{{ codeState.codeName }}
					</el-button>
				</el-col>
			</el-form-item>
			<el-form-item>
				<el-button round type="primary" v-waves class="login-content-submit" @click="onSubmit(mobileFormRef)">
					<span>绑定</span>
				</el-button>
			</el-form-item>
		</el-form>
	</el-dialog>
</template>
<script setup lang="ts" name="userEmail">
import { reactive, ref } from 'vue';
import sliderSilence from '/@/components/Slider/index.vue';
import { verifyEmail } from '/@/utils/toolsValidate';
import { type FormRules, type FormInstance } from 'element-plus';
import type { codeTypes } from '/@/api/register/types';
import type { sliderHeadersTypes, checkEmailTypes } from '/@/api/base/types';
import { useBaseApi } from '/@/api/base/index';
import { ElMessage } from 'element-plus';
import { useUserApi } from '/@/api/user/user';
import { useRegisterApi } from '/@/api/register/index';

const emit = defineEmits(['refresh']);

interface Props {
	emailForm: checkEmailTypes;
}

const props = withDefaults(defineProps<Props>(), {});
const useUserApiCollect = useUserApi();
const useBaseApiCollect = useBaseApi();
const useRegisterCollect = useRegisterApi();
// 定义变量内容
const ruleForm = ref<checkEmailTypes>(props.emailForm);
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	email: [{ trigger: 'blur', validator: verifyEmail }],
	code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
});
const sliderRef = ref();
//阿里云滑块返回的参数
let sliderHeaders = ref<sliderHeadersTypes>({
	sig: '',
	sessionId: '',
	token: '',
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '修改绑定的邮箱',
	submitTxt: '确认',
});
const openDialog = () => {
	dialog.isShowDialog = true;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
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
	await useBaseApiCollect
		.sendEmailCode({ email: ruleForm.value.email }, sliderHeaders.value)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				isGetCode.value = true;
				ElMessage.success('发送邮箱验证码成功');
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
						codeState.totalTime = 10;
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
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!isSliderCheck.value) {
		ElMessage.warning('请滑动滑块移动最右侧!');
		return;
	}

	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			await useRegisterCollect.checkEmailCode(ruleForm.value, sliderHeaders.value).then((res: any) => {
				isSliderCheck.value = false;
				sliderRef.value.resetSider();
				if (res.code === "SUCCESS") {
					useUserApiCollect.updateEmail({ emailToken: res.data.token, isDelete: false }).then((res: any) => {
						if (res.code === "SUCCESS") {
							ElMessage.success('修改邮箱绑定成功');
							emit('refresh');
							closeDialog();
						}
					});
				}
			});
		}
	});
};
defineExpose({
	openDialog,
});
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
</style>
