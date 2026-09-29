<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" destroy-on-close @close="closeDialog">
		<div class="mb20">
			已绑定角色: <el-tag class="ml10" type="success">{{ info.roleName }}</el-tag>
		</div>
		<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
			<el-table-column prop="id" label="角色ID" show-overflow-tooltip width="80" align="left"></el-table-column>
			<el-table-column prop="name" label="角色名称" show-overflow-tooltip width="100" align="center"></el-table-column>
			<el-table-column prop="sign" label="角色标识" show-overflow-tooltip width="100" align="center"></el-table-column>
			<el-table-column prop="describe" label="角色描述" show-overflow-tooltip align="center"> </el-table-column>
			<el-table-column label="操作" width="60" align="center">
				<template #default="scope">
					<el-button size="small" text type="primary" @click="bingRole(scope.row)" :disabled="info.roleID == scope.row.id">绑定</el-button>
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
			:total="state.tableData.param.total"
		>
		</el-pagination
	></el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { roleApi } from '/@/api/system/role';
import type { SysRoleState } from '/@/api/system/role/types';
import { userApi } from '/@/api/system/user';
import { message } from '/@/utils/message';
const userApiCollect = userApi();
const info = ref<any>({});
const roleApiCollect = roleApi();
const emit = defineEmits(['refresh']);
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '绑定角色',
	submitTxt: '',
});

const state = reactive<SysRoleState>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			name: '',
			page: 1,
			pagesize: 10,
			total: 0,
		},
	},
});
const getTableData = () => {
	state.tableData.loading = true;
	roleApiCollect.getRoleList(state.tableData.param).then((res: any) => {
		if (res.code === "SUCCESS") {
			state.tableData.data = res.data.role;
			state.tableData.param.total = res.data.count;
			state.tableData.loading = false;
		}
	});
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const bingRole = (row: any) => {
	userApiCollect.bindRole({ id: info.value.ID, uid: info.value.id, roleID: row.id }).then((res: any) => {
		if (res.code === "SUCCESS") {
			emit('refresh');
			message('绑定成功', { type: 'success' });
			closeDialog();
		}
	});
};

const openDialog = async (row?: any) => {
	await getTableData();
	info.value = JSON.parse(JSON.stringify(row));
	dialog.isShowDialog = true;
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

// 暴露变量
defineExpose({
	openDialog,
});
</script>
