<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-input v-model="state.tableData.param.name" placeholder="请输入角色名称" style="max-width: 180px" clearable> </el-input>
				<el-button size="default" type="primary" class="ml10" @click="handleSearch">
					<el-icon>
						<ele-Search />
					</el-icon>
					查询
				</el-button>
				<el-button size="default" type="success" class="ml10" @click="onOpenAddRole('add')">
					<el-icon>
						<ele-FolderAdd />
					</el-icon>
					新增角色
				</el-button>
				<!--        从数据库更新角色-->
				<el-button size="default" type="success" class="ml10" @click="updateRoleDB">
					<el-icon>
						<RefreshRight />
					</el-icon>
					从数据库更新角色
				</el-button>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="id" label="角色ID" show-overflow-tooltip width="100" align="left"></el-table-column>
				<el-table-column prop="name" label="角色名称" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column prop="sign" label="角色标识" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="belong" label="所属站点" width="180" show-overflow-tooltip align="center">
					<template #default="scope">
						<span>{{ scope.row.belong === 0 ? '用户中心全局' : scope.row.belongName }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="status" label="角色状态" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.status === 1">启用</el-tag>
						<el-tag type="info" v-else>禁用</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="describe" label="角色描述" show-overflow-tooltip width="170" align="center"> </el-table-column>
				<el-table-column prop="createTime" label="创建时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="180" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="roleMoveRole(scope.row, 'top')">上移</el-button>
						<el-button text type="primary" @click="roleMoveRole(scope.row, 'down')">下移</el-button>
						<el-button :disabled="scope.row.roleName === '超级管理员'" text type="primary" @click="onOpenEditRole('edit', scope.row)">修改</el-button>
						<el-button :disabled="scope.row.roleName === '超级管理员'" text type="primary" @click="onRowDel(scope.row)">删除</el-button>
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
		<RoleDialog ref="roleDialogRef" @refresh="getTableData()" />
	</div>
</template>

<script setup lang="ts" name="systemRole">
import { defineAsyncComponent, reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import { roleApi } from '/@/api/system/role/index';
import type { SysRoleState } from '/@/api/system/role/types';
import dayjs from 'dayjs';
import { RefreshRight } from '@element-plus/icons-vue';
// 引入组件
const RoleDialog = defineAsyncComponent(() => import('/@/views/system/role/dialog.vue'));
const roleApiCollect = roleApi();

// 定义变量内容
const roleDialogRef = ref();
const state = reactive<SysRoleState>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			name: '',
			page: 1,
			pagesize: 20,
		},
	},
});
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	roleApiCollect.getRoleList(state.tableData.param).then((res: any) => {
		if (res.code === "SUCCESS") {
			state.tableData.data = res.data.role;
			state.tableData.total = res.data.count;
			state.tableData.loading = false;
			state.tableData.loading = false;
		}
	});
};
//搜索
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
//移动菜单
const roleMoveRole = (row: { id: number }, type: string) => {
	roleApiCollect.roleMove({ id: row.id, isUp: type === 'top' }).then((res: any) => {
		if (res.code === "SUCCESS") {
			type === 'top' ? ElMessage.success('上移成功') : ElMessage.success('下移成功');
			getTableData();
		}
	});
};
// 打开新增角色弹窗
const onOpenAddRole = (type: string) => {
	roleDialogRef.value.openDialog(type);
};
// 打开修改角色弹窗
const onOpenEditRole = (type: string, row: Object) => {
	roleDialogRef.value.openDialog(type, row);
};
// 从数据库更新角色
const updateRoleDB = () => {
	ElMessageBox.confirm(`此操作将从数据库更新角色, 是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			roleApiCollect.updateRoleFromDB().then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('更新成功');
				} else {
					ElMessage.error('更新失败');
				}
				getTableData();
			});
		})
		.catch(() => {});
};
// 删除角色
const onRowDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除角色名称：“${row.name}”，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			roleApiCollect.roleDelAction({ id: row.id }).then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('删除成功');
				} else {
					ElMessage.error('删除失败');
				}
				getTableData();
			});
		})
		.catch(() => {});
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
