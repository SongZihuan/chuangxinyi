<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form @submit.native.prevent :inline="true" :model="state.tableData.param" class="demo-form-inline">
					<el-form-item>
						<el-input v-model="state.tableData.param.domain" size="default" placeholder="请输入域名" style="max-width: 180px" clearable></el-input>
					</el-form-item>
					<el-form-item>
						<el-button type="primary" class="ml10" @click="handleSearch">
							<el-icon>
								<ele-Search />
							</el-icon>
							查询
						</el-button>
						<el-button type="success" class="ml10" @click="openDialog('add')">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增
						</el-button>
						<!--            从数据库更新站点-->
						<el-button type="success" class="ml10" @click="updateSite">
							<el-icon>
								<RefreshRight />
							</el-icon>
							从数据库更新站点
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
        <el-table-column prop="uid" label="站点ID" show-overflow-tooltip align="center"></el-table-column>
        <el-table-column prop="name" label="名称" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="describe" label="描述" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="createTime" label="创建时间" show-overflow-tooltip align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="pubkey" label="公钥" align="center" width="60">
					<template #default="scope">
						<el-tooltip effect="dark" content="点击公钥查看" placement="top">
							<el-button text type="primary" @click="showPubkey(scope.row.pubkey)">公钥</el-button>
						</el-tooltip>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="200" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="openDialog('edit', scope.row)">编辑</el-button>
						<el-button text type="primary" @click="showIpList(scope.row)">编辑IP</el-button>
						<el-button text type="primary" @click="onEditDomainList(scope.row)">编辑域名</el-button>
						<el-button text type="primary" @click="updateSecret(scope.row)">更新公钥</el-button>
						<el-button text type="primary" @click="onTabelRowDel(scope.row)">删除</el-button>
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
				:total="state.tableData.total"
			>
			</el-pagination>
			<!-- 新增编辑弹窗 -->
			<ActionDialog ref="actionDialog" @refresh="getTableData()" />
			<IpDialog ref="ipDialogRef" @refresh="getTableData" />
			<DomainDialog ref="domainDialogRef" @refresh="getTableData" />
			<el-dialog v-model="dialogVisible" title="公钥" width="650px">
				<el-text type="primary" @click="handleCopy">一键复制</el-text>
				<el-input v-model="pubKey" class="mt10" :autosize="{ minRows: 2, maxRows: 10 }" readonly type="textarea"></el-input>
			</el-dialog>
		</el-card>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, defineAsyncComponent, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useSiteApi } from '/@/api/website/site/index';
import type { websiteListTypes } from '/@/api/website/site/types';
import commonFunction from '/@/utils/commonFunction';
import dayjs from 'dayjs';
import { Base64 } from 'js-base64';
import { RefreshRight } from '@element-plus/icons-vue';

const ActionDialog = defineAsyncComponent(() => import('./component/actionDialog.vue'));
const IpDialog = defineAsyncComponent(() => import('./component/ipDialog.vue'));
const DomainDialog = defineAsyncComponent(() => import('./component/domainDialog.vue'));
const useSiteApiCollect = useSiteApi();
const actionDialog = ref();
const ipDialogRef = ref();
const domainDialogRef = ref();
const pubKey = ref<any>('');
const domainStatus = ref<number |null>(null)
const ipStatus = ref<number |null>(null)
const rowData = ref<any>(null)
const { copyText } = commonFunction();
const state = reactive<websiteListTypes>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			domain: '',
		},
	},
});
const dialogVisible = ref();
// 初始化表格数据
const getTableData = (status:boolean=false,type:string='') => {
	state.tableData.loading = true;
	useSiteApiCollect.siteList(state.tableData.param).then((res: any) => {
		if (res.code === "SUCCESS") {
			state.tableData.data = res.data.website;
			state.tableData.total = res.data.count;
			state.tableData.loading = false;
      if(status){
        if(type==='domain'){
          rowData.value = state.tableData.data.find((item:any)=>item.id===domainStatus.value)
          onEditDomainList(rowData.value)
        }else{
          rowData.value = state.tableData.data.find((item:any)=>item.id===ipStatus.value)
          showIpList(rowData.value)
        }
      }
		}
	});
};
const showIpList = (row: any) => {
	ipDialogRef.value.openDialog(row);
  ipStatus.value=row.id
};
const onEditDomainList = (row: any) => {
	domainDialogRef.value.openDialog(row);
  domainStatus.value=row.id
};
const showPubkey = (pubkey: string) => {
	pubKey.value = Base64.decode(pubkey);
	dialogVisible.value = true;
};
//搜索
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
const openDialog = (type: string, row?: any) => {
	actionDialog.value.openDialog(type, row);
};
const handleCopy = () => {
	copyText(pubKey.value);
};
const onTabelRowDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除站点：${row.name}, 是否继续?`, '提示', {
		confirmButtonText: '删除',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useSiteApiCollect.delSite({ id: row.id }).then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('删除成功');
					getTableData();
				}
			});
		})
		.catch(() => {});
};
const updateSite = () => {
	ElMessageBox.confirm(`此操作将从数据库更新站点，是否继续?`, '提示', {
		confirmButtonText: '更新',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useSiteApiCollect.updateSiteFromDB().then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('更新成功');
					getTableData();
				}
			});
		})
		.catch(() => {});
};
const updateSecret = (row: any) => {
	ElMessageBox.prompt('请输入公钥', '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		inputErrorMessage: '请输入公钥',
		inputType: 'textarea',
	})
		.then(({ value }) => {
			if (value != null) {
				useSiteApiCollect.updateSiteSecret({ id: row.id, pubkey: "base64:" + Base64.encode(value) }).then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage({
							type: 'success',
							message: '更新公钥成功',
						});
					}
				});
			} else {
				ElMessage({
					type: 'warning',
					message: '请输入公钥',
				});
			}
		})
		.catch(() => {});
};
// 分页改变
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
	getTableData();
};
// 页面加载时
onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss">
:deep(.el-textarea__inner) {
	font-family: 'Courier New', Courier, monospace;
}
</style>
