<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" destroy-on-close @close="closeDialog">
		<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
			<el-table-column prop="id" label="角色名称ID" show-overflow-tooltip width="120" align="center"></el-table-column>
			<el-table-column prop="name" label="角色名称" show-overflow-tooltip width="100" align="center"></el-table-column>
			<el-table-column prop="sign" label="角色标识" show-overflow-tooltip width="100" align="center"></el-table-column>
			<el-table-column prop="describe" label="角色描述" show-overflow-tooltip align="center"> </el-table-column>
		</el-table>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import type { SysRoleState } from '/@/api/system/role/types';
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
			pagesize: 10,
			total: 0,
		},
	},
});

const closeDialog = () => {
	dialog.isShowDialog = false;
};
const openDialog = async (row?: any) => {
	state.tableData.data = row.roles;

	dialog.isShowDialog = true;
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>
