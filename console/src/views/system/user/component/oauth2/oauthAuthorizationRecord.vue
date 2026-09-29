<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="869px" destroy-on-close @close="closeDialog">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
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
				<el-table-column prop="webID" label="网站ID" show-overflow-tooltip min-width="100" align="left"></el-table-column>
				<el-table-column prop="webName" label="网站名称" show-overflow-tooltip width="150" align="left"></el-table-column>
				<el-table-column prop="ip" label="ip地址" show-overflow-tooltip width="150" align="left"></el-table-column>
				<el-table-column prop="geo" label="地理位置" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="loginTime" label="登录时间" show-overflow-tooltip width="200" align="center">
					<template #default="scope">
						<span v-if="scope.row.loginTime">{{ dayjs.unix(scope.row.loginTime).format('YYYY-MM-DD HH:mm:ss') }}</span>
						<span v-else>-</span>
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
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue';
import dayjs from 'dayjs';
import { userApi } from '/@/api/system/user';
import { authorizationRecordStateType } from '/@/views/oauth2/authorizationRecord/types';

const state = reactive<authorizationRecordStateType>({
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
      uid: '',
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
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '用户授权记录列表',
	submitTxt: '',
});
const closeDialog = () => {
	dialog.isShowDialog = false;
};

const openDialog = async (row?: any) => {
	state.tableData.param.uid = row.id;
	dialog.isShowDialog = true;
	getTableData();
};
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	userApi()
		.getUserOauth2RecordList({ uid: state.tableData.param.uid, ...state.tableData.param })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.record;
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
// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss"></style>
