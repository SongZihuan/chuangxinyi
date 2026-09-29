<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="869px" destroy-on-close @close="closeDialog">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-row>
						<!--        // 来源网站-->
						<el-col :span="24">
							<el-form-item label="来源网站">
								<el-select v-model="webID" placeholder="请选择网站来源">
									<el-option label="全部" :value="0"></el-option>
									<el-option v-for="item in websiteList" :key="item.id" :label="item.name" :value="item.id"></el-option>
								</el-select>
							</el-form-item>
						</el-col>
						<el-form-item>
							<el-button type="danger" @click="onDelToken">
								<el-icon>
									<DeleteFilled />
								</el-icon>
								清除在线用户
							</el-button>
						</el-form-item>
						<el-form-item>
							<el-button type="danger" @click="clearAllToken">
								<el-icon>
									<DeleteFilled />
								</el-icon>
								清除外站所有在线用户
							</el-button>
						</el-form-item>
						<el-form-item>
							<el-button type="danger" @click="onDelFatherToken">
								<el-icon>
									<DeleteFilled />
								</el-icon>
								清除所有在线父用户
							</el-button>
						</el-form-item>
					</el-row>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="tokenType" label="token类型" show-overflow-tooltip width="100" align="left">
          <template #default="scope">
            <span v-if="scope.row.tokenType === 1"> 普通登录 </span>
            <span v-if="scope.row.tokenType === 2"> 父用户登录 </span>
            <span v-if="scope.row.tokenType === 3"> 外站授权 </span>
          </template>
        </el-table-column>
        <el-table-column prop="tokenType" label="是否在线" show-overflow-tooltip width="100" align="left">
          <template #default="scope">
            <el-tag v-if="scope.row.isLogin" type="success"> 在线 </el-tag>
            <el-tag v-else type="warning"> 下线 </el-tag>
          </template>
        </el-table-column>
				<el-table-column prop="ip" label="ip地址" show-overflow-tooltip width="150" align="center"></el-table-column>
				<el-table-column prop="geo" label="地理位置" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column prop="nowIP" label="当前ip地址" show-overflow-tooltip width="150" align="center"></el-table-column>
				<el-table-column prop="nowGeo" label="当前地理位置" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="father" label="父用户ID" show-overflow-tooltip align="center" width="90">
          <template #default="scope">
            <span v-if="scope.row.tokenType === 2"> {{ scope.father.id }} </span>
            <span v-else>- </span>
          </template>
        </el-table-column>
        <el-table-column prop="father" label="父用户" show-overflow-tooltip align="center" width="90">
          <template #default="scope">
            <span v-if="scope.row.tokenType === 2"> {{ scope.father.phone }} </span>
            <span v-else>- </span>
          </template>
        </el-table-column>
        <el-table-column prop="webName" label="授权外站" show-overflow-tooltip width="150" align="left">
          <template #default="scope">
            <span v-if="scope.row.tokenType === 3"> {{ scope.row.webName }} </span>
            <span v-else>- </span>
          </template>
        </el-table-column>
				<!--        取消授权-->
				<el-table-column label="操作" width="90" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="onTabelRowDel(scope.row)">下线</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-card>
	</el-dialog>
</template>

<script setup lang="ts" name="tokenList">
import { reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { userApi } from '/@/api/system/user';
import { websiteListType } from '/@/views/oauth2/authorizationList/types';
import { useoauth2Api } from '/@/api/oauth2';
import { DeleteFilled } from '@element-plus/icons-vue';

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '用户在线列表',
	submitTxt: '',
});
const state = reactive({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		uid: 0,
		param: {
			page: 1,
			pagesize: 10,
		},
	},
});
const websiteList = ref<websiteListType>([]);
const webID = ref<number>(0);
// 获取网站列表
const getWebsiteList = () => {
	useoauth2Api()
		.oauth2DomainList({ page: 1, pagesize: 10000 })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				websiteList.value = res.data.website;
			}
		});
};
//搜索
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	userApi()
		.getUserOnlineToken({ uid: state.tableData.uid })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.token;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const userApiCollect = userApi();
// 删除父用户Token
const onDelFatherToken = () => {
	ElMessageBox.confirm('确定要删除父用户Token吗？', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			userApiCollect.deleteUserOnlineFatherUser({ uid: state.tableData.uid }).then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('删除成功');
				}
			});
		})
		.catch(() => {});
};
// 删除用户所有Token
const onDelToken = () => {
	ElMessageBox.confirm('确定要删除用户所有Token吗？', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			userApiCollect.deleteUserOnlineToken({ uid: state.tableData.uid }).then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('删除成功');
				}
			});
		})
		.catch(() => {});
};
const clearAllToken = () => {
	ElMessageBox.confirm('确定要将该用户所有用户都下线吗？', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			userApi()
				.deleteUserOnlineTokenBySite({ webID: webID.value, uid: state.tableData.uid })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('清除成功');
						getTableData();
					}
				});
		})
		.catch(() => {});
};
// 取消授权
const onTabelRowDel = (row: any) => {
	ElMessageBox.confirm('确定要将该用户踢下线吗？', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			userApi()
				.kickUserOnlineToken({ token: row.deleteToken })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('取消成功');
						getTableData();
					}
				});
		})
		.catch(() => {});
};
const openDialog = async (row?: any) => {
	state.tableData.uid = row.id;
	dialog.isShowDialog = true;
	getWebsiteList();
	getTableData();
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss"></style>
