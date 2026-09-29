<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="search" class="demo-form-inline" @submit.prevent>
					<el-form-item label="应用名称">
						<el-input placeholder="请输入应用名称" style="max-width: 180px" v-model="search.name" clearable> </el-input>
					</el-form-item>
					<el-form-item>
						<el-button type="primary" class="ml10" @click="handleSearch">
							<el-icon>
								<ele-Search />
							</el-icon>
							查询
						</el-button>
						<el-button type="success" class="ml10" @click="onOpenAddMenu('add')">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增应用
						</el-button>
						<!--            从数据库更新应用-->
						<el-button type="success" class="ml10" @click="updateMenuDB">
							<el-icon>
								<RefreshRight />
							</el-icon>
							从数据库更新应用
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
				<el-table-column prop="id" label="应用ID" show-overflow-tooltip width="100" align="left"></el-table-column>
				<el-table-column label="应用图标" show-overflow-tooltip width="85" align="center">
					<template #default="scope">
						<SvgIcon :name="scope.row.icon" />
						<span class="ml10">{{ scope.row.title }}</span>
					</template>
				</el-table-column>
				<el-table-column label="应用名称" show-overflow-tooltip width="200" prop="name" align="center"> </el-table-column>
				<el-table-column prop="webName" label="网站" show-overflow-tooltip width="180" align="center"></el-table-column>
				<el-table-column prop="describe" label="描述" show-overflow-tooltip min-width="180" align="center"></el-table-column>
				<el-table-column label="操作" show-overflow-tooltip width="220" align="center" fixed="right">
					<template #default="scope">
						<el-button text type="primary" @click="onMoveMenu(scope.row, 'top')">上移</el-button>
						<el-button text type="primary" @click="onMoveMenu(scope.row, 'down')">下移</el-button>
						<el-button text type="primary" @click="onOpenEditMenu('edit', scope.row)">修改</el-button>
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
		<!-- 应用弹窗 -->
		<MenuDialog ref="menuDialogRef" @refresh="getTableData()" :menuList="state.tableData.data" />
	</div>
</template>

<script setup lang="ts" name="systemMenu">
import { defineAsyncComponent, ref, onMounted, reactive } from 'vue';
import { RouteRecordRaw } from 'vue-router';
import { ElMessageBox, ElMessage } from 'element-plus';
import { useApplicationApi } from '/@/api/system/application/index';
import type { searchTypes } from '/@/api/system/application/types';

import { RefreshRight } from '@element-plus/icons-vue';
// 引入组件
const MenuDialog = defineAsyncComponent(() => import('./dialog.vue'));

const applicationApi = useApplicationApi();

// 定义变量内容
const menuDialogRef = ref();
const state = reactive({
	tableData: {
		data: [] as RouteRecordRaw[],
		loading: true,
	},
});
const search = reactive<searchTypes>({
	name: '',
	page: 1,
	pagesize: 1000,
	total: 0,
});

const getTableData = () => {
	state.tableData.loading = true;
	applicationApi.applicationList(search).then((res: any) => {
		if (res.code === "SUCCESS" && res.data) {
			(state.tableData.data = res.data.application), (search.total = res.data.count);
		}
		state.tableData.loading = false;
	});
};
//查询
const handleSearch = () => {
	getTableData();
};
// 打开新增应用弹窗
const onOpenAddMenu = (type: string) => {
	menuDialogRef.value.openDialog(type);
};
// 打开编辑应用弹窗
const onOpenEditMenu = (type: string, row: RouteRecordRaw) => {
	menuDialogRef.value.openDialog(type, row);
};

//移动应用
const onMoveMenu = (row: { id: number }, type: string) => {
	applicationApi.applicationMove({ id: row.id, isUp: type === 'top' }).then((res: any) => {
		if (res.code === "SUCCESS") {
			type === 'top' ? ElMessage.success('上移成功') : ElMessage.success('下移成功');
			getTableData();
		}
	});
};

// 从数据库更新应用
const updateMenuDB = () => {
	ElMessageBox.confirm(`此操作将从数据库更新应用，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			applicationApi.updateapplicationFromDB().then((res: any) => {
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
	ElMessageBox.confirm(`此操作将永久删除应用名称：${row.name}, 是否继续?`, '提示', {
		confirmButtonText: '删除',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			applicationApi.applicationDelete({ id: row.id }).then((res: any) => {
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
