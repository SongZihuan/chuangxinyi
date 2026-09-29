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
						<el-button type="primary" @click="handleSearch">
							<el-icon>
								<ele-Search />
							</el-icon>
							查询
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="id" label="ID" show-overflow-tooltip width="100" align="left"></el-table-column>
        <el-table-column prop="userID" label="用户ID" show-overflow-tooltip width="100" align="left"></el-table-column>
				<el-table-column prop="name" label="名称" show-overflow-tooltip align="center" min-width="100"></el-table-column>
        <el-table-column prop="type" label="优惠形式" show-overflow-tooltip align="center" min-width="140">
          <template #default="scope">
            <el-tag type="primary" v-if="scope.row.type === 1">充值立赠</el-tag>
            <el-tag type="primary" v-else-if="scope.row.type === 2">消费满减</el-tag>
            <el-tag type="primary" v-else-if="scope.row.type === 3">消费打折</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="优惠内容" show-overflow-tooltip align="center" min-width="200">
          <template #default="scope">
						<span v-if="scope.row.type == 1">
							<el-text tag="b">充值满</el-text> <span class="mr5"> ￥{{ formatAmount(scope.row.content.bottom) }} </span> <el-text tag="b">立赠</el-text>
							<span class="mr5"> ￥{{ formatAmount(scope.row.content.send) }} </span></span
            >
            <span v-if="scope.row.type == 2">
							<el-text tag="b">消费满</el-text> <span class="mr5"> ￥{{ formatAmount(scope.row.content.bottom) }} </span> <el-text tag="b">立减</el-text>
							<span class="mr5"> ￥{{ scope.row.content.discount }} </span></span
            >
            <span v-if="scope.row.type == 3">
							<el-text tag="b">消费满</el-text> <span class="mr5"> ￥{{ formatAmount(scope.row.content.bottom) }} </span> <el-text tag="b">立即打折</el-text>
							<span class="mr5"> {{ scope.row.content.pre }}% </span></span
            >
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
			<!-- 分配弹窗 -->
			<el-dialog v-model="dialogVisible" title="简短描述" width="30%">
				<div v-html="currentInfo.shortDescribe"></div>
			</el-dialog>
		</el-card>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, watch } from 'vue';
import { useCouponApi } from '/@/api/cardRoll/discountPackage/index';
import dayjs from 'dayjs';
import {formatAmount} from "/@/utils/formatAmount";
const dateTypeDict = ref([{ label: '创建时间', value: 1 }]);
const state = reactive({
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
const currentInfo = ref();
const dialogVisible = ref<boolean>(false);

//搜索
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};

// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	useCouponApi()
		.getCouponList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.coupons;
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
