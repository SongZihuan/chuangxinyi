<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item v-if="authUser.deleteAllToken">
						<el-button type="warning" @click="clearOtherToken">
							<el-icon>
								<DeleteFilled />
							</el-icon>
							清除其他用户token
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="tokenType" label="Token类型" show-overflow-tooltip width="100" align="left">
					<template #default="scope">
						<el-tag v-if="scope.row.tokenType == 1" type="success">用户</el-tag>
						<el-tag v-else-if="scope.row.tokenType == 2" type="warning">父用户</el-tag>
						<el-tag v-else-if="scope.row.tokenType == 3" type="info">外站授权</el-tag>
					</template>
				</el-table-column>
				<!--        IsSelf为true代表自己-->
				<el-table-column prop="isSelf" label="是否自己" width="100" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.isSelf" type="success">是</el-tag>
						<el-tag v-else type="danger">否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="ip" label="IP" width="200" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="geo" label="地理位置" width="200" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="nowIP" label="当前IP" min-width="100" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="nowGeo" label="当前地理位置" width="200" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="father" label="所属用户" width="100" show-overflow-tooltip align="center">
					<template #default="scope">
						<!--            查看父级信息-->
						<el-button v-if="scope.row.tokenType == 2" type="text" @click="checkFatherInfo(scope.row)"> 查看 </el-button>
						<el-tag v-else type="info">-</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="webID" label="站点信息" width="100" show-overflow-tooltip align="center">
					<template #default="scope">
						<!--            查看站点信息-->
						<el-button v-if="scope.row.tokenType == 3" type="text" @click="checkWebInfo(scope.row)"> 查看 </el-button>
						<el-tag v-else type="info">-</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="isLogin" label="是否在线" width="100" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.isLogin" type="success">在线</el-tag>
						<el-tag v-else type="danger">离线</el-tag>
					</template>
				</el-table-column>
				<!--        踢下线-->
				<el-table-column label="操作" width="110" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="onTabelRowDel(scope.row)">下线</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-card>
		<father-info ref="fatherInfoRef" />
		<domain-dialog ref="domainDialogRef" />
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import useOnlineUser from '/@/api/onlineUser';
import { DeleteFilled } from '@element-plus/icons-vue';
import useSubAuth from '/@/hooks/useSubAuth';
import FatherInfo from '/@/views/onlineUser/components/fatherInfo.vue';
import DomainDialog from '/@/views/onlineUser/components/domainDialog.vue';

const fatherInfoRef = ref();
const authUser = useSubAuth();
const domainDialogRef = ref();
const state = reactive({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			signin: false,
		},
	},
});
//搜索
const clearOtherToken = () => {
	ElMessageBox.confirm('确定要退出除自己之外的所有在线用户token吗？', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useOnlineUser()
				.deleteAllOnlineUser()
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('清除成功');
						getTableData();
					}
				});
		})
		.catch(() => {});
};
// 查看父级信息
const checkFatherInfo = (row: any) => {
	fatherInfoRef.value.openDialog(row.father);
};
// 查看站点信息
const checkWebInfo = (row: any) => {
	domainDialogRef.value.openDialog({
		webID: row.webID,
		webName: row.webName,
		webDomain: row.webDomain,
	});
};
// 下线
const onTabelRowDel = (row: any) => {
	ElMessageBox.confirm('确定要下线该用户吗？', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useOnlineUser()
				.kickOnlineUser({ token: row.deleteToken })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('踢下线成功');
						getTableData();
					}
				});
		})
		.catch(() => {});
};
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	useOnlineUser()
		.getOnlineUserList()
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.token;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};

// 页面加载时
onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss"></style>
