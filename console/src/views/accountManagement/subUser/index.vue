<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-form :inline="true" class="demo-form-inline">
					<el-form-item v-if="authUser.registerSon">
						<el-button type="success" @click="openDialog()">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table
				:data="state.tableData.data"
				v-loading="state.tableData.loading"
				style="width: 100%"
				row-key="id"
				default-expand-all
				:tree-props="{ children: 'son', hasChildren: 'hasChildren' }"
			>
				<el-table-column prop="id" label="用户ID" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="phone" label="手机号" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="roleName" label="角色" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="lineal" label="关系" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.lineal">子账号</el-tag>
						<el-tag type="info" v-else>协作账号</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="status" label="状态" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.status == 'NORMAL'">正常</el-tag>
						<el-tag type="info" v-else>禁用</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="inviteCount" label="邀请人数" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="lastInviteAt" label="最后邀请时间" show-overflow-tooltip width="180" align="center">
					<template #default="scope">
						<span v-if="scope.row.lastInviteAt">{{ dayjs.unix(scope.row.lastInviteAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
						<span v-else>-</span>
					</template>
				</el-table-column>
				<el-table-column prop="createAt" label="创建时间" show-overflow-tooltip align="center">
					<template #default="scope">
						<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
			</el-table>
		</el-card>
		<action-dialog ref="uncleActionDialog" @refresh="getTableData" />
	</div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { subUserAccountStatsTypes } from '/@/views/accountManagement/subUser/types';
import { useSubUserCenterApi } from '/@/api/subUser';
import ActionDialog from '/@/views/accountManagement/subUser/component/actionDialog.vue';
import dayjs from 'dayjs';
import useSubAuth from '/@/hooks/useSubAuth';

const state = reactive<subUserAccountStatsTypes>({
	tableData: {
		data: [],
		loading: false,
		param: {
			id: '',
			uid: '',
		},
	},
});
const uncleActionDialog = ref();
const authUser = useSubAuth();
const getTableData = () => {
	state.tableData.loading = true;
	useSubUserCenterApi()
		.getSubUserList()
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				state.tableData.data = res.data.user;
				state.tableData.loading = false;
			}
		});
};
// 打开
const openDialog = () => {
	uncleActionDialog.value.openDialog();
};

onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss"></style>
