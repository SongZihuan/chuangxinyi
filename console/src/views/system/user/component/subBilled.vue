<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="110px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="用户ID" prop="id">
						<el-input v-model="ruleForm.id" placeholder="请输入用户ID" clearable disabled></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="开票金额" prop="amount">
						<template #label>开票金额( <el-text class="mx-1" type="danger" tag="b">分</el-text>)</template>
						<el-input-number
							v-model="ruleForm.amount"
							placeholder="请输入开票金额"
							clearable
							:min="0"
							style="width: 100%"
							:step="1"
							step-strictly
							controls-position="right"
						></el-input-number>
					</el-form-item>
				</el-col>
			</el-row> </el-form
		><template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">确认</el-button>
			</span>
		</template>
		<checkPhone ref="checkPhoneRef" @checkPhoneSuccess="checkPhoneSuccess" />
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance, ElMessage } from 'element-plus';
import checkPhone from '/@/components/checkPhone/index.vue';
import { useDefraytApi } from '/@/api/defray';
const formRef = ref();
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '减少开票额度',
	submitTxt: '',
});
const isShowDialog = ref(false);
const checkPhoneRef = ref();
const rules = reactive<FormRules>({
	amount: [
		{
			required: true,
			message: '请输入增加金额',
			trigger: 'blur',
		},
	],
});
const ruleForm = ref({
	amount: 1,
	id: null,
});

const closeDialog = () => {
	dialog.isShowDialog = false;
	isShowDialog.value = false;
};
const reset = () => {
	ruleForm.value = {
		amount: 1,
		id: null,
	};
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			checkPhoneRef.value.openDialog({}, '减少开票额度确认');
		} else {
			return false;
		}
	});
};
const checkPhoneSuccess = (phoneToken: string) => {
	useDefraytApi()
		.subBilled({ ...ruleForm.value, phoneToken: phoneToken })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				ElMessage({
					type: 'success',
					message: '减少开票额度成功',
				});
				closeDialog();
			}
		});
};

const openDialog = async (row?: any) => {
	reset();
	ruleForm.value.id = JSON.parse(JSON.stringify(row)).numberID;
	dialog.isShowDialog = true;
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>
