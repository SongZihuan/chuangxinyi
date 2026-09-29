<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<!-- 用户登录状态筛选-->
					<el-form-item label="用户登录状态">
						<el-select v-model="state.tableData.param.isLogin" placeholder="请选择用户登录状态" @change="getTableData">
							<el-option label="全部" :value="0"></el-option>
							<el-option label="已登录" :value="1"></el-option>
							<el-option label="未登录" :value="2"></el-option>
						</el-select>
					</el-form-item>
					<!--  来源网站-->
					<el-form-item label="来源网站">
						<el-select v-model="webID" placeholder="请选择网站来源" @change="getSelectWebsite">
							<el-option label="全部" :value="0"></el-option>
							<el-option v-for="item in websiteList" :key="item.id" :label="item.name" :value="item.id"></el-option>
						</el-select>
					</el-form-item>
					<el-form-item>
						<el-button type="danger" @click="clearAllAuth">
							<el-icon>
								<DeleteFilled />
							</el-icon>
							清除所有用户授权
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<!--        isLogin: false
webID: 10000-->
				<el-table-column prop="webID" label="用户ID" show-overflow-tooltip width="150" align="center"></el-table-column>
				<el-table-column prop="webName" label="网站名称" show-overflow-tooltip align="left"></el-table-column>
				<el-table-column prop="isLogin" label="是否登录" show-overflow-tooltip align="center">
					<template #default="scope">
						<span v-if="scope.row.isLogin">是</span>
						<span v-else>否</span>
					</template>
				</el-table-column>
				<el-table-column prop="ip" label="ip地址" show-overflow-tooltip align="left"></el-table-column>
				<el-table-column prop="geo" label="地理位置" show-overflow-tooltip align="center"></el-table-column>
				<!--        取消授权-->
				<el-table-column label="操作" width="220" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="onTabelRowDel(scope.row)">取消授权</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-card>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { useoauth2Api } from '/@/api/oauth2';
import { ElMessage, ElMessageBox } from 'element-plus';
import { ref } from 'vue-demi';
import { websiteListType } from './types';
import { DeleteFilled } from '@element-plus/icons-vue';

const state = reactive({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 10,
			isLogin: 1,
			limit: 0,
		},
	},
});
const websiteList = ref<websiteListType>([]);
const webID = ref<number>(0);
const tableData = ref<any>([]);
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

// 初始化表格数据
const getTableData = async () => {
	state.tableData.loading = true;
	await useoauth2Api()
		.oauth2List(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				tableData.value = res.data.record;
				state.tableData.data = res.data.record;
				state.tableData.total = res.data.count;

        if (state.tableData.param.isLogin == 1) {
          state.tableData.data = tableData.value.filter((item) => item.isLogin == true);
        } else if (state.tableData.param.isLogin == 2) {
          state.tableData.data = tableData.value.filter((item) => item.isLogin == false);
        }

        if (webID.value !== 0) {
          state.tableData.data = tableData.value.filter((item: any) => item.webID == webID.value);
        }
        state.tableData.total = state.tableData.data.length;

        state.tableData.loading = false;
      }
		});
};

// 取消所有授权
const clearAllAuth = () => {
	const webName = webID.value == 0 ? '全部' : websiteList.value.find((item: any) => item.id == webID.value)?.name;
	ElMessageBox.confirm(`确定要取消${webName}的所有授权吗，此操作将导致所有授权清空？`, '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useoauth2Api()
				.oauth2Clear({ webID: webID.value })
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
	ElMessageBox.confirm('确定要取消授权吗？', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useoauth2Api()
				.oauth2Cancel({ token: row.deleteToken })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('取消成功');
            getTableData();
					}
				});
		})
		.catch(() => {});
};

// 页面加载时
onMounted(() => {
  getTableData();
	getWebsiteList();
});
</script>

<style scoped lang="scss"></style>
