<template>
	<div class="layout-container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="search" class="demo-form-inline" @submit.prevent>
					<el-form-item>
						<el-button type="success" @click="onOpenAddPath('add')">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增路由
						</el-button>
						<el-button type="success" class="ml10" @click="updatePathDB">
							<el-icon>
								<RefreshRight />
							</el-icon>
							从数据库更新路由
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="id" label="路由ID" show-overflow-tooltip width="100" align="left"></el-table-column>
				<el-table-column prop="path" label="路由路径" show-overflow-tooltip width="180"></el-table-column>
				<el-table-column prop="path" label="路由方法" show-overflow-tooltip width="180">
					<template #default="scope">
						<span>{{ scope.row.method && scope.row.method.length > 0 ? scope.row.method.toString() : '--' }}</span>
					</template>
				</el-table-column>

				<el-table-column prop="isOr" label="路由匹配" show-overflow-tooltip width="180" align="center">
					<template #default="scope">
						<span>{{ pathModeDict.filter((item) => item.value == scope.row.mode)[0].label }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="status" label="路由状态" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.status === 1">启用</el-tag>
						<el-tag type="info" v-else-if="2">禁用</el-tag>
						<el-tag type="info" v-else>不放行</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="describe" label="路由描述" show-overflow-tooltip min-width="180" align="center"></el-table-column>

				<el-table-column prop="path" label="policy关系" show-overflow-tooltip width="110" align="center">
					<template #default="scope">
						<span>{{ scope.row.isOr ? 'or' : 'and' }}</span>
					</template>
				</el-table-column>

				<el-table-column label="操作" show-overflow-tooltip width="170" align="center" fixed="right">
					<template #default="scope">
						<el-button text type="primary" @click="onOpenBindRole(scope.row)">查看角色</el-button>
						<el-button text type="primary" @click="onOpenEditPath('edit', scope.row)">修改</el-button>
						<el-button text type="primary" @click="onTabelRowDel(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
			<el-pagination hide-on-single-page
				@size-change="onHandleSizeChange"
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				v-model:current-page="search.page"
				background
				v-model:page-size="search.pagesize"
				:total="search.total"
			>
			</el-pagination>
		</el-card>
		<!-- 路由弹窗 -->
		<PathDialog ref="pathDialogRef" @refresh="getTableData()" :pathList="state.tableData.data" />
		<PolicyDialog ref="policyRef" />
	</div>
</template>

<script setup lang="ts" name="systemBackPath">
import { defineAsyncComponent, ref, onMounted, reactive } from 'vue';
import { RouteRecordRaw } from 'vue-router';
import { pathModeDict } from '/@/api/system/path/types';
import { ElMessageBox, ElMessage } from 'element-plus';
import { usePathApi } from '/@/api/system/backPath/index';
import { RefreshRight } from '@element-plus/icons-vue';
import PolicyDialog from '/@/views/system/permission/policyDialog.vue';
// 引入组件
const PathDialog = defineAsyncComponent(() => import('./dialog.vue'));
const pathApi = usePathApi();

// 定义变量内容
const pathDialogRef = ref();

const policyRef = ref();
const state = reactive({
	tableData: {
		data: [] as RouteRecordRaw[],
		loading: true,
	},
});
const search = reactive({
	page: 1,
	pagesize: 20,
	total: 0,
});

const getTableData = () => {
	state.tableData.loading = true;
	pathApi.pathList(search).then((res: any) => {
		if (res.code === "SUCCESS" && res.data) {
			state.tableData.data = res.data.path;
			search.total = res.data.count;
		}
		state.tableData.loading = false;
	});
};

// 打开新增路由弹窗
const onOpenAddPath = (type: string) => {
	pathDialogRef.value.openDialog(type);
};
// 打开编辑路由弹窗
const onOpenEditPath = (type: string, row: RouteRecordRaw) => {
	pathDialogRef.value.openDialog(type, row);
};

const onOpenBindRole = (row: any) => {
	policyRef.value.openDialog(row);
};

// 从数据库更新路由
const updatePathDB = () => {
	ElMessageBox.confirm(`此操作将从数据库更新路由，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			pathApi.updatepathFromDB().then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('更新成功');
					getTableData();
				}
			});
		})
		.catch(() => {});
};
// 删除当前行
const onTabelRowDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除路由：${row.path}, 是否继续?`, '提示', {
		confirmButtonText: '删除',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			pathApi.pathDelete({ id: row.id }).then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('删除成功');
					getTableData();
				}
			});
		})
		.catch(() => {});
};
const onHandleSizeChange = (val: number) => {
	search.pagesize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	search.page = val;
	getTableData();
};
onMounted(() => {
	getTableData();
});
</script>
