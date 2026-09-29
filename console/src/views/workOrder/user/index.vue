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
						<el-button type="success" @click="openDialog('add')">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="orderID" label="ID" width="220" show-overflow-tooltip align="left"></el-table-column>
				<el-table-column prop="title" label="标题" width="180" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="from" label="来源" min-width="180" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="status" label="状态" width="180" show-overflow-tooltip align="center">
					<!--          1等待用户回复，2等待网站回复，3已完成-->
					<template #default="scope">
						<el-tag v-if="scope.row.status == 1">等待用户回复</el-tag>
						<el-tag v-if="scope.row.status == 2" type="warning">等待管理员回复</el-tag>
						<el-tag v-if="scope.row.status == 3" type="success">已完成</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="lastReplyAt" label="最后回复时间" width="180" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.lastReplyAt">{{ dayjs.unix(scope.row.lastReplyAt).format('YYYY-MM-DD HH:mm:ss') }}</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="finishAt" label="完成时间" width="180" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.finishAt">{{ dayjs.unix(scope.row.finishAt).format('YYYY-MM-DD HH:mm:ss') }}</el-tag>
						<el-tag v-else type="info"> 未完成</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="createAt" label="创建时间" width="180" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="80" fixed="right" align="center">
					<template #default="scope">
						<el-dropdown>
							<span class="el-dropdown-link">
								操作
								<el-icon class="el-icon--right">
									<arrow-down />
								</el-icon>
							</span>
							<template #dropdown>
								<el-dropdown-menu>
									<el-dropdown-item @click="onCommunicate(scope.row)">沟通记录</el-dropdown-item>
									<el-dropdown-item v-if="!scope.row.finishAt" @click="onFinish(scope.row)"> 完成工单</el-dropdown-item>
								</el-dropdown-menu>
							</template>
						</el-dropdown>
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
		<ActionDialog ref="actionDialogRef" @refresh="getTableData"></ActionDialog>
		<CommunicateDialog ref="communicateDialogRef"></CommunicateDialog>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, watch, ref } from 'vue';

import dayjs from 'dayjs';
import { ArrowDown } from '@element-plus/icons-vue';
import { useUserWorkOrderApi } from '/@/api/workOrder/user';
import ActionDialog from './component/actionDialog.vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import CommunicateDialog from './component/communicateDialog.vue';
import { IWorkOrderUserListResData } from '/@/views/workOrder/user/types';
import { useRouter } from 'vue-router';
const dateTypeDict = ref([
	{ label: '创建时间', value: 1 },
	{ label: '完成时间', value: 7 },
	{ label: '上次回复时间', value: 8 },
]);
const state = reactive<IWorkOrderUserListResData>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			starttime: 0,
			endtime: 0,
			range: [],
			timetype: '',
      fromID: -1,
		},
	},
});
const router = useRouter();
const actionDialogRef = ref();
const communicateDialogRef = ref();
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
const openDialog = (row: any) => {
	actionDialogRef.value.openDialog(row);
};
const onCommunicate = (row: any) => {
	router.push({ path: '/workOrder/user/communicateIndex', query: { id: row.orderID, finishStatus: row.finishAt ? true : false } });
	// communicateDialogRef.value.openDialog(row);
};

const onFinish = (row: any) => {
	ElMessageBox.confirm(`此操作将完成工单：${row.title}, 是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useUserWorkOrderApi()
				.userOrderFinsh({ orderID: row.orderID })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('完成工单成功');
						getTableData();
					}
				});
		})
		.catch(() => {});
};
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	useUserWorkOrderApi()
		.userWorkOrderList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.order;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
const handleSearch = () => {
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
