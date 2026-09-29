<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="850px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="工单内容" prop="content">
						<el-input v-model="ruleForm.content" type="textarea" :autosize="{ minRows: 2, maxRows: 6 }" placeholder="请输入内容"></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="文件上传" prop="content">
						<multipleUpload @fileSuccess="fileSuccess" />
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

<script setup lang="ts" name="invoiceDialog">
import { reactive, ref, onMounted } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import multipleUpload from '/@/components/multipleUpload/index.vue';
import { useUserWorkOrderApi } from '/@/api/workOrder/user';
import { encryptFileName } from '/@/utils/getFileType';

const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<any>({
	title: '',
	content: '',
	filename: [] as any[],
	file: [] as any[],
});

const rules = reactive<FormRules>({
	title: [
		{
			required: true,
			message: '请输入工单标题',
			trigger: 'blur',
		},
	],
	content: [
		{
			required: true,
			message: '请输入工单内容',
			trigger: 'blur',
		},
	],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'replay',
	title: '回复工单',
	submitTxt: '新增',
});

const openDialog = (row?: any) => {
	reset();
	dialog.isShowDialog = true;
	ruleForm.value = JSON.parse(JSON.stringify(row));
};
const fileSuccess = (rawFile: any) => {
	ruleForm.value.filename = [];
	ruleForm.value.file = [];
	if (rawFile && rawFile.length > 0) {
		rawFile.map((item: any) => {
			ruleForm.value.filename.push(encryptFileName(item.name));
			ruleForm.value.file.push(item.raw);
		});
	}
};
//重置
const reset = () => {
	ruleForm.value = {
		title: '',
		content: '',
		filename: [],
		file: [],
	};
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (ruleForm.value.orderID) {
		ruleForm.value.id = ruleForm.value.orderID;
	} else {
		ElMessage({
			type: 'error',
			message: '工单ID不存在',
		});
	}
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useUserWorkOrderApi()
				.userReplyWorkOrder(ruleForm.value)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '新增工单成功',
						});
						emit('refresh');
					}
				});
		} else {
			return false;
		}
	});
};

onMounted(() => {});
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
	padding: 10px 10px;
}
</style>
