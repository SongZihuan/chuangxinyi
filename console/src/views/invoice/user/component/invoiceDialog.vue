<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<!-- 剩余余额  -->
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="90px" :rules="rules">
			<el-row :gutter="35">
				<el-col :span="24" class="ml5 mb10">
					<div class="tip">
						<el-text type="success" tag="b">剩余开票金额：￥{{ formatAmount(balance.notBilled) }}</el-text>
					</div>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="开票金额" prop="amount">
						<el-input-number
							v-model="ruleForm.amount"
							controls-position="right"
							style="width: 100%"
							:precision="2"
							:step="0.01"
							:max="balance.notBilled"
							:min="balance.notBilled < 0.01 ? balance.notBilled : 0.01"
						/>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="开票类型" prop="type">
						<el-select v-model="ruleForm.type" placeholder="请选择开票类型" style="width: 100%">
							<el-option label="增值税普通发票（个人）" :value="1"></el-option>
							<el-option label="增值税普通发票（企业）" :value="2"></el-option>
							<el-option label="增值税专用发票（企业）" :value="3"></el-option>
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
import { reactive, ref, onMounted } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { invoiceTypes } from '/@/api/invoice/user/types';
import { useInvoiceUserApi } from '/@/api/invoice/user';
import { formatAmount } from '/@/utils/formatAmount';

interface balanceTypes {
	balance: number;
	billed: number;
	hasBilled: number;
	notBilled: number;
}
const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<invoiceTypes>({
	amount: 0,
	type: 1,
});
const balance = ref<balanceTypes>({
	balance: 0,
	billed: 0,
	hasBilled: 0,
	notBilled: 0,
});
const rules = reactive<FormRules>({
	amount: [
		{
			required: true,
			message: '请输入开票金额',
			trigger: 'blur',
		},
	],
	type: [
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
	title: '申请开票',
	submitTxt: '申请',
});

const openDialog = (row: any) => {
	reset();
	balance.value = row;
	dialog.isShowDialog = true;
};
//重置
const reset = () => {
	ruleForm.value = {
		amount: 0,
		type: 1,
	};
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			if (ruleForm.value.amount > balance.value.billed) {
				ElMessage({
					type: 'error',
					message: '开票金额不能大于剩余开票余额',
				});
				return false;
			}
			if (ruleForm.value.amount < 0.01) {
				ElMessage({
					type: 'error',
					message: '开票金额不能小于0.01',
				});
				return false;
			}
			let data = {
				amount: ruleForm.value.amount,
				type: ruleForm.value.type,
			};
			if (ruleForm.value.amount) {
				data.amount = ruleForm.value.amount * 100 | 0;
			}
			useInvoiceUserApi()
				.Invoice(data)
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
