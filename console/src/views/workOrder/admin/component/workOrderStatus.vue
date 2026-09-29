<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<!-- 剩余余额  -->
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="78px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="对话状态" prop="status">
						<el-select v-model="ruleForm.status" placeholder="请选择完成状态">
							<el-option label="等待用户回复" :value="1"></el-option>
							<el-option label="等待管理员回复" :value="2"></el-option>
							<el-option label="工单完成" :value="3"></el-option>
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
	</el-dialog>
</template>

<script setup lang="ts" name="invoiceDialog">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { adminProcessInvoiceTypes } from '/@/api/invoice/admin/types';
import { useAdminWorkOrderApi } from '/@/api/workOrder/admin';

const formRef = ref();
const emit = defineEmits(['refresh']);
const status = ref<number>(1);
let ruleForm = ref<adminProcessInvoiceTypes>({
	orderID: '',
	status: 1,
});
const rules = reactive<FormRules>({
	status: [
		{
			required: true,
			message: '请选择工单状态',
			trigger: 'blur',
		},
	],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '工单完成状态',
	submitTxt: '确定',
});

const openDialog = (row: { orderID: string; status: number }) => {
	if (!row.orderID) {
		ElMessage({
			type: 'error',
			message: '工单ID获取失败',
		});
		return;
	}
	status.value = row.status;
	ruleForm.value.orderID = row.orderID;
	dialog.isShowDialog = true;
};
//重置
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useAdminWorkOrderApi()
				.userOrderFinsh({ orderID: ruleForm.value.orderID,status:ruleForm.value.status })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('工单状态设置成功');
						closeDialog();
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
	padding: 10px 10px;
}
</style>
