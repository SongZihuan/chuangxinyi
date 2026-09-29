<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" height="500px">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="110px" :rules="rules">
			<el-row :gutter="20">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="充值金额(元)" prop="cny">
						<template #label>充值金额( <el-text class="mx-1" type="danger" tag="b">元</el-text>)</template>
						<el-input v-model="ruleForm.cny" placeholder="请输入充值金额(元)" style="width: 300px" disabled> </el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="优惠券" prop="couponsID">
						<el-select v-model="ruleForm.couponsID" placeholder="请选择优惠券" style="width: 300px">
							<el-option v-for="(item, index) in couponsIDDict" :value="item.id" :key="index">
								<template #default>
									<span>{{ item.name }}（充值 ￥{{ formatAmount(item.content.bottom) }} 立得 ￥{{ formatAmount(item.content.send) }}） </span>
								</template>
							</el-option>
						</el-select>
					</el-form-item>
				</el-col>

				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="支付方式" prop="payType">
						<div class="tab">
							<div
								class="tab-item"
								v-for="(item, index) in payDict"
								:key="index"
								:class="{ active: currentIndex === index }"
								@click="tabCclick(index)"
							>
								<img :src="weixin" v-if="item.value === 1" />
								<img :src="apliay" v-else />
								{{ item.label }}
							</div>
						</div>
					</el-form-item>
				</el-col>

				<el-form-item v-if="activeName == 'first'">
					<div style="width: 300px"><silenceSlider ref="sliderRef" @nvcValEmit="nvcValEmit" @siderEmit="siderEmit" /></div>
				</el-form-item>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item v-if="authUser.aliPay"><el-button type="primary" class="btn" @click="submit(formRef)">确认</el-button> </el-form-item>
					<el-form-item v-else><el-button type="primary" class="btn" disabled>抱歉，暂无权限</el-button> </el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<AlipayDialog ref="alpayDialogRef" @queryOrderBack="queryOrderBack" />
		<WeixinDialog ref="weixinDialogRef" @queryOrderBack="queryOrderBack" />
	</el-dialog>
</template>

<script setup lang="ts">
import { defineAsyncComponent, ref, reactive, toRefs, watch, onUnmounted } from 'vue';
import { type FormRules, ElMessage, type FormInstance } from 'element-plus';
import silenceSlider from '/@/components/Slider/silence.vue';
import { useAlipayApi } from '/@/api/alipay/index';
import { useWeixinApi } from '/@/api/weixin/index';
import { useBaseApi } from '/@/api/base/index';
import { useCouponApi } from '/@/api/cardRoll/coupon/index';
import { type formTypes } from '/@/api/recharge/types';
import weixin from '/@/assets/recharge/weixin.png';
import apliay from '/@/assets/recharge/apliay.png';
import commonFunction from '/@/utils/commonFunction';
import useSubAuth from '/@/hooks/useSubAuth';
import { formatAmount } from '/@/utils/formatAmount';
const emit = defineEmits(['refresh', 'paySuccess']);
import { isWeiXin } from '/@/utils/weixin';
const AlipayDialog = defineAsyncComponent(() => import('./component/alipayDialog.vue'));
const WeixinDialog = defineAsyncComponent(() => import('./component/weixinDialog.vue'));
const alpayDialogRef = ref();
const weixinDialogRef = ref();
const tradeID = ref<string>('');
const activeName = ref('first');
const useAlipayApiCollect = useAlipayApi();
const useWeixinApiCollect = useWeixinApi();
const useBaseApiCollect = useBaseApi();
const formRef = ref();
const { isMobile } = commonFunction();
const payDict = [
	{ label: '微信支付', value: 1 },
	{ label: '支付宝', value: 2 },
];
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '在线支付',
	},
});
const { dialog } = toRefs(state);
const openDialog = (cny?: number) => {
	reset();
	if (cny) {
		ruleForm.value.cny = formatAmount(cny);
		newForm.value.cny = formatAmount(cny);
	}
	getCouponList();
	dialog.value.isShowDialog = true;
};
const currentIndex = ref<number>(0);
let ruleForm = ref<formTypes>({
	cny: 1.0,
	payType: 1,
	couponsID: null,
	payMode: 1,
});
let newForm = ref({
	cny: 1,
	payWay: '',
	couponsID: null,
});
const currentVal = ref();

const authUser = useSubAuth();
const couponsIDDict = ref<any>([]);
const rules = reactive<FormRules>({
	cny: [{ required: true, message: '请输入充值金额', trigger: 'blur' }],
	payType: [{ required: true, message: '请选择支付类型', trigger: 'blur' }],
});
const newRules = reactive<FormRules>({
	cny: [{ required: true, message: '请输入充值金额', trigger: 'blur' }],
	payWay: [{ required: true, message: '请输入支付方式', trigger: 'blur' }],
});

const sliderRef = ref();
const newSliderRef = ref();
const isSliderCheck = ref<boolean>(false);
const isSecond = ref(false);
const siderEmit = (val: string) => {
	currentVal.value = val;
	isSecond.value = true;
	isSliderCheck.value = true;
};

const nvcValEmit = (val: string) => {
	currentVal.value = val;
	isSliderCheck.value = true;
};
const tabCclick = (index: number) => {
	currentIndex.value = index;
};
const reset = () => {
	activeName.value == 'second' ? newSliderRef.value && newSliderRef.value.resetSider() : sliderRef.value && sliderRef.value.resetSider();
	ruleForm.value = {
		cny: 1.0,
		payType: 1,
		couponsID: null,
		payMode: 1,
	};
	newForm.value = {
		cny: 1,
		payWay: '',
		couponsID: null,
	};
	isSliderCheck.value = false;
};
const submit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	if (!isSliderCheck.value && isSecond.value) {
		ElMessage.warning('请滑动滑块移动最右侧!');
		return;
	}
	formEl.validate(async (valid) => {
		if (valid) {
			sliderRef.value.registerClick();
			setTimeout(() => {
				if (isMobile()) {
					if (isWeiXin() && currentIndex.value == 0) {
						wxPay();
						return;
					}
					currentIndex.value ? mobileAlipay() : mobileWeixinPay();
				} else {
					currentIndex.value ? pcAlipay() : pcWeixinPay();
				}
			}, 500);
		}
	});
};
const wxPay = () => {
	let miaddleData = { ...ruleForm.value, cny: (ruleForm.value.cny * 100) | 0, h5Type: 'Wap' };
	useWeixinApiCollect.weixinPayJsapi(miaddleData).then((res: any) => {
		if (res.code === 'SUCCESS') {
			tradeID.value = res.data.id;
			if (typeof WeixinJSBridge == 'undefined') {
				if (document.addEventListener) {
					document.addEventListener('WeixinJSBridgeReady', onBridgeReady, false);
				} else if (document.attachEvent) {
					document.attachEvent('WeixinJSBridgeReady', onBridgeReady);
					document.attachEvent('onWeixinJSBridgeReady', onBridgeReady);
				}
			} else {
				onBridgeReady(res);
			}
		}

		// 调用微信支付
	});
};
const onBridgeReady = (res: any) => {
	WeixinJSBridge.invoke(
		'getBrandWCPayRequest',
		{
			appId: res.data.appId, //公众号ID，由商户传入
			timeStamp: res.data.timeStamp, //时间戳，自1970年以来的秒数
			nonceStr: res.data.nonceStr, //随机串
			package: res.data.package,
			signType: res.data.signType, //微信签名方式：
			paySign: res.data.paySign,
		},
		function (res) {
			if (res.err_msg == 'get_brand_wcpay_request:ok') {
				// 使用以上方式判断前端返回,微信团队郑重提示：
				//res.err_msg将在用户支付成功后返回ok，但并不保证它绝对可靠。
				emit('refresh');
			}
		}
	);
};
const pcAlipay = () => {
	const domain = window.location.origin;
	let miaddleData = { ...ruleForm.value, cny: (ruleForm.value.cny * 100) | 0, nvc: currentVal.value, returnUrl: `${domain}/paysuccess` };
	useAlipayApiCollect.alipay(miaddleData).then((res: any) => {
		if (res.code === 'LOGIC_DENY' && res.subCode === 'CAPTCHA_SECOND_CHECK') {
			sliderRef.value.siderCheck();
			return;
		}
		if (res.code === 'SUCCESS') {
			tradeID.value = res.data.id;
			alpayDialogRef.value.openDialog(res.data.url);
		}
	});
};
const getCouponList = () => {
	couponsIDDict.value = [];
	useCouponApi()
		.couponDict()
		.then((res: any) => {
			if (res.code === 'SUCCESS' && res.data.coupons.length > 0) {
				res.data.coupons.forEach((element: any) => {
					if (element.type === 1) {
						couponsIDDict.value.push(element);
					}
				});
			}
		});
};
const mobileAlipay = () => {
	let miaddleData = {
		...ruleForm.value,
		cny: (ruleForm.value.cny * 100) | 0,
		nvc: currentVal.value,
		returnUrl: `${window.location.href}&paysucess=1`,
		quiteUrl: `${window.location.href}`,
	};

	useAlipayApiCollect.wapAlipay(miaddleData).then((res: any) => {
		if (res.code === 'LOGIC_DENY' && res.subCode === 'CAPTCHA_SECOND_CHECK') {
			sliderRef.value.siderCheck();
			return;
		}
		if (res.code === 'SUCCESS') {
			tradeID.value = res.data.id;
			window.location.href = res.data.payUrl;
		}
	});
};

const pcWeixinPay = () => {
	let miaddleData = { ...ruleForm.value, cny: (ruleForm.value.cny * 100) | 0, nvc: currentVal.value };
	useWeixinApiCollect.weixinPay(miaddleData).then((res: any) => {
		if (res.code === 'LOGIC_DENY' && res.subCode === 'CAPTCHA_SECOND_CHECK') {
			sliderRef.value.siderCheck();
			return;
		}
		if (res.code === 'SUCCESS') {
			tradeID.value = res.data.id;
			weixinDialogRef.value.openDialog(res.data.url);
		}
	});
};
const mobileWeixinPay = () => {
	let miaddleData = {
		...ruleForm.value,
		cny: (ruleForm.value.cny * 100) | 0,
		nvc: currentVal.value,
		h5Type: 'Wap',
		returnUrl: `${window.location.href}&paysucess=1`,
		quiteUrl: `${window.location.href}`,
	};
	useWeixinApiCollect.weixinPayMobie(miaddleData).then((res: any) => {
		if (res.code === 'LOGIC_DENY' && res.subCode === 'CAPTCHA_SECOND_CHECK') {
			sliderRef.value.siderCheck();
			return;
		}
		if (res.code === 'SUCCESS') {
			tradeID.value = res.data.id;
			const url= window.location.href +'&paysucess=1'
			window.location.href = res.data.h5Url +'&redirect_url='+encodeURIComponent(url);
		}
	});
};
const queryOrderBack = async () => {
	if (tradeID.value) {
		await useBaseApiCollect.checkOrder({ tradeid: tradeID.value }).then((res: any) => {
			if (res.code === 'SUCCESS' && res.data.success) {
				clearInterval(timer.value);
				alpayDialogRef.value.closeDialog();
				weixinDialogRef.value.closeDialog();
				emit('paySuccess');
			} else {
				// ElMessage.warning('您暂未支付成功 !');
			}
		});
	} else {
		ElMessage.warning('请先支付 !');
	}
};
const timer = ref();
watch(
	() => tradeID.value,
	async () => {
		if (tradeID.value) {
			timer.value = setInterval(async () => {
				await queryOrderBack();
			}, 1000);
		}
	}
);
onUnmounted(() => {
	clearInterval(timer.value);
});
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss">
.price {
	margin-top: 10px;
	margin-bottom: 10px;
	.price-item {
		font-weight: 600;
		border-radius: 100px;
		color: #432b0c;
		font-size: 16px;
		text-align: center;
		line-height: 40px;
		box-sizing: border-box;
		color: #432b0c;
		.price-label {
			padding: 10px;
			margin: 10px 0px;
			box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.1);
			background-color: #fff;
			border-radius: 8px;
		}
		.active {
			background-color: #f9f8df;
			border: 3px solid #e8d0a7;
		}
	}
}
.pay {
	width: 100%;
	margin: auto;
	padding: 10px 0;
	text-align: center;
	background: linear-gradient(180deg, #fae3bf 0%, #f5c379 100%);
	border-radius: 100px;
	color: #432b0c;
	font-weight: 700;
	font-size: 16px;
	margin: 20px 0px;
}
.tab {
	display: flex;
	flex-direction: row;

	.tab-item {
		width: 130px;
		height: 50px;
		background: #fafafb;
		margin-right: 10px;
		display: flex;
		flex-direction: row;
		align-items: center;
		color: #414960;
		font-weight: 600;
		border-radius: 6px;
		position: relative;
		img {
			width: 32px;
			height: 32px;
			margin-left: 15px;
			margin-right: 16px;
		}
	}
	.active {
		border: 2px solid var(--el-color-primary);
	}
}
//.btn {
//	width: 80px;
//}
//.btn-primary {
//  width: ;
//}
.go {
	text-decoration: underline;
}
.action {
	display: flex;
	flex-direction: row;
	.agreement {
		color: var(--el-color-primary);
		text-decoration: underline;
		margin-left: 2px;
	}
}
</style>
