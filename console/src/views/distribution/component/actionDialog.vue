<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="400px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="100px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-text class="mx-1" type="warning">分销层级相同会进行覆盖操作</el-text>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="分销层级" prop="level">
						<el-input-number
							v-model="ruleForm.level"
							placeholder="请输入分销层级"
							clearable
							controls-position="right"
							:disabled="dialog.type == 'edit'"
						></el-input-number>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" label="返佣比例">
					<el-form-item label="返佣比例(%)" prop="pre"
						><el-input-number v-model="ruleForm.pre" placeholder="请输入返佣比例" clearable controls-position="right"></el-input-number
					></el-form-item>
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
import { useDistributionApi } from '/@/api/distribution';
import { fromType } from '/@/api/distribution/types';

const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<fromType>({
	level: 1,
	pre: 1,
});
const rules = reactive<FormRules>({
	level: [{ required: true, message: '请输入分销层级', trigger: 'blur' }],
	pre: [{ required: true, message: '请输入分销比例', trigger: 'blur' }],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '分销新增',
	submitTxt: '新增',
});

const openDialog = (type: string, row: any) => {
	if (type === 'add') {
		reset();
		dialog.type = 'add';
		dialog.title = '分销新增';
		dialog.submitTxt = '新增';
	} else {
		dialog.type = 'edit';
		dialog.title = '分销编辑';
		dialog.submitTxt = '编辑';
		ruleForm.value = JSON.parse(JSON.stringify(row));
	}
	dialog.isShowDialog = true;
};
const reset = () => {
	ruleForm.value.level = 1;
	ruleForm.value.pre = 1;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};

const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useDistributionApi()
				.distributionUpdate(ruleForm.value)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '操作成功',
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

<style lang="scss" scoped></style>
