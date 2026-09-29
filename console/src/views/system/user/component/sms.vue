<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" height="500px" center>
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="模版ID" prop="templateID">
						<el-input v-model="ruleForm.templateID" placeholder="请输入模版ID" clearable></el-input>
					</el-form-item>
				</el-col>

				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="短信签名" prop="sig">
						<el-input v-model="ruleForm.sig" placeholder="请输入短信签名" clearable ></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="参数列表" prop="data" v-for="(item, index) in ruleForm.data" :key="index">
						<el-input v-model="item.label" placeholder="参数" style="width: 35%"></el-input>
						<el-input v-model="item.value" placeholder="参数值" style="width: 35%" class="ml10"></el-input>
						<el-button type="primary" @click="handleAdd" class="ml5" text v-if="index === 0">新增</el-button>
						<el-button type="danger" @click="handleDel(index)" class="ml5" text v-else>删除</el-button>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">发送</el-button>
			</span>
		</template></el-dialog
	>
</template>

<script setup lang="ts">
import { reactive, toRefs, ref } from 'vue';
import { type FormRules, type FormInstance, ElMessage } from 'element-plus';
import { useMessageAdminApi } from '/@/api/message/admin';

const state = reactive({
	dialog: {
		isShowDialog: false,
		type: '',
		title: '发送短信',
		submitTxt: '',
	},
});
let ruleForm = ref({
	templateID: '',
	sig: '',
	data: [{ label: '', value: '' }],
  id: null,
});
const formRef = ref();
const rules = reactive<FormRules>({
	templateID: [
		{
			required: true,
			message: '请输入模版ID',
			trigger: 'blur',
		},
	],
	sig: [
		{
			required: true,
			message: '请输入短信签名',
			trigger: 'blur',
		},
	],
	id: [
		{
			required: true,
			message: '请输入用户ID',
			trigger: 'blur',
		},
	],
});
const closeDialog = () => {
	dialog.value.isShowDialog = false;
};
const handleAdd = () => {
	ruleForm.value.data.push({ value: '', label: '' });
};
const handleDel = (index: number) => {
	ruleForm.value.data.splice(index, 1);
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useMessageAdminApi()
				.sendSms({ ...ruleForm.value, uid: ruleForm.value.uid })
				.then((res: any) => {
					if (res.code === "SUCCESS" && res.data.success) {
						ElMessage({
							type: 'success',
							message: '发送短信成功!',
						});
						dialog.value.isShowDialog = false;
					} else {
						ElMessage({
							type: 'error',
							message: '发送失败',
						});
					}
				});
		} else {
			return false;
		}
	});
};
const { dialog } = toRefs(state);
const openDialog = (row: any) => {
	ruleForm.value = {
		templateID: '',
		sig: '',
		data: [{ label: '', value: '' }],
		id: row.numberID || row.userID || row.id,
	};

	dialog.value.isShowDialog = true;
};

defineExpose({
	openDialog,
});
</script>

<style lang="scss" scoped>
.subtitle {
	color: var(--el-text-color-secondary);
	text-align: center;
	padding-bottom: 20px;
}
</style>
