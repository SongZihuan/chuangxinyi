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
import { verifyPhone } from '/@/utils/toolsValidate';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { userApi } from '/@/api/system/user';
import { userPhoneTypes } from '/@/views/system/user/types';

const emit = defineEmits(['refresh']);
interface Props {
	phoneForm: userPhoneTypes;
}

const props = withDefaults(defineProps<Props>(), {});
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
const ruleForm = ref<userPhoneTypes>(props.phoneForm);

const mobileFormRef = ref();
const rules = reactive<FormRules>({
	phone: [{ trigger: 'blur', validator: verifyPhone }],
});
const confirmLogOff = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	if (!ruleForm.value.uid) {
		ElMessage.error('用户获取失败');
		return;
	}
	formEl.validate(async (valid) => {
		if (valid) {
			await userApi()
				.updateUserPhone(ruleForm.value)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('修改绑定手机号成功');
						emit('refresh');
						closeDialog();
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
