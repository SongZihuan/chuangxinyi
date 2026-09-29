<template>
	<div>
		<div class="container layout-padding warp" v-if="!paySuccess">
			<el-card class="box-card">
				<div class="title">订单信息</div>
				<div class="purchase-box">
					<el-table :data="tableData" border stripe style="width: 100%">
						<el-table-column prop="subject" label="标题" min-width="180" align="center"></el-table-column>
						<el-table-column prop="quantity" label="数量" width="180" align="center"></el-table-column>
						<el-table-column prop="unitPrice" label="单价" width="180" align="center">
							<template #default="scope">
								<span>￥{{ formatAmount(scope.row.unitPrice) }}</span>
							</template>
						</el-table-column>
						<el-table-column prop="describe" label="描述" width="180" align="center"></el-table-column>
						<el-table-column prop="supplier" label="供应商" width="180" align="center"></el-table-column>
						<el-table-column prop="price" label="价格" width="180" align="center">
							<template #default="scope">
								<span class="table-money">￥{{ formatAmount(scope.row.price) }}</span>
							</template>
						</el-table-column>
						<el-table-column prop="owner" label="购买人信息" width="180" align="center">
							<template #default="scope">
								<el-button v-if="scope.row.hasOwner" type="text" @click="openOwnerDialog(scope.row.owner)">查看</el-button>
								<span v-else> 无购买人 </span>
							</template>
						</el-table-column>
						<el-table-column prop="mustSelfDefray" label="是否必须购买人支付" width="180" align="center">
							<template #default="scope">
								<span v-if="scope.row.mustSelfDefray">是</span>
								<span v-else>否</span>
							</template>
						</el-table-column>
					</el-table>
				</div>
			</el-card>
			<el-card class="mt10">
				<div class="card-block">
					<div class="w-[80px]">优惠券</div>
					<div class="card-block-content">
						<el-select v-model="ruleForm.couponsID" placeholder="请选择优惠券" style="w-full">
							<el-option v-for="(item, index) in couponsIDDict" :value="item.id" :label="item.name" :key="index" />
						</el-select>
					</div>
				</div>
			</el-card>
			<div class="payment">
				<div class="payment-box">
					<div class="left">
						应付款:<span class="money">￥{{ changeMoney(purchaseData.price) || 0 }}</span>
					</div>
					<el-button type="primary" class="right" @click="onPay">立即支付</el-button>
					<el-button text type="primary" :icon="Share" @click="onShare">分享好友支付</el-button>
				</div>
			</div>
			<purchase-info ref="purchaseInfoRef" />
			<doubleCheck ref="doubleCheckRef" @checkPhoneSuccess="doubleCheckSuccess" />
		</div>
		<div v-else>
			<div class="container">
				<div class="main">
					<SvgIcon name="iconfont icon-zhifuchenggong" :size="100" class="icon left-item-animation"> </SvgIcon>
					<div class="tip left-item-animation">支付成功</div>
					<div class="left-item-animation pay-text">您已完成支付</div>
					<div v-if="redirectUrl" class="left-item-animation left-item-btn mt10">
						<el-button type="primary" size="default" round @click="onGoHome">返回商户页面</el-button>
					</div>
				</div>
			</div>
		</div>
		<rechargeDialog ref="rechargeDialogRef" @refresh="getUserInfo" @paySuccess="onPay()" v-if="!paySuccess" />
	</div>
</template>
<script setup lang="ts">
// 获取订单列表
import { purchaseApi } from '/@/api/purchase';
import { onMounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useCouponApi } from '/@/api/cardRoll/coupon';
import { formTypes } from '/@/api/recharge/types';
import PurchaseInfo from './components/purchaseInfo.vue';
import { useUserApi } from '/@/api/user/user';
import { Session } from '/@/utils/storage';
import doubleCheck from '/@/components/doubleCheck/index.vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { formatAmount } from '/@/utils/formatAmount';
import { Share } from '@element-plus/icons-vue';
import commonFunction from '/@/utils/commonFunction';
import { useUserInfo } from '/@/stores/userInfo';
import rechargeDialog from '/@/views/purchase/components/recharge/index.vue';
const router = useRouter();
const route = useRoute();
const token = ref<string>('');
const tradeID = ref<string>('');
const tableData = ref<any[]>([]);
const purchaseInfoRef = ref<any>({});
const couponsIDDict = ref<any>([]);
const rechargeDialogRef = ref();
const paySuccess = ref(false);
const redirectUrl = ref('');
tableData.value = [];
const stores = useUserInfo();
let ruleForm = ref<formTypes>({
	cny: 1.0,
	payType: 1,
	couponsID: null,
	payMode: 1,
});
// 将分转换为元
const changeMoney = (money: number) => {
	return (money / 100).toFixed(2);
};
const onGoHome = () => {
	window.open(redirectUrl.value, '_self');
};
const purchaseData = ref<any>({
	subject: '',
	quantity: 0,
	unitPrice: 0,
	describe: '',
	supplier: '',
	price: 0,
	hasOwner: false,
	owner: '',
	mustSelfDefray: false,
});

const getCouponList = async () => {
	couponsIDDict.value = [];
	await useCouponApi()
		.couponDict()
		.then((res) => {
			if (res.code === 'SUCCESS' && res.data.coupons.length > 0) {
				res.data.coupons.forEach((element: any) => {
					if (element.type === 2 || element.type === 3) {
						couponsIDDict.value.push(element);
					}
				});
			}
		});
};
const useUserApiCollect = useUserApi();
const getUserInfo = async () => {
	await useUserApiCollect.userInfo().then((res: any) => {
		if (res.code === 'SUCCESS') {
			stores.setUserData(res.data);
			Session.set('userInfo', res.data);
			Session.set('userData', res.data);
			stores.setUserInfos();
		}
	});
};
const onPay = () => {
	if (!token.value) {
		ElMessage.error('token不存在');
		return;
	}
	let data = {
		couponsID: ruleForm.value.couponsID,
		token: token.value,
	};
	purchaseApi()
		.purchase(data)
		.then((res) => {
			if (res.code === 'SUCCESS') {
				ElMessage.success('购买成功');
				paySuccess.value = true;
			} else if (res.code === 'NOT_TOKEN_DENY') {
				router.push(`/login?redirect=/purchase&params=${JSON.stringify(route.query)}`);
			} else if (res.code === 'LOGIC_DENY' && res.subCode === 'DEFRAY_INSUFFICIENT') {
				rechargeDialogRef.value.openDialog(res.data.cny);
			} else if (res.code === 'LOGIC_DENY' && res.subCode === 'PAY_MUST_VERIFY') {
				ElMessageBox.confirm('您暂未实名认证，是否前往实名认证？', 'Warning', {
					confirmButtonText: '前往',
					cancelButtonText: '取消',
					type: 'warning',
				})
					.then(() => {
						router.push('/home');
					})
					.catch(() => {});
			}else{
				// router.push('/home');
			}
		});
};
const { copyText } = commonFunction();
const onShare = () => {
	// 获取当前域名，带http或者https的
	const domain = window.location.origin;
	copyText(domain + '/purchase?tradeID=' + tradeID.value + '&token=' + token.value);
};

const getOrders = async () => {
	purchaseApi()
		.getOrderInfo({ token: token.value })
		.then((res) => {
			if (res.code === 'SUCCESS') {
				tableData.value = [res.data.info];
				purchaseData.value = res.data.info;
			}
		});
};
const openOwnerDialog = (row: any) => {
	purchaseInfoRef.value.openDialog(row);
};
onMounted(async () => {
	// 获取URL参数
	token.value = router.currentRoute.value.query.token as string;
	tradeID.value = router.currentRoute.value.query.tradeID as string;
	redirectUrl.value = router.currentRoute.value.query.redirectUrl as string;

	await getOrders();
	await getCouponList();
	if (router.currentRoute.value.query.paysucess && router.currentRoute.value.query.paysucess =='1') {
		onPay();
	}
	
});
</script>
<style scoped lang="scss">
.main {
	width: 100%;
	height: calc(100vh + 40px);
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	background: var(--next-bg-color);
	margin-top: -40px;
	overflow: hidden;
	.icon {
		color: var(--el-color-primary);
	}
	.tip {
		color: var(--el-color-primary);
		font-size: 20px;
		font-weight: 600;
		margin-top: -15px;
	}
	.pay-text {
		color: var(--el-text-color-primary);
		margin-top: 4px;
	}
	.left-item-animation {
		opacity: 0;
		animation-name: error-num;
		animation-duration: 0.5s;
		animation-fill-mode: forwards;
	}
	.left-item-btn {
		animation-delay: 0.2s;
	}
}
.el-table th.el-table__cell {
	background-color: #f5f7fa;
	color: #909399;
	font-weight: 400;
	font-size: 12px;
	padding: 10px 0;
}
.box-card {
	width: 100%;
	.title {
		font-size: 18px;
		color: #373d41;
		border-bottom: 1px solid #ebeef5;
		font-weight: 600;
	}
	.purchase-box {
		padding: 20px 0;
	}
}
.card-block {
	display: flex;
	justify-content: flex-start;
	align-items: center;
	.card-block-title {
		width: 120px;
		text-align: left;
		font-size: 14px;
		color: #373d41;
		font-weight: 600;
		padding-right: 40px;
	}
	.card-block-content {
		width: 100%;
		display: flex;
		justify-content: flex-start;
		align-items: center;
		.el-select {
			width: 240px;
			.el-input {
				width: 100%;
				.el-input__inner {
					width: 100%;
				}
			}
		}
	}
}
.payment {
	position: fixed;
	bottom: 0;
	left: 0;
	width: 100%;
	height: 100px;
	background-color: #fff;
	box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.1);
	.payment-box {
		height: 100%;
		display: flex;
		justify-content: flex-end;
		align-items: center;
		padding: 0 20px;
		.left {
			font-size: 12px;
			color: #373d41;
			padding: 0 30px;
			.money {
				font-size: 24px;
				color: #ff8a00;
			}
		}
		.right {
			height: 40px;
			line-height: 40px;
			text-align: center;
			border-radius: 0;
			padding: 0 16px;
			color: #fff;
			cursor: pointer;
		}
	}
}
.table-money {
	color: #ff8a00;
}
</style>
