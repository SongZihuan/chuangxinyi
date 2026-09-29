<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-text type="success" tag="b">剩余开票金额：￥{{ formatAmount(balance.notBilled) }}</el-text>
			<div class="search-header mb15 mt15">
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
						<el-button type="success" @click="openDialog()">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							申请开票
						</el-button>
						<el-button type="success" @click="openDialog('editInvoice')">
							<el-icon>
								<ele-Edit />
							</el-icon>
							抬头
						</el-button>
						<el-button type="success" @click="openDialog('editAddress')" class="max-sm:mt-[10px] ml-0">
							<el-icon>
								<ele-Edit />
							</el-icon>
							地址
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="invoiceID" label="发票ID" show-overflow-tooltip min-width="300" align="left"></el-table-column>
				<el-table-column prop="type" label="发票类型" align="center" min-width="120">
					<template #default="scope">
						<el-tag type="primary">{{ formatInvoiceType(scope.row.type) }}</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="title.name" label="发票抬头" min-width="150" align="center" show-overflow-tooltip></el-table-column>
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
				<el-table-column prop="createAt" label="申请时间" show-overflow-tooltip width="160" align="center">
					<template #default="scope">
						<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="billingAt" label="开票时间" show-overflow-tooltip width="160" align="center">
					<template #default="scope">
						<span v-if="scope.row.billingAt">{{ dayjs.unix(scope.row.billingAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
						<span v-else>-</span>
					</template>
				</el-table-column>
				<el-table-column prop="returnAt" label="退票时间" show-overflow-tooltip width="160" align="center">
					<template #default="scope">
						<span v-if="scope.row.returnAt">{{ dayjs.unix(scope.row.returnAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
						<span v-else>-</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="120" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="openInfo(scope.row)">查看</el-button>
						<!--            下载-->
						<el-button text type="primary" @click="download(scope.row)">下载</el-button>
					</template>
				</el-table-column>
			</el-table>
			<el-pagination
				hide-on-single-page
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
		<invoice-dialog ref="invoiceRef" @refresh="getTableData"></invoice-dialog>
		<invoice-info ref="invoiceInfoRef"></invoice-info>
		<!-- 发票抬头 -->
		<invoice ref="invoiceEditRef" :invoice-form="invoiceForm"></invoice>
		<!-- 地址 -->
		<EditAddress ref="editAddressRef" :userInfo="addressInfo" @refresh="getUserInfo" />
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { invoiceAdminState } from './types';
import { useInvoiceUserApi } from '/@/api/invoice/user';
import InvoiceDialog from '/@/views/invoice/user/component/invoiceDialog.vue';
import InvoiceInfo from '/@/views/invoice/user/component/invoiceInfo.vue';
import { formatAmount, formatInvoiceType } from '/@/utils/formatAmount';
import invoice from '/@/views/userCenter/countInfo/component/invoice.vue';
import EditAddress from '/@/views/userCenter/userInfo/component/edit.vue';
import { Session } from '/@/utils/storage';
import { useUserApi } from '/@/api/user/user';
const useUserApiCollect = useUserApi();
const dateTypeDict = ref([
	{ label: '创建时间', value: 1 },
	{ label: '开票时间', value: 4 },
	{ label: '退票时间', value: 5 },
]);
const balance = ref({
	balance: '',
	notBilled: '',
});
const invoiceRef = ref();
const invoiceEditRef = ref();
const invoiceInfoRef = ref();
const editAddressRef = ref();
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
const addressInfo = ref({
	name: '',
	phone: '',
	email: '',
	province: '',
	city: '',
	district: '',
	address: '',
	areas: [],
});
const invoiceForm = ref({
	name: '',
	tax_id: '',
	band_id: '',
	band: '',
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
const getTableData = async () => {
	state.tableData.loading = true;
	await useInvoiceUserApi()
		.getInvoiceList(state.tableData.param)
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				state.tableData.data = res.data.invoice;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};

// 查看详情
const openInfo = (row: { invoiceID: string }) => {
	invoiceInfoRef.value.openDialog(row);
};
// 下载
const download = (row: { invoiceID: string }) => {
	useInvoiceUserApi()
		.downloadInvoice({ invoiceID: row.invoiceID })
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				window.open(res.data.url);
			}
		});
};
const getUserInfo = async () => {
	await useUserApiCollect.userInfo().then((res: any) => {
		if (res.code === 'SUCCESS') {
			addressInfo.value = res.data.address;
			Session.set('userData', res.data);
			if (res.data.balance) {
				balance.value = res.data.balance;
			}
			if (res.data?.title) {
				invoiceForm.value.name = res.data.title.name;
				invoiceForm.value.tax_id = res.data.title.taxID;
				invoiceForm.value.band_id = res.data.title.bandID;
				invoiceForm.value.band = res.data.title.band;
			}
		}
	});
};
//搜索
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
// 申请开票
const openDialog = (type?: string) => {
	if (type == 'editInvoice') {
		invoiceEditRef.value.openDialog(balance.value);
		return;
	}
	if (type == 'editAddress') {
		editAddressRef.value.openDialog(balance.value);
		return;
	}
	invoiceRef.value.openDialog(balance.value);
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
onMounted(async () => {
	await getTableData();
	await getUserInfo();
});
</script>

<style scoped lang="scss"></style>
