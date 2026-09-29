<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="800px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="标题" prop="title">
						<el-input v-model="ruleForm.title" placeholder="请输入标题" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="开始时间" prop="startAt">
						<el-date-picker v-model="ruleForm.startAt" type="date" placeholder="选择日期时间" clearable style="width: 100%"></el-date-picker>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="结束时间" prop="stopAt">
						<el-date-picker v-model="ruleForm.stopAt" type="date" placeholder="选择日期时间" clearable style="width: 100%"></el-date-picker>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="公告内容" prop="content">
						<WangEditor v-model:get-html="ruleForm.content" />
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
import WangEditor from '/@/components/Editor/index.vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { announcementCreateTypes } from '/@/api/announcement/types';
import { useAnnouncementApi } from '/@/api/announcement';
import dayjs from 'dayjs';
const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<announcementCreateTypes>({
	title: '',
	content: '',
	startAt: 0,
	stopAt: 0,
});
let allCheck = ref<any>([]);
const rules = reactive<FormRules>({
	title: [
		{ required: true, message: '请输入标题', trigger: 'blur' },
		{ min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' },
	],
	startAt: [{ required: true, message: '请选择开始时间', trigger: 'blur' }],
	content: [{ required: true, message: '请输入公告内容', trigger: 'blur' }],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '公告新增',
	submitTxt: '新增',
});

const openDialog = (type: string, row: any) => {
	if (type === 'add') {
		reset();
		dialog.title = '公告新增';
		dialog.submitTxt = '新增';
	} else {
		ruleForm.value = JSON.parse(JSON.stringify(row));
		ruleForm.value.startAt = dayjs.unix(ruleForm.value.startAt);
		ruleForm.value.stopAt = dayjs.unix(ruleForm.value.stopAt);
		dialog.title = '公告编辑';
		dialog.submitTxt = '编辑';
	}
	dialog.type = type;
	dialog.isShowDialog = true;
};
//重置
const reset = () => {
	allCheck.value = false;
	ruleForm.value = {
		title: '',
		content: '',
		startAt: 0,
		stopAt: 0,
	};
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			let data = { ...ruleForm.value, startAt: dayjs(ruleForm.value.startAt).unix(), stopAt: dayjs(ruleForm.value.stopAt).unix() };
			if (dialog.type === 'add') {
				useAnnouncementApi()
					.createAnnouncement(data)
					.then((res: any) => {
						if (res.code === "SUCCESS") {
							closeDialog();
							ElMessage({
								type: 'success',
								message: '新增公告成功',
							});
							emit('refresh');
						}
					});
			} else {
				useAnnouncementApi()
					.updateAnnouncement(data)
					.then((res: any) => {
						if (res.code === "SUCCESS") {
							closeDialog();
							ElMessage({
								type: 'success',
								message: '编辑公告成功',
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
