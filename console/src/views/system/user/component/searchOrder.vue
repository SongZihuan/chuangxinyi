<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="900px" destroy-on-close @close="closeDialog">
		<div class="search-header mb15">
			<el-form :inline="true" :model="param">
				<el-form-item> <el-input placeholder="请输入订单号" style="width: 300px" v-model="param.tradeID" clearable> </el-input></el-form-item>
				<el-form-item>
					<el-button type="primary" @click="handleSearch">
						<el-icon>
							<ele-Search />
						</el-icon>
						查询
					</el-button>
				</el-form-item>
			</el-form>
		</div>
		<el-descriptions title="订单详情" v-if="info.tradeID" :column="column" class="dialog-content">
			<el-descriptions-item label="订单ID">{{ info.tradeID }}</el-descriptions-item>
			<el-descriptions-item label="订单名称">{{ info.subject }}</el-descriptions-item>
			<el-descriptions-item label="订单金额">{{ formatAmount(info.cny) }}</el-descriptions-item>
			<el-descriptions-item label="支付状态">
				<el-tag type="primary" v-if="info.tradeStatus === 1 || info.tradeStatus === 6">待支付</el-tag>
				<el-tag type="success" v-else-if="info.tradeStatus === 2 || info.tradeStatus === 3">成功</el-tag>
				<el-tag type="danger" v-else-if="info.tradeStatus === 4">支付关闭</el-tag>
				<el-tag type="danger" v-else-if="info.tradeStatus === 5">支付退款</el-tag>
				<el-tag type="danger" v-else>失败</el-tag>
			</el-descriptions-item>
			<el-descriptions-item label="支付方式">{{ info.payWay }}</el-descriptions-item>
			<el-descriptions-item label="支付时间">{{ dayjs.unix(info.payAt).format('YYYY-MM-DD HH:mm:ss') }}</el-descriptions-item>
			<el-descriptions-item label="创建时间">{{ dayjs.unix(info.createAt).format('YYYY-MM-DD HH:mm:ss') }}</el-descriptions-item>
		</el-descriptions>
		<el-empty description="暂无数据" v-else></el-empty>
	</el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import dayjs from 'dayjs';
import { formatAmount } from '/@/utils/formatAmount';
import { ElMessage } from 'element-plus';
import { useDefraytApi } from '/@/api/defray';
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '充值订单查询',
	submitTxt: '',
});
const column = ref(3);
const info = ref({
	tradeID: '',
	subject: '',
	cny: '',
	tradeStatus: 1,
	payAt: '',
	payWay: '',
	createAt: '',
});
const param = reactive({
	tradeID: '',
});

const handleSearch = () => {
	getTableData();
};
const getTableData = () => {
	if (!param.tradeID) {
		ElMessage({
			type: 'warning',
			message: '请输入订单号',
		});
		return;
	}
	if (dialog.type === 'pay') {
		searchPay();
	} else {
		searchDefray();
	}
};
const searchPay = async () => {
	await useDefraytApi()
		.adminPayInfo({ id: param.tradeID })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				info.value = res.data.pay;
			}
		});
};
const searchDefray = () => {
	useDefraytApi()
		.adminDefrayInfo({ id: param.tradeID })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				info.value = res.data.defray;
			}
		});
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};

const openDialog = async (type: string) => {
	if (type === 'pay') {
		dialog.title = '充值订单查询';
	} else {
		dialog.title = '支付订单查询';
	}
	dialog.type = type;
	dialog.isShowDialog = true;
};
onMounted(() => {
	const clientWidth = document.body.clientWidth;
	if (clientWidth < 1000) {
		column.value = 1;
	}
});
// 暴露变量
defineExpose({
	openDialog,
});
</script>
<style lang="scss" scoped>
.dialog-content {
	height: 293px;
}
</style>
