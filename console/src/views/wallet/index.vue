<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
						<el-input v-model="state.tableData.param.src" placeholder="请输入关键词" style="width: 100%"> </el-input>
					</el-form-item>
					<el-form-item>
						<el-input v-model="state.tableData.param.fundingID" placeholder="请输入资金ID" style="width: 100%"> </el-input>
					</el-form-item>
					<el-form-item>
						<el-select v-model="state.tableData.param.timetype" placeholder="请选择日期类型" style="width: 100%">
							<el-option v-for="(item, index) in dateTypeDict" :value="item.value" :label="item.label" :key="index" />
						</el-select>
					</el-form-item>
					<el-form-item>
						<el-select v-model="state.tableData.param.status" placeholder="请选择状态" style="width: 100%">
							<el-option v-for="(item, index) in statusDict" :value="item.value" :label="item.label" :key="index" />
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
						<el-button type="primary" class="ml10" @click="handleSearch">
							<el-icon>
								<ele-Search />
							</el-icon>
							查询
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="id" label="ID" show-overflow-tooltip width="100" align="left"></el-table-column>
				<el-table-column prop="walletID" label="钱包ID" show-overflow-tooltip width="100" align="left"></el-table-column>
				<el-table-column prop="userID" label="用户ID" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="user" label="用户名称" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="type" label="类型" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="funding_id" label="资金ID" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="reason" label="原因" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="balance" label="余额" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="cny" label="CNY" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="notBilled" label="未开票" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="billed" label="开票总额" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-text class="mx-1" type="success" tag="b">{{ formatAmount(scope.row.billed) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column prop="hasBilled" label="已开票" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.hasBilled === 1" type="success">是</el-tag>
						<el-tag v-else type="danger">否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="withdraw" label="提现金额" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-text class="mx-1" type="success" tag="b">{{ formatAmount(scope.row.withdraw) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column prop="notWithdraw" label="未提现" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.notWithdraw === 1" type="success">是</el-tag>
						<el-tag v-else type="danger">否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="hasWithdraw" label="已提现" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.hasWithdraw === 1" type="success">是</el-tag>
						<el-tag v-else type="danger">否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="BeforeBalance" label="变更前余额" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.BeforeBalance) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column prop="BeforeCny" label="变更前CNY" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.BeforeCny) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column prop="BeforeNotBilled" label="变更前未开票" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.BeforeNotBilled === 1" type="success">是</el-tag>
						<el-tag v-else type="danger">否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="BeforeBilled" label="变更前开票总额" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-text class="mx-1" type="success" tag="b">{{ formatAmount(scope.row.BeforeBilled) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column prop="BeforeHasBilled" label="变更前已开票" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.BeforeHasBilled === 1" type="success">是</el-tag>
						<el-tag v-else type="danger">否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="BeforeWithdraw" label="变更前提现金额" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-text class="mx-1" type="success" tag="b">{{ formatAmount(scope.row.BeforeWithdraw) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column prop="BeforeNotWithdraw" label="变更前未提现" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.BeforeNotWithdraw === 1" type="success">是</el-tag>
						<el-tag v-else type="danger">否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="BeforeHasWithdraw" label="变更前已提现" show-overflow-tooltip min-width="200" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.BeforeHasWithdraw === 1" type="success">是</el-tag>
						<el-tag v-else type="danger">否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="remark" label="备注" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column label="创建时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
						<span v-else>-</span>
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
	</div>
</template>

<script setup lang="ts">
import { onMounted, reactive, watch } from 'vue';
import { formatAmount } from '/@/utils/formatAmount';
import dayjs from 'dayjs';
import { useWalletRecordApi } from '/@/api/wallet/record';
const state = reactive<any>({
	tableData: {
		data: [],
		loading: false,
		param: {
			range: [],
			page: 1,
			pagesize: 20,
			timetype: '',
			starttime: null,
			endtime: null,
			status: '',
			src: '',
			fundingID: '',
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
const statusDict = [
	{ label: '等待提现', value: 1 },
	{ label: '提现成功', value: 2 },
	{ label: '提现失败', value: 3 },
];
const dateTypeDict = [
	{ label: '创建时间', value: 1 },
	{ label: '支付时间', value: 2 },
	{ label: '提现时间', value: 23 },
];
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
const getTableData = () => {
	state.tableData.loading = true;
	useWalletRecordApi()
		.userWalletList(state.tableData.param)
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				state.tableData.data = res.data.record;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
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
onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss"></style>
