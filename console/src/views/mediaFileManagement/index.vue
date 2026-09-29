<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-form :inline="true" class="demo-form-inline">
					<el-form-item>
						<el-button type="success" class="ml10" @click="openDialog()">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column label="文件名" width="350px" show-overflow-tooltip align="center">
					<template #default="scope">
						{{ scope.row.fid }}
					</template>
				</el-table-column>
				<el-table-column label="MIME类型" show-overflow-tooltip align="center">
					<template #default="scope">
						{{ scope.row.mediaType }}
					</template>
				</el-table-column>
        <el-table-column label="预览" max-width="200px" align="center">
          <template #default="scope">
            <el-image
                v-if="isImage(scope.row.mediaType)"
                class="img"
                :src="getFile(scope.row.fid)"
                :preview-src-list="[getFile(scope.row.fid)]"
                :preview-teleported="true"
                fit="cover"
            />
            <vue3-video-player
                v-else-if="isVideo(scope.row.mediaType)"
                :src="getFile(scope.row.fid)"
            />
            <audio-player
                v-else-if="isAudio(scope.row.mediaType)"
                :audio-list="[getFile(scope.row.fid)]"
            />
            <div v-else> 类型：{{ scope.row.mediaType }} 暂时不支持预览</div>
          </template>
        </el-table-column>
				<el-table-column label="操作" width="180" align="center" fixed="right">
					<template #default="scope">
						<el-button text type="primary" @click="onDoawload(scope.row)">下载</el-button>
						<!--   覆盖         -->
						<el-button text type="primary" @click="editDialog(scope.row)">覆盖</el-button>
						<el-button text type="primary" @click="onRowDel(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
			<el-pagination hide-on-single-page
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				v-model:current-page="state.tableData.param.page"
				background
				v-model:page-size="state.tableData.param.pagesize"
				layout="total, prev, pager, next, jumper"
				:total="state.tableData.total"
			>
			</el-pagination>
		</el-card>
		<actionDialog ref="fileActionDialog" @refresh="getTableData" />
	</div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { fileStatsTypes } from '/@/views/mediaFileManagement/types';
import { useMediaFileApi } from '/@/api/mediaFileManagement';
import actionDialog from '/@/views/mediaFileManagement/component/actionDialog.vue';
import { Session } from '/@/utils/storage';
import useDownload from '/@/hooks/useDownload';
import commonFunction from '/@/utils/commonFunction';
import useFile from "/@/hooks/useFile";
const state = reactive<fileStatsTypes>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 5,
			timetype: 1,
			starttime: 0,
			endtime: 0,
			range: [],
		},
	},
});
const isImage = (mimeType: string) => {
  return mimeType.startsWith("image/")
}

const isVideo = (mimeType: string) => {
  return mimeType.startsWith("video/")
}

const isAudio = (mimeType: string) => {
  return mimeType.startsWith("audio/")
}

const fileActionDialog = ref();
// 打开
const openDialog = () => {
	fileActionDialog.value.openDialog({}, 'add');
};
// 覆盖
const editDialog = (row: { fid: string }) => {
	fileActionDialog.value.openDialog(row, 'edit');
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
	getTableData();
};
const getTableData = () => {
	state.tableData.loading = true;
	useMediaFileApi()
		.getMediaFileList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.file;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
const apiUrl = ref();
const { downFile } = useDownload();
const onDoawload = (row: any) => {
	const url = getFile(row.fid);
	downFile(url, true);
};
const {getFile} = useFile();
const onRowDel = (row: any) => {
	ElMessageBox.confirm('此操作将永久删除该图片, 是否继续?', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useMediaFileApi()
				.deleteMediaFile({ fid: row.fid })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
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
	apiUrl.value = import.meta.env.VITE_API_URL;
	getTableData();
});
</script>

<style scoped lang="scss">
.img {
	width: 100px;
	height: 100px;
	object-fit: cover;
	border-radius: 5px;
	cursor: pointer;
}
</style>
