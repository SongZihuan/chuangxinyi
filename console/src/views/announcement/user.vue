<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="title" label="标题" show-overflow-tooltip align="left"></el-table-column>
				<el-table-column prop="content" label="公告内容" show-overflow-tooltip align="center">
					<template #default="scope"> <el-button text type="primary" @click="showContent(scope.row)">公告内容</el-button> </template>
				</el-table-column>
				<el-table-column prop="startAt" label="公告时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.startAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
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
		.userAnnouncementList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.announcement;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};

const showContent = (row: any) => {
	currentInfo.value = row.content;
	showHtmlRef.value.openDialog();
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
