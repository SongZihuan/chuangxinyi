<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<!-- 剩余余额  -->
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="78px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="开票状态" prop="status">
						<el-select v-model="ruleForm.status" placeholder="请选择开票类型">
							<el-option v-if="status === 1" label="待开票" :value="1"></el-option>
							<el-option v-if="status === 1 || status === 2" label="已开票" :value="2"></el-option>
							<el-option v-if="status === 1 || status === 2 || status === 3 || status === 4" label="已退票" :value="3"></el-option>
							<el-option v-if="status === 1 || status === 2 || status === 3 || status === 4" label="信息错误" :value="4"></el-option>
						</el-select>
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
		<checkPhone ref="checkPhoneRef" @checkPhoneSuccess="checkPhoneSuccess" />
	</el-dialog>
</template>

<script setup lang="ts" name="invoiceDialog">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import checkPhone from '/@/components/checkPhone/index.vue';
import { ElMessage } from 'element-plus';
import { adminProcessInvoiceTypes } from '/@/api/invoice/admin/types';
import { useInvoiceAdminApi } from '/@/api/invoice/admin';

const formRef = ref();
const emit = defineEmits(['refresh']);
const status = ref<number>(1);
const upStatus = ref<number>(1);
const checkPhoneRef = ref();
let ruleForm = ref<adminProcessInvoiceTypes>({
	id: '',
	status: 1,
});
const rules = reactive<FormRules>({
	status: [
		{
			required: true,
			message: '请选择开票类型',
			trigger: 'blur',
		},
	],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '开票',
	submitTxt: '开票',
});

const openDialog = (row: { invoiceID: string; status: number }) => {
	if (!row.invoiceID) {
		ElMessage({
			type: 'error',
			message: '发票ID获取失败',
		});
		return;
	}
	status.value = Number(row.status);
	upStatus.value = Number(row.status);
	ruleForm.value.id = row.invoiceID;
	ruleForm.value.status = row.status;
	dialog.isShowDialog = true;
};
//重置
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (upStatus.value === ruleForm.value.status) {
		closeDialog();
		return;
	}
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			checkPhoneRef.value.openDialog({}, '开票确认');
		} else {
			return false;
		}
	});
};
const checkPhoneSuccess = (phoneToken: string) => {
	useInvoiceAdminApi()
		.processInvoice({ ...ruleForm.value, phoneToken: phoneToken })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				closeDialog();
				ElMessage({
					type: 'success',
					message: '开票成功',
				});
				emit('refresh');
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
	padding: 10px 10px;
}
</style>
