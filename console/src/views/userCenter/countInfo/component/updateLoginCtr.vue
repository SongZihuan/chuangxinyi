<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="400px">
		<el-form ref="roleDialogFormRef" :model="ruleForm" size="default" label-width="130px" :rules="rules" @close="closeDialog">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="手机登录权限" prop="name">
						<el-radio-group v-model="ruleForm.allowPhone">
							<el-radio :label="true">允许</el-radio>
							<el-radio :label="false">禁止</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="邮箱登录权限" prop="name">
						<el-radio-group v-model="ruleForm.allowEmail">
							<el-radio :label="true">允许</el-radio>
							<el-radio :label="false">禁止</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="密码登录权限" prop="name">
						<el-radio-group v-model="ruleForm.allowPassword">
							<el-radio :label="true">允许</el-radio>
							<el-radio :label="false">禁止</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="微信登录权限" prop="name">
						<el-radio-group v-model="ruleForm.allowWeChat">
							<el-radio :label="true">允许</el-radio>
							<el-radio :label="false">禁止</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="二次验证登录权限" prop="name">
						<el-radio-group v-model="ruleForm.allowSecondFA">
							<el-radio :label="true">允许</el-radio>
							<el-radio :label="false">禁止</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="onCancel" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(roleDialogFormRef)" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script setup lang="ts" name="updateLoginCtr">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance, ElMessage } from 'element-plus';
import { bannedUpdateTypes } from '/@/views/oauth2/banned/types';
import { userApi } from '/@/api/system/user';
import { checkPhoneRes } from '/@/views/userCenter/countInfo/types';
import { useRegisterApi } from '/@/api/register';

const emit = defineEmits(['refresh']);
// 定义变量内容
const roleDialogFormRef = ref();
let ruleForm = ref<checkPhoneRes>({
	allowPhone: false,
	allowEmail: false,
	allowPassword: false,
	allowWeChat: false,
	allowSecondFA: false,
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '登录控制',
	submitTxt: '修改',
});
const rules = reactive<FormRules>({
	allowPhone: [{ required: true, message: '请选择是否允许手机登录', trigger: 'change' }],
	allowEmail: [{ required: true, message: '请选择是否允许邮箱登录', trigger: 'change' }],
	allowPassword: [{ required: true, message: '请选择是否允许密码登录', trigger: 'change' }],
	allowWeChat: [{ required: true, message: '请选择是否允许微信登录', trigger: 'change' }],
	allowSecondFA: [{ required: true, message: '请选择是否允许二次验证登录', trigger: 'change' }],
});

// 打开弹窗
const openDialog = (row: bannedUpdateTypes, type: string = 'user') => {
	if (type == 'user') {
		dialog.title = '登录控制';
	} else {
		dialog.title = '管理员登录控制';
	}
	dialog.type = type;
	ruleForm.value = row;
	dialog.isShowDialog = true;
};

const closeDialog = () => {
	roleDialogFormRef.value?.resetFields();
	dialog.isShowDialog = false;
};

// 取消
const onCancel = () => {
	closeDialog();
};
// 提交
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			if (dialog.type == 'user') {
				useRegisterApi()
					.updateLoginController(ruleForm.value)
					.then((res: any) => {
						if (res.code === "SUCCESS") {
							ElMessage.success('修改成功');
							closeDialog();
							emit('refresh');
						}
					});
			} else {
				userApi()
					.updateLoginController(ruleForm.value)
					.then((res: any) => {
						if (res.code === "SUCCESS") {
							ElMessage.success('修改成功');
							closeDialog();
							emit('refresh');
						}
					});
			}
		} else {
			return false;
		}
	});
};
// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss"></style>
