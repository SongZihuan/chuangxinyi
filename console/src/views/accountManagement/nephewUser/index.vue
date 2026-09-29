<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-form :inline="true" class="demo-form-inline">
					<el-form-item>
						<el-button type="success" @click="openDialog()">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border auto>
				<el-table-column prop="id" label="用户ID" show-overflow-tooltip align="left"></el-table-column>
				<el-table-column prop="roleName" label="角色" align="center" min-width="170"></el-table-column>
				<el-table-column prop="phone" label="手机号" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="userName" label="用户名" show-overflow-tooltip align="center" min-width="170"></el-table-column>
				<el-table-column prop="nickName" label="昵称" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="status" label="状态" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.uncleStatus === 1">待确认</el-tag>
						<el-tag type="success" v-else-if="scope.row.status === 'NORMAL'">正常</el-tag>
						<el-tag type="info" v-else>禁用</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="150" align="center" fixed="right" v-if="authUser.delUncle">
					<template #default="scope">
<!--            账号确认-->
            <el-button text type="primary" @click="onRowConfirm(scope.row)" v-if="scope.row.uncleStatus === 1">确认</el-button>
						<el-button text type="primary" @click="onRowDel(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-card>
		<action-dialog ref="uncleActionDialog" @refresh="getTableData"></action-dialog>
	</div>
</template>

<script setup lang="ts">
import { defineAsyncComponent, onMounted, reactive, ref } from 'vue';
import { uncleAccountStatsTypes } from '/@/views/accountManagement/uncleUser/types';
import { ElMessage, ElMessageBox } from 'element-plus';
import useSubAuth from '/@/hooks/useSubAuth';
import { useNephewUserApi } from '/@/api/nephewUser';
const ActionDialog = defineAsyncComponent(() => import('/@/views/accountManagement/uncleUser/component/actionDialog.vue'));

const authUser = useSubAuth();
const state = reactive<uncleAccountStatsTypes>({
	tableData: {
		data: [],
		loading: false,
	},
});
const uncleActionDialog = ref();
// 打开
const openDialog = () => {
	uncleActionDialog.value.openDialog();
};
const getTableData = () => {
	state.tableData.loading = true;
	useNephewUserApi()
		.getNephewUserInfo()
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				state.tableData.data = res.data.user;
				state.tableData.loading = false;
			}
		});
};
const onRowConfirm = (row: any) => {
  ElMessageBox.confirm('此操作将确认该用户, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      useNephewUserApi()
        .addNephewUser({ nephewID: row.id })
        .then((res: any) => {
          if (res.code === 'SUCCESS') {
            ElMessage.success('确认成功');
            getTableData();
          }
        });
    })
    .catch(() => {
      ElMessage.info('已取消确认');
    });
};
const onRowDel = (row: any) => {
	ElMessageBox.confirm('此操作将永久删除该用户, 是否继续?', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useNephewUserApi()
				.delNephewUser({ nephewID: row.id })
				.then((res: any) => {
					if (res.code === 'SUCCESS') {
						ElMessage.success('删除成功');
						getTableData();
					}
				});
		})
		.catch(() => {
			ElMessage.info('已取消删除');
		});
};
onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss"></style>
