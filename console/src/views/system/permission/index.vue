<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="search" class="demo-form-inline" @submit.prevent>
					<el-form-item>
						<el-button type="success" class="ml10" @click="onOpenAddPermission('add')">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增权限
						</el-button>
						<!--            从数据库更新权限-->
						<el-button type="success" class="ml10" @click="updatePermissionDB">
							<el-icon>
								<RefreshRight />
							</el-icon>
							从数据库更新权限
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table
				:data="state.tableData.data"
				v-loading="state.tableData.loading"
				style="width: 100%"
				row-key="id"
				:tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
				border
			>
				<el-table-column prop="id" label="权限ID" show-overflow-tooltip width="80" align="left"></el-table-column>
				<el-table-column prop="name" label="权限名称" show-overflow-tooltip width="180"></el-table-column>
				<el-table-column label="权限标识" prop="sign" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="describe" label="权限描述" show-overflow-tooltip min-width="180"></el-table-column>
				<el-table-column prop="status" label="是否陌生用户权限" show-overflow-tooltip width="140" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.isAnonymous">是</el-tag>
						<el-tag type="info" v-else>否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="status" label="是否普通用户权限" show-overflow-tooltip width="140" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.isUser">是</el-tag>
						<el-tag type="info" v-else>否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="status" label="状态" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.status === 1">启用</el-tag>
						<el-tag type="info" v-else>禁用</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="操作" show-overflow-tooltip width="170" align="center" fixed="right">
					<template #default="scope">
						<el-button text type="primary" @click="onOpenBindRole(scope.row)">查看角色</el-button>
						<el-button text type="primary" @click="onOpenEditPermission('edit', scope.row)">修改</el-button>
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
		<!-- 权限弹窗 -->
		<PermissionDialog ref="permissionDialogRef" @refresh="getTableData()" :permissionList="state.tableData.data" />
		<PolicyDialog ref="policyRef" />
	</div>
</template>

<script setup lang="ts" name="systemPermission">
import { defineAsyncComponent, ref, onMounted, reactive } from 'vue';
import { RouteRecordRaw } from 'vue-router';
import { ElMessageBox, ElMessage } from 'element-plus';
import { usePermissionApi } from '/@/api/system/permission/index';
import type { searchTypes } from '/@/api/system/permission/types';
import PolicyDialog from '/@/views/system/permission/policyDialog.vue';
import { RefreshRight } from '@element-plus/icons-vue';
// 引入组件
const PermissionDialog = defineAsyncComponent(() => import('./dialog.vue'));
const permissionApi = usePermissionApi();
// 定义变量内容
const permissionDialogRef = ref();
const policyRef = ref();
const state = reactive({
	tableData: {
		data: [] as RouteRecordRaw[],
		loading: true,
	},
});
const search = reactive<searchTypes>({
	page: 1,
	pagesize: 20,
	total: 0,
});
const onOpenBindRole = (row: any) => {
	policyRef.value.openDialog(row);
};
const getTableData = () => {
	state.tableData.loading = true;
	permissionApi.permissionList(search).then((res: any) => {
		if (res.code === "SUCCESS" && res.data) {
			state.tableData.data = res.data.permission;
			search.total = res.data.count;
		}
		state.tableData.loading = false;
	});
};

// 打开新增权限弹窗
const onOpenAddPermission = (type: string) => {
	permissionDialogRef.value.openDialog(type);
};
// 打开编辑权限弹窗
const onOpenEditPermission = (type: string, row: RouteRecordRaw) => {
	permissionDialogRef.value.openDialog(type, row);
};

// 从数据库更新权限
const updatePermissionDB = () => {
	ElMessageBox.confirm(`此操作将从数据库更新权限，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			permissionApi.updatePermissionFromDB().then((res: any) => {
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
	ElMessageBox.confirm(`此操作将永久删除权限：${row.name}, 是否继续?`, '提示', {
		confirmButtonText: '删除',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			permissionApi.permissionDelete({ id: row.id }).then((res: any) => {
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
