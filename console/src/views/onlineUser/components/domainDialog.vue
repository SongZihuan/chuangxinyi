<template>
	<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="869px" @close="closeDialog" destroy-on-close>
		<el-descriptions :column="column" border>
			<el-descriptions-item>
				<template #label>
					<div class="cell-item">
						<el-icon class="mr2"> </el-icon>
						网站ID
					</div>
				</template>
				{{ state.domainInfo.webID }}
			</el-descriptions-item>
			<el-descriptions-item>
				<template #label>
					<div class="cell-item">
						<el-icon class="mr2"> </el-icon>
						网站名称
					</div>
				</template>
				{{ state.domainInfo.webName }}
			</el-descriptions-item>
			<el-descriptions-item>
				<template #label>
					<div class="cell-item">
						<el-icon class="mr2"> </el-icon>
						网站域名
					</div>
				</template>
				{{ state.domainInfo.webDomain }}
			</el-descriptions-item>
		</el-descriptions>
	</el-dialog>
</template>
<script setup lang="ts">
import { reactive } from 'vue';

const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '站点信息查看',
		submitTxt: '',
	},
	domainInfo: {} as {
		webDomain: string;
		webID: number;
		webName: string;
	},
});
const column = 1;
const openDialog = (row?: { webDomain: string; webID: number; webName: string }) => {
	if (!row) return;
	state.domainInfo = JSON.parse(JSON.stringify(row));
	state.dialog.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
	state.dialog.isShowDialog = false;
};
defineExpose({
	openDialog,
	closeDialog,
});
</script>
<style scoped lang="scss"></style>
