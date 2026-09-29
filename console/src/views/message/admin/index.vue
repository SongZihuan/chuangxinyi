<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-tabs v-model="messageActiveName" @tab-click="handleSearch" type="border-card">
				<el-tab-pane label="站内信" name="1"
					><message-table v-if="messageActiveName == 1" :data="state" @search="handleMessageSeach" @update="getTableData"></message-table
				></el-tab-pane>
				<el-tab-pane label="手机短信" name="2"
					><phone-message-table v-if="messageActiveName == 2" :data="state" @search="handleMessageSeach" @update="getTableData"></phone-message-table
				></el-tab-pane>
				<el-tab-pane label="邮箱邮件" name="3"
					><email-message-table v-if="messageActiveName == 3" :data="state" @search="handleMessageSeach" @update="getTableData"></email-message-table
				></el-tab-pane>
				<el-tab-pane label="服务号模板消息" name="4">
					<fuwuhao-message-table v-if="messageActiveName == 4" :data="state" @search="handleMessageSeach" @update="getTableData"></fuwuhao-message-table
				></el-tab-pane>
				<el-tab-pane label="企微机器人消息" name="5"
					><wxrobot-message-table v-if="messageActiveName == 5" :data="state" @search="handleMessageSeach" @update="getTableData"></wxrobot-message-table
				></el-tab-pane>
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
			</el-tabs>
		</el-card>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { messageAdminState } from './types';
import { useMessageAdminApi } from '/@/api/message/admin';
import { TabsPaneContext } from 'element-plus';
import MessageTable from '/@/views/message/admin/component/messageTable.vue';
import PhoneMessageTable from '/@/views/message/admin/component/phoneMessageTable.vue';
import EmailMessageTable from '/@/views/message/admin/component/emailMessageTable.vue';
import FuwuhaoMessageTable from '/@/views/message/admin/component/fuwuhaoMessageTable.vue';
import WxrobotMessageTable from '/@/views/message/admin/component/wxrobotMessageTable.vue';
const messageActiveName = ref<any>('1');
const state = reactive<messageAdminState>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			timetype: 1,
			starttime: 0,
			endtime: 0,
			range: [],
			phone: '',
      senderID: -1,
		},
	},
});
watch(
	() => state.tableData.param.range,
	() => {
		if (state.tableData.param.range.length > 0) {
			state.tableData.param.starttime = dayjs(state.tableData.param.range[0]).unix();
			state.tableData.param.endtime = dayjs(state.tableData.param.range[1]).unix();
		}
	},
	{
		deep: true,
	}
);
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	if (messageActiveName.value == undefined) {
		messageActiveName.value = '1';
	}
	getTableDataByMessageActiveName();
};
const handleSearch = (tab: TabsPaneContext) => {
	messageActiveName.value = tab.props.name;
	state.tableData.param.page = 1;
	getTableData();
};
const handleMessageSeach = (param: any) => {
	state.tableData.param = param;
	getTableData();
};
// 根据messageActiveName获取不同的数据
const getTableDataByMessageActiveName = () => {
	if (messageActiveName.value == '1') {
		// 站内信
		useMessageAdminApi()
			.getMessageList(state.tableData.param)
			.then((res: any) => {
				if (res.code === "SUCCESS") {
					state.tableData.data = res.data.message;
					state.tableData.total = res.data.count;
					state.tableData.loading = false;
				}
			});
	} else if (messageActiveName.value == '2') {
		useMessageAdminApi()
			.getSmsMessageList(state.tableData.param)
			.then((res: any) => {
				if (res.code === "SUCCESS") {
					state.tableData.data = res.data.message;
					state.tableData.total = res.data.count;
					state.tableData.loading = false;
				}
			});
	} else if (messageActiveName.value == '3') {
		// 邮件站内信
		useMessageAdminApi()
			.getEmailMessageList(state.tableData.param)
			.then((res: any) => {
				if (res.code === "SUCCESS") {
					state.tableData.data = res.data.message;
					state.tableData.total = res.data.count;
					state.tableData.loading = false;
				}
			});
	} else if (messageActiveName.value == '4') {
		// 服务号站内信
		useMessageAdminApi()
			.getFuwuhaoMessageList(state.tableData.param)
			.then((res: any) => {
				if (res.code === "SUCCESS") {
					state.tableData.data = res.data.message;
					state.tableData.total = res.data.count;
					state.tableData.loading = false;
				}
			});
	} else if (messageActiveName.value == '5') {
		// 公众号站内信
		useMessageAdminApi()
			.getWxrobotMessageList(state.tableData.param)
			.then((res: any) => {
				if (res.code === "SUCCESS") {
					state.tableData.data = res.data.message;
					state.tableData.total = res.data.count;
					state.tableData.loading = false;
				}
			});
	}
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
