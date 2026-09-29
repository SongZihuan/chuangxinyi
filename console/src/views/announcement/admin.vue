<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-button size="default" type="success" @click="onOpenAddAnnouncement('add')">
					<el-icon>
						<ele-FolderAdd />
					</el-icon>
					新增公告
				</el-button>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="title" label="标题" show-overflow-tooltip min-width="180" align="left"></el-table-column>
				<el-table-column prop="content" label="公告内容" show-overflow-tooltip align="center">
					<template #default="scope"> <el-button text type="primary" @click="showContent(scope.row)">公告内容</el-button> </template>
				</el-table-column>
				<el-table-column prop="startAt" label="开始时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.startAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="stopAt" label="结束时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.stopAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="180" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="roleMoveAnnouncement(scope.row, 'top')">上移</el-button>
						<el-button text type="primary" @click="roleMoveAnnouncement(scope.row, 'down')">下移</el-button>
						<el-button :disabled="scope.row.roleName === '超级管理员'" text type="primary" @click="onOpenEditAnnouncement('edit', scope.row)"
							>修改</el-button
						>
						<el-button :disabled="scope.row.roleName === '超级管理员'" text type="primary" @click="onRowDel(scope.row)">删除</el-button>
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
		</el-card>
		<action-dialog ref="announcementDialogRef" @refresh="getTableData()" />
		<ShowHtml ref="showHtmlRef" title="公告内容" :contentHtml="currentInfo" />
	</div>
</template>

<script setup lang="ts" name="systemAnnouncement">
import { reactive, onMounted, ref } from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import dayjs from 'dayjs';
import { useAnnouncementApi } from '/@/api/announcement';
import ActionDialog from '/@/views/announcement/component/actionDialog.vue';
import { AnnouncementState } from '/@/views/announcement/types';
import ShowHtml from '/@/components/ShowHtml/index.vue';
// 定义变量内容
const announcementDialogRef = ref();
const showHtmlRef = ref();
const currentInfo = ref();
const state = reactive<AnnouncementState>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
		},
	},
});
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	useAnnouncementApi()
		.getAnnouncementList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.announcement;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
//移动菜单
const roleMoveAnnouncement = (row: { id: number }, type: string) => {
	useAnnouncementApi()
		.moveAnnouncement({ id: row.id, isUp: type === 'top' })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				type === 'top' ? ElMessage.success('上移成功') : ElMessage.success('下移成功');
				getTableData();
			}
		});
};
const showContent = (row: any) => {
	currentInfo.value = row.content;
	showHtmlRef.value.openDialog();
};
// 打开新增公告弹窗
const onOpenAddAnnouncement = (type: string) => {
	announcementDialogRef.value.openDialog(type);
};
// 打开修改公告弹窗
const onOpenEditAnnouncement = (type: string, row: Object) => {
	announcementDialogRef.value.openDialog(type, row);
};
// 删除公告
const onRowDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除公告名称：“${row.title}”，是否继续?`, '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useAnnouncementApi()
				.delAnnouncement({ id: row.id })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('删除成功');
					} else {
						ElMessage.error('删除失败');
					}
					getTableData();
				});
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

<style scoped lang="scss"></style>
