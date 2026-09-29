<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="850px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="文件上传" prop="content">
						<multipleUpload :fileList="fileList" @fileSuccess="fileSuccess" />
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit()" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script setup lang="ts" name="invoiceDialog">
import { reactive, ref, onMounted } from 'vue';
import { type FormRules } from 'element-plus';
import multipleUpload from '/@/components/multipleUpload/index.vue';

const formRef = ref();
const emit = defineEmits(['sendSuccess']);
let ruleForm = ref<any>({
	filename: [] as any[],
	file: [] as any[],
});
const fileList = ref<any>([]);

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

const openDialog = (row?: any, fileLists?: any) => {
	dialog.isShowDialog = true;
	fileList.value = fileLists;
	ruleForm.value = JSON.parse(JSON.stringify(row));
};
const fileSuccess = (rawFile: any) => {
	ruleForm.value.filename = [];
	ruleForm.value.file = [];
	if (rawFile && rawFile.length > 0) {
		rawFile.map((item: any) => {
			ruleForm.value.filename.push(item.name);
			ruleForm.value.file.push(item.raw);
		});
		fileList.value = rawFile;
	}
};
//重置

const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = () => {
	emit('sendSuccess', ruleForm.value, fileList.value);
	closeDialog();
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
