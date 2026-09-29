<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" destroy-on-close @close="closeDialog">
		<el-table :data="info.roles" v-loading="state.tableData.loading" style="width: 100%" border>
			<el-table-column prop="id" label="角色名称ID" show-overflow-tooltip width="120" align="center"></el-table-column>
			<el-table-column prop="name" label="角色名称" show-overflow-tooltip width="100" align="center"></el-table-column>
			<el-table-column prop="sign" label="角色标识" show-overflow-tooltip width="100" align="center"></el-table-column>
			<el-table-column prop="describe" label="角色描述" show-overflow-tooltip align="center"> </el-table-column>
		</el-table>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { roleApi } from '/@/api/system/role/index';
import type { SysRoleState } from '/@/api/system/role/types';

const roles = ref();
const info = ref<any>({});
const roleApiCollect = roleApi();
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '角色',
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
			pageSize: 10,
			total: 0,
		},
	},
});
const getTableData = () => {
	state.tableData.loading = true;
	roleApiCollect.getRoleList(state.tableData.param).then((res: any) => {
		if (res.code === "SUCCESS") {
			let data = res.data.role;
			data = data.filter((item: any) => {
				return roles.value.includes(item.roleID);
			});
			state.tableData.data = data;
			state.tableData.param.total = res.data.count;
			state.tableData.loading = false;
		}
	});
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const openDialog = async (row?: any) => {
	await getTableData();
	info.value = JSON.parse(JSON.stringify(row));
	roles.value = info.value.roles.map((item: any) => {
		return item.id;
	});
	dialog.isShowDialog = true;
};
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pageSize = val;
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
