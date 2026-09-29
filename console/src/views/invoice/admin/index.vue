<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
						<el-select v-model="state.tableData.param.timetype" placeholder="请选择日期类型" style="width: 100%">
							<el-option v-for="(item, index) in dateTypeDict" :value="item.value" :label="item.label" :key="index" />
						</el-select>
					</el-form-item>
					<el-form-item>
						<el-date-picker
							v-model="state.tableData.param.range"
							type="datetimerange"
							range-separator="至"
							start-placeholder="开始日期"
							end-placeholder="结束日期"
							style="width: 100%"
						/>
					</el-form-item>
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
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="invoiceID" label="发票ID" show-overflow-tooltip min-width="300" align="left"></el-table-column>
				<el-table-column prop="type" label="发票类型" align="center" min-width="160">
					<template #default="scope">
						<el-tag type="primary">{{ formatInvoiceType(scope.row.type) }}</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="title.name" label="发票抬头" min-width="150" align="center"></el-table-column>
				<el-table-column prop="address.address" label="地址" min-width="200" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="amount" label="金额" align="center" min-width="100">
					<template #default="scope">
						<el-tag type="success">￥{{ formatAmount(scope.row.amount) }}</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="status" label="状态" align="center" min-width="100">
					<template #default="scope">
						<el-tag type="primary" v-if="scope.row.status === 1">待开票</el-tag>
						<el-tag type="success" v-else-if="scope.row.status === 2">已开票</el-tag>
						<el-tag type="danger" v-else-if="scope.row.status === 3">已退票</el-tag>
						<el-tag type="danger" v-else-if="scope.row.status === 4">信息错误</el-tag>
						<el-tag type="danger" v-else>失败</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="createAt" label="添加时间" show-overflow-tooltip width="160" align="center">
					<template #default="scope">
						<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="billingAt" label="开票时间" show-overflow-tooltip width="160" align="center">
					<template #default="scope">
						<span v-if="scope.row.billingAt">{{ dayjs.unix(scope.row.billingAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="returnAt" label="退票时间" show-overflow-tooltip width="160" align="center">
					<template #default="scope">
						<span v-if="scope.row.returnAt">{{ dayjs.unix(scope.row.returnAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
            <span v-else>-</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="100" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="openInfo(scope.row)">查看</el-button>
						<el-button v-if="!(scope.row.status === 3) || !(scope.row.status === 4)" text type="primary" @click="openStatus(scope.row)"
							>开票</el-button
						>
					</template>
				</el-table-column>
			</el-table>
			<el-pagination hide-on-single-page
				@size-change="onHandleSizeChange"
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				:page-sizes="[10, 20, 30]"
				v-model:current-page="state.tableData.param.page"
				background
				v-model:page-size="state.tableData.param.pagesize"
				layout="total, sizes, prev, pager, next, jumper"
				:total="state.tableData.total"
			>
			</el-pagination>
		</el-card>
		<invoice-dialog ref="invoiceDialogRef" @refresh="getTableData"></invoice-dialog>
		<invoice-info ref="invoiceInfoRef"></invoice-info>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { invoiceAdminState } from './types';
import { useInvoiceAdminApi } from '/@/api/invoice/admin';
import { formatAmount, formatInvoiceType } from '/@/utils/formatAmount';
import InvoiceInfo from '/@/views/invoice/admin/component/invoiceInfo.vue';
import InvoiceDialog from '/@/views/invoice/admin/component/invoiceDialog.vue';

const invoiceInfoRef = ref();
const invoiceDialogRef = ref();
const dateTypeDict = ref([
	{ label: '创建时间', value: 1 },
	{ label: '开票时间', value: 4 },
	{ label: '退票时间', value: 5 },
]);
const state = reactive<invoiceAdminState>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			timetype: 1,
			starttime: 0,
			endtime: 0,
			range: [],
		},
	},
});
watch(
	() => state.tableData.param.range,
	() => {
		if (state.tableData.param.range.length > 0) {
			state.tableData.param.starttime = dayjs(state.tableData.param.range[0]).unix();
			state.tableData.param.endtime = dayjs(state.tableData.param.range[1]).unix();
		}
	},
	{
		deep: true,
	}
);
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	useInvoiceAdminApi()
		.getInvoiceList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.invoice;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
//搜索
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
// 查看详情
const openInfo = (row: { invoiceID: string }) => {
	invoiceInfoRef.value.openDialog(row);
};
// 申请开票
const openStatus = (row: { invoiceID: string; status: number }) => {
	invoiceDialogRef.value.openDialog(row);
};
// 分页改变
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
	getTableData();
};
// 页面加载时
onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss"></style>
