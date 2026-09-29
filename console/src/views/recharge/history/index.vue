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
				<el-table-column prop="tradeID" label="订单ID" show-overflow-tooltip min-width="240" align="left"></el-table-column>
				<el-table-column prop="subject" label="订单名称" show-overflow-tooltip align="center" min-width="140"></el-table-column>
				<el-table-column prop="user.id" label="发起人ID" show-overflow-tooltip align="center" min-width="140"></el-table-column>
				<el-table-column prop="user.phone" label="发起人" show-overflow-tooltip align="center" min-width="140"></el-table-column>
				<el-table-column prop="tradeStatus" label="支付状态" show-overflow-tooltip align="center" width="100">
					<template #default="scope">
						<el-tag type="primary" v-if="scope.row.tradeStatus === 1">等待支付</el-tag>
						<el-tag type="success" v-else-if="scope.row.tradeStatus === 2 || scope.row.tradeStatus === 3">支付成功</el-tag>
						<el-tag type="danger" v-else-if="scope.row.tradeStatus === 4">支付关闭</el-tag>
						<el-tag type="info" v-else-if="scope.row.tradeStatus === 5">等待退款</el-tag>
						<el-tag type="success" v-else-if="scope.row.tradeStatus === 6">退款成功</el-tag>
						<el-tag type="info" v-else-if="scope.row.tradeStatus === 7">退款失败</el-tag>
						<el-tag type="danger" v-else>支付未找到</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="payWay" label="支付方式" show-overflow-tooltip align="center" width="130"></el-table-column>
				<el-table-column prop="cny" label="订单金额" show-overflow-tooltip align="center" width="130">
					<template #default="scope">
						<el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.cny) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column prop="get" label="实际获得金额" show-overflow-tooltip align="center" width="130">
					<template #default="scope">
						<el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.get) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column label="支付时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span v-if="scope.row.payAt">{{ dayjs.unix(scope.row.payAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="退款时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span v-if="scope.row.refundAt">{{ dayjs.unix(scope.row.refundAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="createTime" label="创建时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="150" fixed="right" align="center">
					<template #default="scope">
						<el-button
							text
							type="primary"
							@click="itemRefund(scope.row)"
							v-if="scope.row.tradeStatus === 2 || scope.row.tradeStatus === 3 || scope.row.tradeStatus === 7"
							>申请退款</el-button
						>
						<el-button text type="primary" @click="itemSearch(scope.row)" v-if="scope.row.tradeStatus === 1">支付查询</el-button>
						<el-button text type="primary" @click="searchRefund(scope.row)" v-if="scope.row.tradeStatus === 5">退款查询</el-button>
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
		<!-- 手机验证 -->
		<phoneDialog ref="phoneDialogRef" @refresh="getTableData()" />
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, watch, defineAsyncComponent } from 'vue';
import { rechargeApi } from '/@/api/recharge/index';
import { useBaseApi } from '/@/api/base/index';
import type { rechargeStateTypes } from '/@/api/recharge/types';
import { ElMessage } from 'element-plus';
import { formatAmount } from '/@/utils/formatAmount';
const phoneDialog = defineAsyncComponent(() => import('./component/phoneDialog.vue'));
import dayjs from 'dayjs';
const rechargeApiCollect = rechargeApi();
const useBaseApiCollect = useBaseApi();
const phoneDialogRef = ref();
const dateTypeDict = ref([
	{ label: '创建时间', value: 1 },
	{ label: '支付时间', value: 2 },
]);
const state = reactive<rechargeStateTypes>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			range: [],
			page: 1,
			pagesize: 20,
			timetype: '',
			starttime: null,
			endtime: null,
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
	rechargeApiCollect.getRechargeList(state.tableData.param).then((res: any) => {
		if (res.code === 'SUCCESS') {
			state.tableData.data = res.data.pay;
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
//查询这条记录是否支付成功
const itemSearch = (item: any) => {
	useBaseApiCollect.checkOrder({ tradeid: item.tradeID }).then((res: any) => {
		if (res.code === 'SUCCESS' && res.data.success) {
			ElMessage.success('支付成功');
		} else {
			ElMessage.warning('您暂未支付成功');
		}
	});
};
const searchRefund = (item: any) => {
	useBaseApiCollect.checkRefund({ tradeid: item.tradeID }).then((res: any) => {
		if (res.code === 'SUCCESS' && res.data.refund) {
			ElMessage.success('退款成功');
		} else {
			ElMessage.warning('您暂未退款成功');
		}
	});
};
const itemRefund = (item: any) => {
	phoneDialogRef.value.openDialog(item);
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
