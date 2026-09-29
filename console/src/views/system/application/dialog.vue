<template>
	<div class="system-menu-dialog-container">
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" destroy-on-close @close="closeDialog">
			<el-form ref="menuDialogFormRef" :model="ruleForm" size="default" label-width="90px" :rules="rules">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="应用名称" prop="name">
							<el-input v-model="ruleForm.name" placeholder="请输入应用名称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="应用图标" prop="icon">
							<IconSelector placeholder="请输入应用图标" v-model="ruleForm.icon" />
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="网站" prop="webID">
							<el-select v-model="ruleForm.webID" placeholder="请选择" clearable class="w100">
								<el-option :label="item.name" :value="item.id" v-for="(item, index) in webIDDict" :key="index"></el-option>
							</el-select>
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
					<el-col :span="24" class="mb20">
						<el-form-item label="应用地址" prop="url">
							<el-input v-model="ruleForm.url" placeholder="应用地址" clearable @input="onVerifyUrl($event)"> </el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="描述" prop="describe">
							<el-input type="textarea" v-model="ruleForm.describe" placeholder="请输入描述" clearable :rows="4"></el-input>
						</el-form-item>
					</el-col>
				</el-row>
			</el-form>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit(menuDialogFormRef)" size="default">{{ dialog.submitTxt }}</el-button>
				</span>
			</template>
		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="systemMenuDialog">
import { defineAsyncComponent, reactive, ref, nextTick, onMounted } from 'vue';
import { useApplicationApi } from '/@/api/system/application/index';
import { verifyUrl, verifyAndSpace } from '/@/utils/toolsValidate';
import { type FormRules, type FormInstance } from 'element-plus';
import { message } from '/@/utils/message';
const emit = defineEmits(['refresh']);
const IconSelector = defineAsyncComponent(() => import('/@/components/iconSelector/index.vue'));
const menuDialogFormRef = ref();
let ruleForm = ref({
	name: '',
	webID: '',
	url: '',
	icon: '',
	status: 1,
	describe: '',
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '',
	submitTxt: '',
});
const webIDDict = ref<any>([]);
const rules = reactive<FormRules>({
	name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
	icon: [{ required: true, message: '请输入应用图标', trigger: 'blur' }],
	webID: [{ required: true, message: '请选择网站', trigger: 'blur' }],
	url: [{ required: true, message: '请输入地址', trigger: 'blur' }],
	status: [{ required: true, message: '请选择状态', trigger: 'blur' }],
});
const applicationApi = useApplicationApi();
const onVerifyUrl = (val: string) => {
	ruleForm.value.url = verifyAndSpace(val);
};
const openDialog = (type: string, row?: any) => {
	dialog.isShowDialog = true;
	dialog.type = type;
	if (type == 'edit') {
		nextTick(() => {
			ruleForm.value = JSON.parse(JSON.stringify(row));
		});
		dialog.title = '修改菜单';
		dialog.submitTxt = '修 改';
	} else {
		reset();

		dialog.title = '新增菜单';
		dialog.submitTxt = '新 增';
	}
};
const reset = () => {
	ruleForm.value = {
		name: '',
		webID: '',
		url: '',
		icon: '',
		status: 1,
		describe: '',
	};
};

const closeDialog = () => {
	menuDialogFormRef.value?.resetFields();
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
			if (!verifyUrl(ruleForm.value.url)) return message('url格式不挣正确', { type: 'warning' });
			if (dialog.type === 'add') {
				applicationApi.addAdminapplication({ ...ruleForm.value }).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('新增成功', { type: 'success' });
						closeDialog(); // 关闭弹窗
						emit('refresh');
					}
				});
			} else if (dialog.type === 'edit') {
				applicationApi.applicationUpdate({ ...ruleForm.value }).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('编辑成功', { type: 'success' });
						closeDialog(); // 关闭弹窗
						emit('refresh');
					}
				});
			}
		} else {
			return false;
		}
	});
};
const getSubAll = () => {
	applicationApi.websiteAll().then((res: any) => {
		if (res.code === "SUCCESS") {
			webIDDict.value = res.data.website;
		}
	});
};
onMounted(() => {
	getSubAll();
});
// 暴露变量
defineExpose({
	openDialog,
});
</script>
