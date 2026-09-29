<template>
	<div>
		<!--查看详情-->
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="600px" destroy-on-close @close="closeDialog">
			<el-row :gutter="35">
				<el-form>
					<el-col :span="24">
						<el-form-item label="发票ID">
							<span>{{ invoiceInfo.invoiceID }}</span>
						</el-form-item>
					</el-col>
					<el-col :span="24">
						<el-form-item label="发票类型">
							<span>{{ formatInvoiceType(invoiceInfo.type) }}</span>
						</el-form-item>
					</el-col>
					<el-col :span="24">
						<el-form-item label="发票抬头">
							<span>{{ invoiceInfo.title.name }}</span>
						</el-form-item>
					</el-col>
					<el-col :span="24">
						<el-form-item label="发票地址">
							<span>{{ invoiceInfo.address.address }}</span>
						</el-form-item>
					</el-col>
					<el-col :span="24">
						<el-form-item label="发票金额">
							<span
								><el-tag type="success">￥{{ formatAmount(invoiceInfo.amount) }}</el-tag></span
							>
						</el-form-item>
					</el-col>
					<el-col :span="24">
						<el-form-item label="发票状态">
							<el-tag type="primary" v-if="invoiceInfo.status === 1">待开票</el-tag>
							<el-tag type="success" v-else-if="invoiceInfo.status === 2">已开票</el-tag>
							<el-tag type="danger" v-else-if="invoiceInfo.status === 3">已退票</el-tag>
							<el-tag type="danger" v-else-if="invoiceInfo.status === 4">信息错误</el-tag>
							<el-tag type="danger" v-else>失败</el-tag>
						</el-form-item>
					</el-col>
					<el-col :span="24">
						<el-form-item label="添加时间">
							<span>{{ invoiceInfo.createAt ? dayjs.unix(invoiceInfo.createAt).format('YYYY-MM-DD HH:mm:ss') : '' }}</span>
						</el-form-item>
					</el-col>
					<el-col :span="24">
						<el-form-item label="开票时间">
							<span>{{ invoiceInfo.billingAt ? dayjs.unix(invoiceInfo.billingAt).format('YYYY-MM-DD HH:mm:ss') : '' }}</span>
						</el-form-item>
					</el-col>
					<el-col :span="24">
						<el-form-item label="退票时间">
							<span>{{ invoiceInfo.returnAt ? dayjs.unix(invoiceInfo.returnAt).format('YYYY-MM-DD HH:mm:ss') : '' }}</span>
						</el-form-item>
					</el-col>
				</el-form>
			</el-row>
		</el-dialog>
	</div>
</template>
<script setup lang="ts" name="invoiceInfo">
import { reactive, ref } from 'vue';
import { formatAmount, formatInvoiceType } from '/@/utils/formatAmount';
import dayjs from 'dayjs';
import { invoiceAdminDataType } from '/@/views/invoice/user/types';
import { ElMessage } from 'element-plus';
import { useInvoiceAdminApi } from '/@/api/invoice/admin';

const invoiceInfo = ref<invoiceAdminDataType>({
	invoiceID: '',
	type: 1,
	title: {
		name: '',
		taxID: '',
		bandID: '',
		band: '',
	},
	address: {
		address: '',
		phone: '',
		province: '',
		email: '',
		city: '',
		district: '',
		name: '',
	},
	amount: 0,
	status: 1,
	createAt: 0,
	billingAt: 0,
	returnAt: 0,
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'look',
	title: '发票详情',
	submitTxt: '查看',
});
const openDialog = (row: { invoiceID: string }) => {
	getInvoiceInfo(row.invoiceID);
};
// 获取详情
const getInvoiceInfo = async (id: string) => {
	await useInvoiceAdminApi()
		.getInvoiceInfo({ id: id })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				invoiceInfo.value = res.data.invoice;
				dialog.isShowDialog = true;
			} else {
				ElMessage.error(res.msg);
			}
		});
};
//重置
const closeDialog = () => {
	dialog.isShowDialog = false;
};
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style scoped lang="scss"></style>
