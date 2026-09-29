<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="webID" label="网站ID" show-overflow-tooltip width="100" align="left"></el-table-column>
				<el-table-column prop="webName" label="网站名称" show-overflow-tooltip min-width="200" align="center"></el-table-column>
				<el-table-column prop="webDomain" label="网站域名" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="allowLogin" label="登录网站" show-overflow-tooltip align="center" width="120">
					<template #default="scope">
						<el-tag v-if="scope.row.allowLogin" type="success">允许</el-tag>
						<el-tag v-else type="danger">禁止</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="allowDefray" label="支付" show-overflow-tooltip align="center" width="80">
					<template #default="scope">
						<el-tag v-if="scope.row.allowDefray" type="success">允许</el-tag>
						<el-tag v-else type="danger">禁止</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="allowMsg" label="通信" show-overflow-tooltip align="center" width="90">
					<template #default="scope">
						<el-tag v-if="scope.row.allowMsg" type="success">允许</el-tag>
						<el-tag v-else type="danger">禁止</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="80" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="updateBannedDialog(scope.row)">更新</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-card>
		<updateBanned ref="bannedBannedRef" @refresh="getTableData"></updateBanned>
	</div>
</template>

<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useBannedApi } from '/@/api/banned';
import { bannedDataType, bannedStateType, bannedUpdateTypes } from '/@/views/oauth2/banned/types';
import UpdateBanned from '/@/views/oauth2/banned/component/updateBanned.vue';
import { ref } from 'vue-demi';

const state = reactive<bannedStateType>({
	tableData: {
		data: [],
		loading: false,
		param: {
			limit: 10000,
		},
	},
});
const bannedBannedRef = ref();
const updateBannedDialog = (row: bannedDataType) => {
	let data: bannedUpdateTypes = {
		webID: row.webID,
		allowLogin: row.allowLogin,
		allowDefray: row.allowDefray,
		allowMsg: row.allowMsg,
	};
	bannedBannedRef.value.openDialog(data);
};
const getTableData = () => {
	state.tableData.loading = true;
	useBannedApi()
		.getBannedList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.record;
				state.tableData.loading = false;
			}
		});
};

onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss"></style>
