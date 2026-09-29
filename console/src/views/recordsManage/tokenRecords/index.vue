<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
						<el-input text v-model="state.tableData.param.token" placeholder="请输入token" />
					</el-form-item>
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
						<!--            根据token查询-->
						<el-button type="primary" class="ml10" @click="handleTokenSearch">
							<el-icon>
								<ele-Search />
							</el-icon>
							根据token查询
						</el-button>
						<el-button type="primary" class="ml10" @click="handleSearch">
							<el-icon>
								<ele-Search />
							</el-icon>
							根据时间查询
						</el-button>
						<!--            重置-->
						<el-button type="primary" class="ml10" @click="handleSearchReset">
							<el-icon>
								<ele-Search />
							</el-icon>
							重置
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="tokenType" label="token类型" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<span v-if="scope.row.tokenType == 1">用户Token</span>
						<span v-else-if="scope.row.tokenType == 2">授权Token</span>
					</template>
				</el-table-column>
				<el-table-column prop="token" label="token" align="center"></el-table-column>
				<el-table-column prop="type" label="类型" show-overflow-tooltip width="150" align="center">
					<template #default="scope">
						<span v-if="scope.row.type == 1">创建</span>
						<span v-else-if="scope.row.type == 2">地理位置变更</span>
						<span v-else-if="scope.row.type == 3">过期删除</span>
					</template>
				</el-table-column>
				<el-table-column prop="data" label="数据" show-overflow-tooltip width="80" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="showJson(scope.row.data, '数据')">查看</el-button>
					</template>
				</el-table-column>
				<el-table-column prop="createAt" label="创建时间" width="170" align="center">
					<template #default="scope">
						<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<!--        复制token-->
				<el-table-column label="操作" width="180" align="center" fixed="right">
					<template #default="scope">
						<el-button text type="primary" @click="copyText(scope.row.token)">根据token获取访问记录</el-button>
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
		<!-- json展示弹窗 -->
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" height="500px" center destroy-on-close>
			<jsonEditor :jsonData="jsonData" v-if="isJson(currentInfo)" />
			<el-card shadow="never" v-else>
				<template #header>
					<div class="card-header">其他数据格式</div>
				</template>
				<span>{{ currentInfo }}</span>
			</el-card>
		</el-dialog>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, watch, defineAsyncComponent, ref } from 'vue';
import dayjs from 'dayjs';
import { useAccessrecordApi } from '/@/api/accessrecord';
import { tokenStateType } from '/@/views/recordsManage/tokenRecords/types';
import commonFunction from '/@/utils/commonFunction';

const { copyText } = commonFunction();

const jsonEditor = defineAsyncComponent(() => import('/@/components/json-editor/index.vue'));
const state = reactive<tokenStateType>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			timetype: '',
			starttime: 0,
			endtime: 0,
			range: [],
			token: '',
		},
	},
});
const dateTypeDict = ref([{ label: '创建时间', value: 1 }]);
const dialog = reactive({
	title: 'JSON',
	isShowDialog: false,
});
const jsonData = ref();
// 判断的是否是JSON字符串
const isJson = (str: any) => {
	if (typeof str == 'string') {
		try {
			var obj = JSON.parse(str);
			// 等于这个条件说明就是JSON字符串 会返回true
			if (typeof obj == 'object' && obj) {
				return true;
			} else {
				//不是就返回false
				return false;
			}
		} catch (e) {
			return false;
		}
	}
	return false;
};
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
// 根据token查询
const handleTokenSearch = () => {
	state.tableData.param.page = 1;
	getTokenTableData();
};
// 重置
const handleSearchReset = () => {
	state.tableData.param.page = 1;
	state.tableData.param.token = '';
	state.tableData.param.timetype = '';
	state.tableData.param.starttime = 0;
	state.tableData.param.endtime = 0;
	state.tableData.param.range = [];
	getTableData();
};
const currentInfo = ref();
const showJson = (row: any, title: string) => {
	currentInfo.value = row;
	if (isJson(row)) {
		jsonData.value = JSON.parse(row);
		dialog.title = title + 'JSON数据';
	}
	dialog.title = title + '数据';
	dialog.isShowDialog = true;
};

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
	useAccessrecordApi()
		.accessrecordTokenList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.record;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
// 根据token获取访问记录
const getTokenTableData = () => {
	state.tableData.loading = true;
	useAccessrecordApi()
		.accessrecordToken({ token: state.tableData.param.token, ...state.tableData.param })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.record;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
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
