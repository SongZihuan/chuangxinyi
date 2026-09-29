<template>
	<div class="container layout-padding">
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
				<el-table-column prop="webID" label="网站ID" show-overflow-tooltip width="150" align="center"></el-table-column>
				<el-table-column prop="webName" label="网站名称" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="webDomain" label="网站域名" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="ip" label="ip地址" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="geo" label="地理位置" show-overflow-tooltip width="300" align="center"></el-table-column>
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
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, watch } from 'vue';
import dayjs from 'dayjs';
import { authorizationRecordStateType } from './types';
import { useoauth2Api } from '/@/api/oauth2';

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
	useoauth2Api()
		.oauth2Record(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				// 排序 webID

				state.tableData.data = res.data.record.sort((a, b) => {
					return b.webID - a.webID;
				});
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
