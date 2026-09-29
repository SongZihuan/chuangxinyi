<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="100px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="协作人" prop="uncleID">
						<el-input v-model="ruleForm.uncleID" placeholder="请输入协作人" style="width: 100%"></el-input>
					</el-form-item>
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
import { addUncleTypes } from '/@/views/accountManagement/uncleUser/types';
import { useUncleUserApi } from '/@/api/uncleUser';

const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<addUncleTypes>({
	uncleID: '',
});
const rules = reactive<FormRules>({
	uncleID: [{ required: true, message: '请输入协作人', trigger: 'blur' }],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '协作人新增',
	submitTxt: '新增',
});

const openDialog = () => {
	dialog.isShowDialog = true;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useUncleUserApi()
				.addUncleUser(ruleForm.value)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '新增协作人成功',
						});
						emit('refresh');
					}
				});
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

<style lang="scss" scoped>
.ifr {
	width: 100%;
	height: 560px;
}

.tip {
	display: flex;
	flex-direction: row;
	justify-content: center;
}
</style>
