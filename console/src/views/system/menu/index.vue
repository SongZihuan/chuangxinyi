<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="search" class="demo-form-inline" @submit.prevent>
					<el-form-item label="菜单名称">
						<el-input placeholder="请输入菜单名称" style="max-width: 180px" v-model="search.name" clearable> </el-input>
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
							新增菜单
						</el-button>
						<!--            从数据库更新菜单-->
						<el-button type="success" class="ml10" @click="updateMenuDB">
							<el-icon>
								<RefreshRight />
							</el-icon>
							从数据库更新菜单
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
				<el-table-column prop="id" label="菜单ID" show-overflow-tooltip width="100" align="left"></el-table-column>
				<el-table-column label="菜单名称" show-overflow-tooltip width="240" align="center">
					<template #default="scope">
						<SvgIcon :name="scope.row.icon" />
						<span class="ml10">{{ scope.row.title }}</span>
					</template>
				</el-table-column>
				<el-table-column label="菜单name" show-overflow-tooltip width="200" prop="name"> </el-table-column>
				<el-table-column prop="path" label="路由路径" show-overflow-tooltip width="180"></el-table-column>
				<el-table-column label="组件路径" show-overflow-tooltip min-width="180">
					<template #default="scope">
						<span>{{ scope.row.componentAlias }}</span>
					</template>
				</el-table-column>
				<el-table-column label="类型" show-overflow-tooltip width="80">
					<template #default="scope">
						<el-tag type="success" size="small" v-if="scope.row.menuType == 1">菜单</el-tag>
						<el-tag size="small" v-if="scope.row.menuType == 2">按钮</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="操作" show-overflow-tooltip width="230" align="center" fixed="right">
					<template #default="scope">
						<el-button text type="primary" @click="onOpenBindRole(scope.row)">角色</el-button>
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
		<!-- 菜单弹窗 -->
		<MenuDialog ref="menuDialogRef" @refresh="getTableData()" :menuList="state.tableData.data" />
		<!-- 绑定角色菜单 -->
		<rolrDialog ref="roleRef" @refresh="getTableData()" />
	</div>
</template>

<script setup lang="ts" name="systemMenu">
import { defineAsyncComponent, ref, onMounted, reactive } from 'vue';
import { RouteRecordRaw } from 'vue-router';
import { ElMessageBox, ElMessage } from 'element-plus';
import { useMenuApi } from '/@/api/system/menu/index';
import type { searchTypes } from '/@/api/system/menu/types';
import commonFunction from '/@/utils/commonFunction';
import { RefreshRight } from '@element-plus/icons-vue';
// 引入组件
const MenuDialog = defineAsyncComponent(() => import('./dialog.vue'));
const rolrDialog = defineAsyncComponent(() => import('./bindRole.vue'));
const menuApi = useMenuApi();
const { flatToTree } = commonFunction();
// 定义变量内容
const menuDialogRef = ref();
const roleRef = ref();
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
	menuApi.menuList(search).then((res: any) => {
		if (res.code === "SUCCESS" && res.data) {
			state.tableData.data = flatToTree(res.data.menu, 'id', 'parentID');
			search.total = res.data.count;
		}
		state.tableData.loading = false;
	});
};
//查询
const handleSearch = () => {
	getTableData();
};
// 打开新增菜单弹窗
const onOpenAddMenu = (type: string) => {
	menuDialogRef.value.openDialog(type);
};
// 打开编辑菜单弹窗
const onOpenEditMenu = (type: string, row: RouteRecordRaw) => {
	menuDialogRef.value.openDialog(type, row);
};
//移动菜单
const onMoveMenu = (row: { id: number }, type: string) => {
	menuApi.menuMove({ id: row.id, isUp: type === 'top' }).then((res: any) => {
		if (res.code === "SUCCESS") {
			type === 'top' ? ElMessage.success('上移成功') : ElMessage.success('下移成功');
			getTableData();
		}
	});
};
const onOpenBindRole = (row: any) => {
	roleRef.value.openDialog(row);
};

// 从数据库更新菜单
const updateMenuDB = () => {
	ElMessageBox.confirm(`此操作将从数据库更新菜单，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			menuApi.updateMenuFromDB().then((res: any) => {
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
			menuApi.menuDelete({ id: row.id }).then((res: any) => {
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
