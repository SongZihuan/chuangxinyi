<template>
	<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px" @close="closeDialog">
		<div>
			<el-card shadow="never">
				<el-descriptions direction="vertical" :column="column">
					<el-descriptions-item label="钱包ID">{{ financeInfo.walletID || '-' }}</el-descriptions-item>
					<el-descriptions-item label="用户余额">
						<el-tag type="success" size="small">￥{{ formatAmount(financeInfo.balance) || '-' }}</el-tag>
					</el-descriptions-item>
          <el-descriptions-item label="已充值金额">
            <el-tag type="success" size="small">￥{{ formatAmount(financeInfo.cny) || '-' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="用户待入账余额">
            <el-tag type="success" size="small">￥{{ formatAmount(financeInfo.waitBalance) || '-' }}</el-tag>
          </el-descriptions-item>
					<el-descriptions-item label="未开票金额">
						<el-tag type="success" size="small">￥{{ formatAmount(financeInfo.notBilled) || '-' }}</el-tag>
					</el-descriptions-item>
					<el-descriptions-item label="已开票金额">
						<el-tag type="success" size="small">￥{{ formatAmount(financeInfo.hasBilled) || '-' }}</el-tag>
					</el-descriptions-item>
					<el-descriptions-item label="可开票金额">
						<el-tag type="success" size="small">￥{{ formatAmount(financeInfo.billed) || '-' }}</el-tag>
					</el-descriptions-item>
          <el-descriptions-item label="未提现金额">
            <el-tag type="success" size="small">￥{{ formatAmount(financeInfo.notWithdraw) || '-' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="已提现金额">
            <el-tag type="success" size="small">￥{{ formatAmount(financeInfo.hasWithdraw) || '-' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="可提现金额">
            <el-tag type="success" size="small">￥{{ formatAmount(financeInfo.withdraw) || '-' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="代入账可提现金额">
            <el-tag type="success" size="small">￥{{ formatAmount(financeInfo.waitWithdraw) || '-' }}</el-tag>
          </el-descriptions-item>
					<el-descriptions-item label="发票抬头">{{ financeInfo.titleName || '-' }}</el-descriptions-item>
					<el-descriptions-item label="税号">{{ financeInfo.titleTaxID || '-' }}</el-descriptions-item>
					<el-descriptions-item label="开户行账号">{{ financeInfo.titleBankID || '-' }}</el-descriptions-item>
					<el-descriptions-item label="开户行">{{ financeInfo.titleBank || '-' }}</el-descriptions-item>
				</el-descriptions>
			</el-card>
		</div>
	</el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { userWalletTypes } from '/@/views/system/user/types';
import { userApi } from '/@/api/system/user';
import { ElMessage } from 'element-plus';
import { formatAmount } from '/@/utils/formatAmount';
const column = ref(1);
const financeInfo = ref<userWalletTypes>({
	walletID: 0,
	balance: 0,
	notBilled: 0,
	billed: 0,
	hasBilled: 0,
});
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '用户资金信息',
		submitTxt: '',
	},
	ruleForm: {} as any,
});
const openDialog = (row?: any) => {
	state.dialog.isShowDialog = true;
	if (!row.id) {
		ElMessage.error('用户ID不能为空');
		return;
	}
	state.ruleForm = JSON.parse(JSON.stringify(row));
	getRealNameInfo();
};
// 关闭弹窗
const closeDialog = () => {
	state.dialog.isShowDialog = false;
};
const getRealNameInfo = () => {
	userApi()
		.getUserFinance({ uid: state.ruleForm.id })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				financeInfo.value = res.data;
			}
		});
};
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style lang="scss" scoped></style>
