<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
						<el-button type="success" @click="openDialog('add')">
							<el-icon>
								<ele-FolderAdd />
							</el-icon>
							新增
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="name" label="优惠包名称" show-overflow-tooltip min-width="140" align="left"></el-table-column>
				<el-table-column prop="dayLimit" label="日限额" show-overflow-tooltip align="center" width="100"></el-table-column>
				<el-table-column prop="monthLimit" label="月限额" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column prop="yearLimit" label="年限额" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column prop="limit" label="总限额" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column prop="type" label="优惠包形式" show-overflow-tooltip align="center" min-width="140">
					<template #default="scope">
						<el-tag type="primary" v-if="scope.row.type === 1">赠送额度</el-tag>
						<el-tag type="primary" v-else-if="scope.row.type === 2">赠送优惠券</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="优惠券类型" show-overflow-tooltip align="center" width="100">
					<template #default="scope">
						<el-tag type="primary" v-if="scope.row.quota.type == 1">充值即送</el-tag>
						<el-tag type="primary" v-else-if="scope.row.quota.type == 2">满减优惠</el-tag>
						<el-tag type="primary" v-else>满打折</el-tag>
					</template>
				</el-table-column>

				<el-table-column prop="needVerify" label="需要实名" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.needVerify">是</el-tag>
						<el-tag type="warning" v-else>否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="needCompany" label="需要企业" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.needCompany">是</el-tag>
						<el-tag type="warning" v-else>否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="needUserFace" label="是否人脸" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.needUserFace">是</el-tag>
						<el-tag type="warning" v-else>否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="needCompanyFace" label="是否企业人脸" show-overflow-tooltip width="120" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.needCompanyFace">是</el-tag>
						<el-tag type="warning" v-else>否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="show" label="是否显示" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.show">是</el-tag>
						<el-tag type="warning" v-else>否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="shortDescribe" label="简短描述" show-overflow-tooltip width="170" align="center"></el-table-column>
				<el-table-column label="描述" show-overflow-tooltip width="95" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="showDes(scope.row)">简短描述</el-button>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="140" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="openDialog('allocation', scope.row)">分配</el-button>
						<el-button text type="primary" @click="openDialog('edit', scope.row)">编辑</el-button>
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
			<!-- 新增编辑弹窗 -->
			<ActionDialog ref="actionDialogRef" @refresh="getTableData()" />
			<!-- 分配弹窗 -->
			<AllocationDialog ref="allocationDialogRef" />
			<el-dialog v-model="dialogVisible" title="描述" width="30%">
				<div v-html="currentInfo.describe"></div>
			</el-dialog>
		</el-card>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, defineAsyncComponent } from 'vue';
import { useCouponApi } from '/@/api/cardRoll/discountPackage/index';
import type { stateTypes } from '/@/api/cardRoll/discountPackage/types';
import { ElMessage, ElMessageBox } from 'element-plus';

const ActionDialog = defineAsyncComponent(() => import('./component/actionDialog.vue'));
const AllocationDialog = defineAsyncComponent(() => import('./component/allocationDialog.vue'));
const useCouponApiCollect = useCouponApi();
const actionDialogRef = ref();
const allocationDialogRef = ref();
const state = reactive<stateTypes>({
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
const currentInfo = ref();
const dialogVisible = ref<boolean>(false);
const showDes = (row: any) => {
	currentInfo.value = row;
	dialogVisible.value = true;
};
const openDialog = (type: string, row?: any) => {
	if (type === 'allocation') {
		allocationDialogRef.value.openDialog(type, row);
		return;
	}
	actionDialogRef.value.openDialog(type, row);
};

// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	useCouponApiCollect.couponList(state.tableData.param).then((res: any) => {
		if (res.code === "SUCCESS") {
			state.tableData.data = res.data.discount;
			state.tableData.total = res.data.count;
			state.tableData.loading = false;
		}
	});
};
const onTabelRowDel = (row: any) => {
	ElMessageBox.confirm(`此操作将永久删除站点：${row.name}, 是否继续?`, '提示', {
		confirmButtonText: '删除',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useCouponApiCollect.delCoupon({ id: row.id }).then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('删除成功');
					getTableData();
				}
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
