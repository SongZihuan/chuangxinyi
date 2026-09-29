<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="800px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="100px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="协议标识" prop="aid">
						<el-input v-model="ruleForm.aid" placeholder="请输入协议标识" clearable :disabled="dialog.type == 'edit'"></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="协议内容" prop="content"> <WangEditor v-model:get-html="ruleForm.content" /></el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { useAgreementApi } from '/@/api/agreement';
import { fromType } from '/@/api/agreement/types';
import WangEditor from '/@/components/Editor/index.vue';
const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<fromType>({
	aid: '',
	content: '',
});
const rules = reactive<FormRules>({
	aid: [{ required: true, message: '请选择协议类型', trigger: 'blur' }],
	content: [{ required: true, message: '请输入协议内容', trigger: 'blur' }],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '协议新增',
	submitTxt: '新增',
});

const openDialog = (type: string, row: any) => {
	if (type === 'add') {
		reset();
		dialog.type = 'add';
		dialog.title = '协议新增';
		dialog.submitTxt = '新增';
	} else {
		dialog.type = 'edit';
		dialog.title = '协议编辑';
		dialog.submitTxt = '编辑';
		ruleForm.value = JSON.parse(JSON.stringify(row));
		getAgreement();
	}
	dialog.isShowDialog = true;
};
const reset = () => {
	ruleForm.value.aid = '';
	ruleForm.value.content = '';
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const getAgreement = () => {
	useAgreementApi()
		.getAgreement({ aid: ruleForm.value.aid })
		.then((res: any) => {
			if (res) {
				ruleForm.value.content = res;
			}
		});
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			if (dialog.type === 'add') {
				useAgreementApi()
					.createAgreement(ruleForm.value)
					.then((res: any) => {
						if (res.code === "SUCCESS") {
							closeDialog();
							ElMessage({
								type: 'success',
								message: '新增协议成功',
							});
							emit('refresh');
						}
					});
			} else {
				useAgreementApi()
					.editAgreement(ruleForm.value)
					.then((res: any) => {
						if (res.code === "SUCCESS") {
							closeDialog();
							ElMessage({
								type: 'success',
								message: '编辑协议成功',
							});
							emit('refresh');
						}
					});
			}
		} else {
			return false;
		}
	});
};
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style lang="scss" scoped></style>
