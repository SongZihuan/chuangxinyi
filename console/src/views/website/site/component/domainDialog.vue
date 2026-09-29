<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="850px" destroy-on-close @close="closeDialog">
		<el-form :inline="true" :model="state.tableData.param" class="demo-form-inline">
			<el-form-item>
				<el-button type="success" @click="create">
					<el-icon>
						<ele-FolderAdd />
					</el-icon>
					新增
				</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
			<el-table-column prop="domain" label="域名地址" show-overflow-tooltip min-width="140" align="center"></el-table-column>
			<el-table-column label="操作" width="80" fixed="right" align="center">
				<template #default="scope">
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
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useDomainnApi } from '/@/api/website/domain/index';
import type { ipListTypes } from '/@/api/website/ip/types';
const emit = defineEmits(['refresh']);
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '登录站点域名列表',
	submitTxt: '',
});

const state = reactive<ipListTypes>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			websiteid: '',
		},
	},
});

const create = () => {
	ElMessageBox.prompt('请输入你要添加的域名地址', '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		inputPattern: /^.+$/,
		inputErrorMessage: '请填写',
	})
		.then(({ value }) => {
			useDomainnApi()
				.createdomain({ websiteID: state.tableData.param.websiteid, domain: value })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('添加域名成功');
						emit('refresh',true,'domain');
					}
				});
		})
		.catch(() => {});
};

const onTabelRowDel = (row: any) => {
	useDomainnApi()
		.deldomain({ id: row.id })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				ElMessage.success('删除域名成功');
        emit('refresh',true,'domain');
			}
		});
};
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};

const openDialog = async (row?: any) => {
	state.tableData.param.websiteid = row.id;
	state.tableData.data = row.domain;
	dialog.isShowDialog = true;
};
// 暴露变量
defineExpose({
	openDialog,
});
</script>
