<template>
	<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px" @close="closeDialog">
		<div>
			<el-card shadow="never">
				<el-descriptions direction="vertical" :column="column">
					<el-descriptions-item label="订单ID">{{ info.tradeID || '-' }}</el-descriptions-item>
					<el-descriptions-item label="订单名称">{{ info.subject || '-' }}</el-descriptions-item>
					<el-descriptions-item label="订单金额"> {{ formatAmount(info.cny) }}元 </el-descriptions-item>
					<el-descriptions-item label="企业ID">{{ info.companyID || '-' }}</el-descriptions-item>
					<el-descriptions-item label="支付状态">
						<el-tag type="primary" v-if="info.tradeStatus === 1 || info.tradeStatus === 6">待支付</el-tag>
						<el-tag type="success" v-else-if="info.tradeStatus === 2 || info.tradeStatus === 3">成功</el-tag>
						<el-tag type="danger" v-else-if="info.tradeStatus === 4">支付关闭</el-tag>
						<el-tag type="danger" v-else-if="info.tradeStatus === 5">支付退款</el-tag>
						<el-tag type="danger" v-else>失败</el-tag>
					</el-descriptions-item>
					<el-descriptions-item label="支付方式">{{ info.payWay || '-' }}</el-descriptions-item>
					<el-descriptions-item label="支付时间">{{ dayjs.unix(info.payAt).format('YYYY-MM-DD HH:mm:ss') || '-' }}</el-descriptions-item>
					<el-descriptions-item label="创建时间">{{ dayjs.unix(info.createAt).format('YYYY-MM-DD HH:mm:ss') || '-' }}</el-descriptions-item>
				</el-descriptions>
			</el-card>
		</div>
	</el-dialog>
</template>
<script setup lang="ts">
import { ref, reactive } from 'vue';
import { tradeTypes } from '/@/views/system/user/types';
import { ElMessage } from 'element-plus';
import { formatAmount } from '/@/utils/formatAmount';
import dayjs from 'dayjs';
const column = ref(1);
const info = ref<tradeTypes>({
	tradeID: '',
	subject: '',
	cny: '',
	companyID: '',
	tradeStatus: 0,
	payWay: '',
	createAt: '',
	payAt: '',
});
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '支付信息',
		submitTxt: '',
	},
});
const openDialog = (row?: any) => {
	state.dialog.isShowDialog = true;
	if (!row) {
		ElMessage.error('支付信息不存在');
		return;
	}
	info.value = JSON.parse(JSON.stringify(row));
};
// 关闭弹窗
const closeDialog = () => {
	state.dialog.isShowDialog = false;
};
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style lang="scss" scoped></style>
