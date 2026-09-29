<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
						<el-select v-model="state.tableData.param.timetype" placeholder="请选择日期类型" style="width: 100%">
							<el-option v-for="(item, index) in dateTypeDict" :value="item.value" :label="item.label" :key="index" />
						</el-select>
					</el-form-item>
					<el-form-item>
						<el-date-picker
							v-model="state.tableData.param.range"
							type="datetimerange"
							range-separator="至"
							start-placeholder="开始日期"
							end-placeholder="结束日期"
							style="width: 100%"
						/>
					</el-form-item>
					<el-form-item>
						<el-button type="primary" class="ml10" @click="handleSearch">
							<el-icon>
								<ele-Search />
							</el-icon>
							查询
						</el-button>
						<!--            一键已读-->
						<el-button type="success" class="ml10" @click="onRead">
							<SvgIcon name="ele-ChatLineRound"></SvgIcon>
							一键已读
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="title" label="标题" show-overflow-tooltip align="center"></el-table-column>
				<el-table-column prop="sender" label="发送者" show-overflow-tooltip align="center">
					<template #default="scope">
						<span :class="{link:scope.row.senderLink}" v-if="scope.row.sender" @click="onSkip({ link: scope.row.senderLink })">{{ scope.row.sender }}</span>
					</template>
				</el-table-column>
				<el-table-column prop="readAt" label="阅读时间" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.readAt" type="success">{{ dayjs.unix(scope.row.readAt).format('YYYY-MM-DD HH:mm:ss') }}</el-tag>
						<el-tag v-else type="info">未阅读</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="createAt" label="发送时间" show-overflow-tooltip align="center">
					<template #default="scope">
						<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="60" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" size="mini" @click="handleView(scope.row)">查看</el-button>
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
		<Detail-dialog ref="detailDialogRef" @refresh="getTableData"></Detail-dialog>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, watch, ref } from 'vue';
import dayjs from 'dayjs';
import { messageUserDataType, messageUserState } from './types';
import { useMessageUserApi } from '/@/api/message/user';
import DetailDialog from '/@/views/message/user/component/detailDialog.vue';
import { ElMessage } from 'element-plus';
import useSkip from '/@/hooks/useSkip';
const state = reactive<messageUserState>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			starttime: 0,
			endtime: 0,
			range: [],
			timetype: undefined,
      senderID: -1,
    },
	},
});
const dateTypeDict = ref([
	{ label: '创建时间', value: 1 },
	{ label: '读取时间', value: 6 },
]);
const detailDialogRef = ref();
const { onSkip } = useSkip();
watch(
	() => state.tableData.param.range,
	() => {
		if (state.tableData.param.range && state.tableData.param.range.length > 0) {
			state.tableData.param.starttime = dayjs(state.tableData.param.range[0]).unix();
			state.tableData.param.endtime = dayjs(state.tableData.param.range[1]).unix();
		}
	},
	{
		deep: true,
	}
);
// 查看
const handleView = (row: messageUserDataType) => {
	detailDialogRef.value.openDialog(row);
};
// 已读
const onRead = () => {
	useMessageUserApi()
		.readAllMessage()
		.then((res: any) => {
			if (res) {
				ElMessage.success('一键已读成功');
				getTableData();
			}
		});
};
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	// 站内信
	useMessageUserApi()
		.getMessageList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.message;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
const handleSearch = () => {
	getTableData();
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
.link {
	color: #409eff;
	cursor: pointer;
}
</style>
