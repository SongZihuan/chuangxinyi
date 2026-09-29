<template>
	<!-- 手机验证码弹窗 -->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
		<el-form size="large" :rules="rules" :model="ruleForm" ref="mobileFormRef">
			<!--      选择手机号、邮箱、2fa验证-->
			<el-form-item class="login-animation1 mb20" prop="type">
				<el-radio-group v-model="ruleForm.type">
					<el-radio label="PhoneCheck">手机验证</el-radio>
					<el-radio label="EmailCheck">邮箱验证</el-radio>
					<el-radio label="2faCheck">2fa验证</el-radio>
				</el-radio-group>
			</el-form-item>
			<el-form-item prop="phone" class="mb20" v-if="ruleForm.type === 'PhoneCheck'">
				<el-input text placeholder="请输入手机号" v-model="ruleForm.phone" clearable autocomplete="off">
					<template #prefix>
						<el-icon><ele-Phone /></el-icon>
					</template>
				</el-input>
			</el-form-item>
			<el-form-item class="login-animation1 mb20" prop="email" v-if="ruleForm.type === 'EmailCheck'">
				<el-input text placeholder="请输入邮箱" v-model="ruleForm.email" clearable autocomplete="off">
					<template #prefix>
						<el-icon class="el-input__icon"><ele-Message /></el-icon>
					</template>
				</el-input>
			</el-form-item>
			<el-form-item class="login-animation1 mb20" prop="code" v-if="ruleForm.type === '2faCheck'">
				<el-input text placeholder="2FA验证码" v-model="ruleForm.code2FA" clearable autocomplete="off">
					<template #prefix>
						<el-icon class="el-input__icon"><ele-Position /></el-icon>
					</template>
				</el-input>
			</el-form-item>
			<el-form-item>
				<el-col :span="24">
					<slider-silence :is-center="true" @siderEmit="siderEmit" ref="sliderRef" />
				</el-col>
			</el-form-item>
			<el-row v-if="ruleForm.type !== '2faCheck'">
				<el-form-item prop="code" v-if="isSliderCheck" class="code-box">
					<el-col :span="15">
						<el-input text maxlength="6" placeholder="请输入验证码" v-model="ruleForm.code" clearable autocomplete="off">
							<template #prefix>
								<el-icon class="el-input__icon"><ele-Position /></el-icon>
							</template>
						</el-input>
					</el-col>
					<el-col :span="1"></el-col>
					<el-col :span="8">
						<el-button v-waves @click="getCode" :disabled="codeState.isSend" class="code"> {{ codeState.codeName }}</el-button>
					</el-col>
				</el-form-item>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="dialog.isShowDialog = false">取消</el-button>
				<el-button type="danger" @click="confirm(mobileFormRef)"> {{ dialog.submitTxt }} </el-button>
			</span>
		</template>
	</el-dialog>
</template>
<script setup lang="ts" name="doubleCheck">
import { reactive, ref } from 'vue';
import sliderSilence from '/@/components/Slider/index.vue';
import { useBaseApi } from '/@/api/base';
import { verifyEmail, verifyPhone } from '/@/utils/toolsValidate';
import { type FormRules, type FormInstance } from 'element-plus';
import type { sliderHeadersTypes, codeTypes } from '/@/api/base/types';
import { Session, Local } from '/@/utils/storage';
import { ElMessage } from 'element-plus';

const useBaseApiCollect = useBaseApi();

const emit = defineEmits(['reset', 'checkPhoneSuccess']);
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '验证',
	submitTxt: '确认',
});
const openDialog = (row: any, title: string) => {
	dialog.title = title;
	dialog.isShowDialog = true;
	reset();
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
// 定义变量内容
const ruleForm = ref({
	phone: '',
	email: '',
	code: '',
	token: '',
	code2FA: '',
	type: 'PhoneCheck',
});
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	phone: [{ trigger: 'blur', validator: verifyPhone }],
	email: [{ trigger: 'blur', validator: verifyEmail }],
	code2FA: [{ required: true, message: '请输入2FA验证码', trigger: 'blur' }],
	code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
});
const userInfo = Session.get('userInfo') || Local.get('userInfo');
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
const reset = () => {
	ruleForm.value.phone = userInfo.user.phone;
	ruleForm.value.email = userInfo.user.email;
	ruleForm.value.code = '';
	isSliderCheck.value = false;
	isGetCode.value = false;
};
const isGetCode = ref<boolean>(false);
const codeState = reactive({
	isSend: false,
	codeName: '获取验证码',
	totalTime: 60, //一般是60
	timer: undefined, //定时器
}) as codeTypes;
const getCode = async () => {
	if (ruleForm.value.type === 'PhoneCheck') {
		await useBaseApiCollect
			.sendPhoneCode({ phone: ruleForm.value.phone }, sliderHeaders.value)
			.then((res: any) => {
				if (res.code === 'SUCCESS') {
					sendCode();
				} else {
					sliderRef.value.resetSider();
				}
			})
			.catch(() => {
				sliderRef.value.resetSider();
			});
	} else if (ruleForm.value.type === 'EmailCheck') {
		await useBaseApiCollect
			.sendEmailCode({ email: ruleForm.value.email }, sliderHeaders.value)
			.then((res: any) => {
				if (res.code === 'SUCCESS') {
					sendCode();
				} else {
					sliderRef.value.resetSider();
				}
			})
			.catch(() => {
				sliderRef.value.resetSider();
			});
	}
};
const sendCode = () => {
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
		}
	}, 1000);
	ElMessage.success(`发送${ruleForm.value.type === 'PhoneCheck' ? '手机' : '邮箱'}验证码成功`);
};
const confirm = (formEl: FormInstance | undefined) => {
	if (!isSliderCheck.value) {
		ElMessage.warning('请滑动滑块移动最右侧!');
		return;
	}
	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			// 清除验证码倒计时
			if (codeState.timer) {
				clearInterval(codeState.timer);
				codeState.codeName = '重新发送';
				codeState.totalTime = 60;
				codeState.isSend = false;
			}
			if (ruleForm.value.type === 'PhoneCheck') {
				await useBaseApiCollect.checkPhoneCode({ ...ruleForm.value, type: 'PhoneCheck' }).then((res: any) => {
					isSliderCheck.value = false;
					sliderRef.value.resetSider();
					if (res.code === 'SUCCESS') {
						const data = {
							'X-Phone-Token': res.data.token,
						};
						emit('checkPhoneSuccess', data);
						closeDialog();
					}
				});
			} else if (ruleForm.value.type === 'EmailCheck') {
				await useBaseApiCollect.checkEmailCode({ ...ruleForm.value, type: 'EmailCheck' }).then((res: any) => {
					isSliderCheck.value = false;
					sliderRef.value.resetSider();
					if (res.code === 'SUCCESS') {
						const data = {
							'X-Email-Token': res.data.token,
						};
						emit('checkPhoneSuccess', data);
						closeDialog();
					}
				});
			} else if (ruleForm.value.type === '2faCheck') {
				ruleForm.value.token = Session.get('center-token');
				if (!ruleForm.value.token) {
					ElMessage.warning('2fa验证失败，请重新登录');
					return;
				}
				await useBaseApiCollect.check2faCode({ ...ruleForm.value }).then((res: any) => {
					isSliderCheck.value = false;
					sliderRef.value.resetSider();
					if (res.code === 'SUCCESS') {
						const data = {
							'X-2FA-Token': res.data.token,
						};
						emit('checkPhoneSuccess', data);
						closeDialog();
					}
				});
			}
		}
	});
};
defineExpose({
	openDialog,
});
</script>

<style lang="scss" scoped>
.code-box {
	width: 500px;
	.code {
		width: 100%;
	}
}
</style>
