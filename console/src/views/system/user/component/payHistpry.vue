<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="1000px" destroy-on-close @close="closeDialog">
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
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="subject" label="标题" show-overflow-tooltip min-width="200" align="left"></el-table-column>
				<el-table-column prop="supplier" label="供应商" show-overflow-tooltip width="150" align="center"></el-table-column>
				<el-table-column prop="supplier" label="状态" show-overflow-tooltip width="120" align="center">
					<template #default="scope">
						<el-tag type="primary" v-if="scope.row.defrayStatus === 1">订单等待支付</el-tag>
						<el-tag type="success" v-if="scope.row.defrayStatus === 2">订单支付完成</el-tag>
						<el-tag type="danger" v-if="scope.row.defrayStatus === 3">订单退款</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="price" label="价格" show-overflow-tooltip width="80" align="center"></el-table-column>
				<el-table-column prop="quantity" label="数量" show-overflow-tooltip width="80" align="center"></el-table-column>
        <el-table-column prop="unitPrice" label="单价" show-overflow-tooltip width="150" align="center">
          <template #default="scope">
            <el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.unitPrice) }}元</el-text>
          </template>
        </el-table-column>
        <el-table-column prop="realPrice" label="实际支付金额" show-overflow-tooltip width="120" align="center">
          <template #default="scope">
            <el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.realPrice) }}元</el-text>
          </template>
        </el-table-column>
				<el-table-column prop="describe" label="描述" show-overflow-tooltip width="200" align="center"></el-table-column>
				<el-table-column prop="createAt" label="创建时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
        <el-table-column label="消费时间" show-overflow-tooltip width="170" align="center">
          <template #default="scope">
            <span v-if="scope.row.defrayAt">{{ dayjs.unix(scope.row.defrayAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="returnAt" label="退款时间" show-overflow-tooltip width="170" align="center">
          <template #default="scope">
            <span v-if="scope.row.returnAt">{{ dayjs.unix(scope.row.returnAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
				<el-table-column label="操作" width="100" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="danger" @click="onRefund(scope.row)" v-if="scope.row.defrayStatus === 2">退款</el-button>
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
		<checkPhone ref="checkPhoneRef" @checkPhoneSuccess="checkPhoneSuccess" />
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { ElMessage } from 'element-plus';
import checkPhone from '/@/components/checkPhone/index.vue';
import { useDefraytApi } from '/@/api/defray';
import {formatAmount} from "/@/utils/formatAmount";

const info = ref<any>({});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '用户消费记录',
	submitTxt: '',
});
const checkPhoneRef = ref();
const state = reactive({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			timetype: null,
			range: [],
			starttime: null,
			endtime: null,
		},
	},
});
const dateTypeDict = ref([
	{ label: '创建时间', value: 1 },
	{ label: '付款时间', value: 1 },
]);
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
	useDefraytApi()
		.admindefrayList({ ...state.tableData.param, uid: info.value.id })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.defray;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};

const openDialog = async (row?: any) => {
	info.value = JSON.parse(JSON.stringify(row));
	await getTableData();
	dialog.isShowDialog = true;
};
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
	getTableData();
};
const currentInfo = ref();
const onRefund = (row: any) => {
	currentInfo.value = row;
	checkPhoneRef.value.openDialog(row, '退款确认');
};
const checkPhoneSuccess = (phoneToken: string) => {
	useDefraytApi()
		.adminDefrayRefund({ tradeID: currentInfo.value.defrayID, phoneToken: phoneToken })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				ElMessage({
					type: 'success',
					message: '退款成功',
				});
				closeDialog();
			}
		});
};
// 暴露变量
defineExpose({
	openDialog,
});
</script>
