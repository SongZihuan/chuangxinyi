<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="869px" destroy-on-close @close="closeDialog">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%">
				<el-table-column prop="webID" label="网站ID" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column prop="webName" label="网站名称" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="allowLogin" label="允许开通网站" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.allowLogin" type="success">允许</el-tag>
						<el-tag v-else type="danger">禁止</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="allowDefray" label="允许支付" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.allowDefray" type="success">允许</el-tag>
						<el-tag v-else type="danger">禁止</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="allowMsg" label="通信" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.allowMsg" type="success">允许</el-tag>
						<el-tag v-else type="danger">禁止</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="220" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="updateBannedDialog(scope.row)">更新</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-card>
		<updateBanned ref="bannedBannedRef" @refresh="getTableData"></updateBanned>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { bannedDataType, bannedStateType, bannedUpdateTypes } from '/@/views/oauth2/banned/types';
import UpdateBanned from '/@/views/oauth2/banned/component/updateBanned.vue';
import { ref } from 'vue-demi';
import { userApi } from '/@/api/system/user';

const state = reactive<bannedStateType>({
	tableData: {
		data: [],
		loading: false,
		param: {
			limit: 10000,
      uid: 0,
		},
	},
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '用户开通列表',
	submitTxt: '',
});
const closeDialog = () => {
	dialog.isShowDialog = false;
};

const openDialog = async (row?: any) => {
	state.tableData.param.uid = row.id;
	dialog.isShowDialog = true;
	getTableData();
};
const bannedBannedRef = ref();
const updateBannedDialog = (row: bannedDataType) => {
	let data: bannedUpdateTypes = {
		id: state.tableData.param.uid,
		webID: row.webID,
		allowLogin: row.allowLogin,
		allowDefray: row.allowDefray,
		allowMsg: row.allowMsg,
	};
	bannedBannedRef.value.openDialog(data, 'admin');
};
const getTableData = () => {
	state.tableData.loading = true;
	userApi()
		.getUserBanRecord({ uid: state.tableData.param.uid, ...state.tableData.param })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.record;
				state.tableData.loading = false;
			}
		});
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss"></style>
