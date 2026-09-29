<template>
	<div class="system-role-dialog-container">
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px">
			<el-form ref="roleDialogFormRef" :model="ruleForm" size="default" label-width="80px" :rules="rules" @close="closeDialog">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="名称" prop="name">
							<el-input v-model="ruleForm.name" placeholder="请输入名称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="权限标识" prop="sign">
							<el-input v-model="ruleForm.sign" placeholder="请输入权限标识" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="状态" prop="status">
							<el-radio-group v-model="ruleForm.status">
								<el-radio :label="1">启用</el-radio>
								<el-radio :label="2">禁用</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="权限描述" prop="describe">
							<el-input v-model="ruleForm.describe" placeholder="请输入权限描述" clearable type="textarea" :rows="4"></el-input>
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
	</div>
</template>

<script setup lang="ts" name="systemRoleDialog">
import { reactive, ref, nextTick, onMounted } from 'vue';
import { message } from '/@/utils/message';
import { type FormRules, type FormInstance } from 'element-plus';
import { usePermissionApi } from '/@/api/system/websitePermission';
const emit = defineEmits(['refresh']);
const permissionApi = usePermissionApi();
// 定义变量内容
const roleDialogFormRef = ref();
let ruleForm = ref<any>({
	name: '',
	sign: '',
	isAnonymous: '',
	isUser: '',
	status: 1,
	describe: '',
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '',
	submitTxt: '',
});

const rules = reactive<FormRules>({
	name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
	sign: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
	status: [{ required: true, message: '请选择', trigger: 'blur' }],
});

const reset = () => {
	ruleForm.value = {
		name: '', // 权限名称
		sign: '', // 权限标识
		status: 1, // 权限状态
		describe: '', // 权限描述
		sort: '', // 排序
	};
};
// 打开弹窗
const openDialog = (type: string, row: any) => {
	dialog.isShowDialog = true;
	dialog.type = type;

	if (type === 'edit') {
		nextTick(() => {
			ruleForm.value = JSON.parse(JSON.stringify(row));
		});
		dialog.title = '修改权限';
		dialog.submitTxt = '修 改';
	} else {
		reset();
		dialog.title = '新增权限';
		dialog.submitTxt = '新 增';
	}
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
			if (dialog.type === 'add') {
				permissionApi.addPermissionMenu(ruleForm.value).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('新增成功', { type: 'success' });
						closeDialog();
						emit('refresh');
					}
				});
			} else {
				permissionApi.permissionUpdate({ ...ruleForm.value, id: ruleForm.value.id }).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('编辑成功', { type: 'success' });
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
onMounted(() => {});
// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss"></style>
