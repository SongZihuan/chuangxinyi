<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
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
				<el-table-column prop="level" label="分销层级" show-overflow-tooltip min-width="200" align="left"></el-table-column>
				<el-table-column prop="pre" label="返佣比例" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column label="操作" width="120" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="openDialog('edit', scope.row)">编辑</el-button>
						<el-button text type="primary" @click="onTabelRowDel(scope.row)">删除</el-button>
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
		<ActionDialog ref="actionDialogRef" @refresh="getTableData" />
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, defineAsyncComponent } from 'vue';
import { useDistributionApi } from '/@/api/distribution';
import { ElMessage, ElMessageBox } from 'element-plus';
const ActionDialog = defineAsyncComponent(() => import('./component/actionDialog.vue'));
const actionDialogRef = ref();
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
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	useDistributionApi()
		.distributionList()
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.coupons;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};

const openDialog = (type: string, row?: any) => {
	actionDialogRef.value.openDialog(type, row);
};
// 分页改变
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
	getTableData();
};
const onTabelRowDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除：${row.level}, 是否继续?`, '提示', {
		confirmButtonText: '删除',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useDistributionApi()
				.distributionDel({ level: row.level })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('删除成功');
						getTableData();
					}
				});
		})
		.catch(() => {});
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
