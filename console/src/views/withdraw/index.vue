<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
						<el-input v-model="state.tableData.param.src" placeholder="请输入关键词" style="width: 100%"> </el-input>
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
				<el-table-column prop="withdrawID" label="提现ID" show-overflow-tooltip align="left"></el-table-column>
				<el-table-column prop="name" label="名称" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="withdrawWay" label="提现方式" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="status" label="状态" show-overflow-tooltip align="center" width="120">
					<template #default="scope">
						<el-tag v-if="scope.row.status === 1" type="warning">等待提现</el-tag>
						<el-tag v-else-if="scope.row.status === 2" type="success">提现成功</el-tag>
						<el-tag v-else type="danger">提现失败</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="cny" label="金额" show-overflow-tooltip width="200" align="center">
					<template #default="scope">
						<el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.cny) }}元</el-text>
					</template>
				</el-table-column>
				<el-table-column label="提现时间" show-overflow-tooltip align="center">
					<template #default="scope">
						<span v-if="scope.row.withdrawAt">{{ dayjs.unix(scope.row.withdrawAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
						<span v-else>-</span>
					</template>
				</el-table-column>
				<el-table-column label="支付时间" show-overflow-tooltip align="center">
					<template #default="scope">
						<span v-if="scope.row.payAt">{{ dayjs.unix(scope.row.payAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
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
import { withdrawApi } from '/@/api/withdraw';
import { formatAmount } from '/@/utils/formatAmount';
import dayjs from 'dayjs';
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
	withdrawApi()
		.getwithdrawList(state.tableData.param)
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				state.tableData.data = res.data.withdraw;
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
