<template>
	<!-- 手机验证码弹窗 -->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
		<el-form size="large" :rules="rules" :model="ruleForm" ref="mobileFormRef">
			<el-form-item prop="phone">
				<el-col :span="24" class="mb20">
					<el-input text placeholder="请输入手机号" v-model="ruleForm.phone" clearable autocomplete="off">
						<template #prefix>
							<el-icon><ele-Phone /></el-icon>
						</template>
					</el-input>
				</el-col>
			</el-form-item>
			<el-form-item>
				<el-col :span="24">
					<slider-silence @siderEmit="siderEmit" ref="sliderRef" />
				</el-col>
			</el-form-item>
			<el-row>
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
				<el-button type="danger" @click="confirmLogOff(mobileFormRef)"> {{ dialog.submitTxt }} </el-button>
			</span>
		</template>
	</el-dialog>
</template>
<script setup lang="ts" name="userMobile">
import { reactive, ref } from 'vue';
import sliderSilence from '/@/components/Slider/index.vue';
import { useBaseApi } from '/@/api/base/index';
import { verifyPhone } from '/@/utils/toolsValidate';
import { type FormRules, type FormInstance } from 'element-plus';
import type { checkPhoneCodeTypes, codeTypes } from '/@/api/register/types';
import type { sliderHeadersTypes } from '/@/api/base/types';
import { useRegisterApi } from '/@/api/register/index';
import { ElMessage } from 'element-plus';
import { useUserApi } from '/@/api/user/user';
import { checkPhoneTypes } from '/@/api/base/types';

const emit = defineEmits(['refresh']);
interface Props {
	phoneForm: checkPhoneTypes;
}

const props = withDefaults(defineProps<Props>(), {});
const useUserApiCollect = useUserApi();
const useBaseApiCollect = useBaseApi();
const useRegisterCollect = useRegisterApi();
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '修改绑定的手机号',
	submitTxt: '确认',
});
const openDialog = () => {
	dialog.isShowDialog = true;
};

const closeDialog = () => {
	dialog.isShowDialog = false;
};
// 定义变量内容
const ruleForm = ref<checkPhoneCodeTypes>(props.phoneForm);

const mobileFormRef = ref();
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
const isGetCode = ref<boolean>(false);
const codeState = reactive({
	isSend: false,
	codeName: '获取验证码',
	totalTime: 60, //一般是60
	timer: undefined, //定时器
}) as codeTypes;
const getCode = async () => {
	await useBaseApiCollect.sendPhoneCode({ phone: ruleForm.value.phone }, sliderHeaders.value).then((res: any) => {
		if (res.code === "SUCCESS") {
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
		}else{
      sliderRef.value.resetSider();
    }
	}).catch(()=>{
    sliderRef.value.resetSider();
  });
};
const confirmLogOff = (formEl: FormInstance | undefined) => {
	if (!isSliderCheck.value) {
		ElMessage.warning('请滑动滑块移动最右侧!');
		return;
	}

	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			await useRegisterCollect.checkPhoneCode(ruleForm.value, sliderHeaders.value).then((res: any) => {
				isSliderCheck.value = false;
				sliderRef.value.resetSider();
				if (res.code === "SUCCESS") {
					useUserApiCollect.updatePhone({ phoneToken: res.data.token }).then((res: any) => {
						if (res.code === "SUCCESS") {
							ElMessage.success('修改绑定手机号成功');
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

<style lang="scss" scoped>
.code-box {
	width: 500px;
	.code {
		width: 100%;
	}
}
</style>
