<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="id" label="用户ID" show-overflow-tooltip min-width="240" align="left"></el-table-column>
				<el-table-column prop="phone" label="手机号" show-overflow-tooltip align="center" width="140"></el-table-column>
				<el-table-column prop="inviteCount" label="邀请人数" show-overflow-tooltip align="center" width="100"></el-table-column>
				<el-table-column prop="status" label="状态" show-overflow-tooltip align="center" width="100">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.status == 'NORMAL'">正常</el-tag>
						<el-tag type="info" v-else>禁用</el-tag>
					</template>
				</el-table-column>

				<el-table-column prop="createTime" label="创建时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
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
import { reactive, onMounted } from 'vue';
import dayjs from 'dayjs';
import { useInviteApi } from '/@/api/invite';

const state = reactive({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
		},
	},
});

const getTableData = () => {
	state.tableData.loading = true;
	useInviteApi()
		.inviterList()
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.user;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};

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
