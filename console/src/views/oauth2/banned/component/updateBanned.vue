<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="569px">
		<el-form ref="roleDialogFormRef" :model="ruleForm" size="default" label-width="120px" :rules="rules" @close="closeDialog">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="网址登录权限" prop="name">
						<el-radio-group v-model="ruleForm.allowLogin">
							<el-radio :label="true">允许</el-radio>
							<el-radio :label="false">禁止</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="网址支付权限" prop="name">
						<el-radio-group v-model="ruleForm.allowDefray">
							<el-radio :label="true">允许</el-radio>
							<el-radio :label="false">禁止</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="网址消息权限" prop="name">
						<el-radio-group v-model="ruleForm.allowMsg">
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

<script setup lang="ts" name="updateBanned">
import { reactive, ref } from 'vue';
import { message } from '/@/utils/message';
import { type FormRules, type FormInstance } from 'element-plus';
import { bannedUpdateTypes } from '/@/views/oauth2/banned/types';
import { useBannedApi } from '/@/api/banned';
import { userApi } from '/@/api/system/user';

const emit = defineEmits(['refresh']);
// 定义变量内容
const roleDialogFormRef = ref();
let ruleForm = ref<bannedUpdateTypes>({
	webID: 0,
	allowLogin: false,
	allowDefray: false,
	allowMsg: false,
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '网站更新',
	submitTxt: '修改',
});
const rules = reactive<FormRules>({
	webID: [{ required: true, message: '请选择网站', trigger: 'change' }],
	allowLogin: [{ required: true, message: '请选择是否允许登录', trigger: 'change' }],
	allowDefray: [{ required: true, message: '请选择是否允许支付', trigger: 'change' }],
	allowMsg: [{ required: true, message: '请选择是否通信', trigger: 'change' }],
});

// 打开弹窗
const openDialog = (row: bannedUpdateTypes, type: string = 'user') => {
	if (type == 'user') {
		dialog.title = '用户更新';
	} else {
		dialog.title = '管理员更新';
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
				useBannedApi()
					.updateBanned(ruleForm.value)
					.then((res: any) => {
						if (res.code === "SUCCESS") {
							message('编辑成功', { type: 'success' });
							emit('refresh');
							closeDialog();
						}
					});
			} else {
				userApi()
					.updateUserBanStatus(ruleForm.value)
					.then((res: any) => {
						if (res.code === "SUCCESS") {
							message('编辑成功', { type: 'success' });
							emit('refresh');
							closeDialog();
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
