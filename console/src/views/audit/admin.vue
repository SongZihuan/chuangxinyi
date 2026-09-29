<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
						<el-input v-model="state.tableData.param.src" placeholder="搜索条件" style="width: 100%" />
					</el-form-item>
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
				<el-table-column prop="userID" label="用户ID" show-overflow-tooltip min-width="100" align="left"></el-table-column>
				<el-table-column prop="content" label="审计内容" show-overflow-tooltip min-width="200" align="left"></el-table-column>
				<el-table-column prop="from" label="来源" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="createAt" label="审计时间" show-overflow-tooltip width="200" align="center">
					<template #default="scope">
						<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
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
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { auditStateType } from './types';
import { useAuditApi } from '/@/api/audit';
const dateTypeDict = ref([{ label: '创建时间', value: 1 }]);
const state = reactive<auditStateType>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			timetype: '',
			starttime: 0,
			endtime: 0,
			range: [],
      src: '',
      fromID: -1,
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
	useAuditApi()
		.adminAuditList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.audit;
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
