<template>
	<div class="lazy-img-container layout-pd">
		<el-card shadow="hover" header="用户优惠包">
			<div class="flex-warp" v-if="state.tableData.data.length > 0">
				<el-row :gutter="15" style="width: 100%">
					<el-col :xs="24" :sm="12" :md="12" :lg="8" :xl="8" class="mb15" v-for="(v, k) in state.tableData.data" :key="k">
						<div class="flex-warp-item">
							<div class="flex-warp-item-box">
								<div class="item-txt">
									<div class="item-txt-title" v-html="v.describe"></div>
									<div class="item-txt-other">
										<el-text type="primary" tag="b" @click="onAllocation(v)">获取</el-text>
										<el-text type="primary" tag="b" class="mr20" @click="onDetail(v)">详情</el-text>
									</div>
								</div>
							</div>
						</div>
					</el-col>
				</el-row>
			</div>
			<el-empty v-else description="暂无数据"></el-empty>
			<template v-if="state.tableData.data.length > 0">
				<el-pagination hide-on-single-page
					style="text-align: right"
					background
					@size-change="onHandleSizeChange"
					@current-change="onHandleCurrentChange"
					:page-sizes="[10, 20, 30]"
					:current-page="state.tableData.param.page"
					:page-size="state.tableData.param.pagesize"
					layout="total, sizes, prev, pager, next, jumper"
					:total="state.tableData.total"
				>
				</el-pagination>
			</template>
		</el-card>
		<!-- 详情 -->
		<DetailDialog ref="detailDialogRef" />
	</div>
</template>

<script setup lang="ts" name="pagesLazyImg">
import { reactive, onMounted, defineAsyncComponent, ref } from 'vue';
import { useCouponApi } from '/@/api/cardRoll/coupon/index';
import other from '/@/utils/other';
import { ElMessage } from 'element-plus';
const DetailDialog = defineAsyncComponent(() => import('./component/detailDialog.vue'));
const useCouponApiCollect = useCouponApi();
const detailDialogRef = ref();
// 定义变量内容
const state = reactive({
	tableData: {
		data: [] as any[],
		total: 99,
		loading: false,
		param: {
			page: 1,
			pagesize: 10,
		},
	},
});
const onDetail = (v: any) => {
	detailDialogRef.value.openDialog(v);
};
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
// 分页点击
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
};
// 分页点击
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
};
const onAllocation = (row: any) => {
	useCouponApiCollect.allocationUser({ discountID: row.id }).then((res: any) => {
		if (res.code === "SUCCESS") {
			ElMessage.success('获取优惠包成功');
		}
	});
};
// 页面加载时
onMounted(() => {
	other.lazyImg('[data-lazy-img-list]', state.tableData.data);
	getTableData();
});
</script>

<style scoped lang="scss">
.lazy-img-container {
	.flex-warp {
		display: flex;
		flex-wrap: wrap;
		align-content: flex-start;
		margin: 0 -5px;
		.flex-warp-item {
			padding: 5px;
			width: 100%;
			height: 360px;
			.flex-warp-item-box {
				border: 1px solid var(--next-border-color-light);
				width: 100%;
				height: 100%;
				border-radius: 2px;
				display: flex;
				flex-direction: column;
				transition: all 0.3s ease;
				flex: 1;
				&:hover {
					cursor: pointer;
					border: 1px solid var(--el-color-primary);
					transition: all 0.3s ease;
					box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.03);
					.item-txt-title {
						color: var(--el-color-primary) !important;
						transition: all 0.3s ease;
					}
					.item-img {
						img {
							transition: all 0.3s ease;
							transform: translateZ(0) scale(1.05);
						}
					}
				}
				.item-img {
					width: 100%;
					height: 215px;
					overflow: hidden;
					img {
						transition: all 0.3s ease;
						width: 100%;
						height: 100%;
					}
				}
				.item-txt {
					flex: 1;
					padding: 15px;
					display: flex;
					flex-direction: column;
					overflow: hidden;
					.item-txt-title {
						text-overflow: ellipsis;
						overflow: hidden;
						-webkit-line-clamp: 2;
						-webkit-box-orient: vertical;
						display: -webkit-box;
						color: #666666;
						transition: all 0.3s ease;
						flex: 1;
						&:hover {
							color: var(--el-color-primary);
							transition: all 0.3s ease;
						}
					}
					.item-txt-other {
						display: flex;
						flex-direction: row-reverse;
						.item-txt-msg {
							font-size: 12px;
							color: #8d8d91;
						}
						.item-txt-price {
							display: flex;
							justify-content: space-between;
							align-items: center;
							.font-price {
								color: #ff5000;
								.font {
									font-size: 22px;
								}
							}
						}
					}
				}
			}
		}
	}
}
</style>
