<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" height="500px">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="110px" :rules="rules">
			<el-row :gutter="20">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="提现金额(元)" prop="cny">
						<template #label
							>提现金额(
							<el-text class="mx-1" type="danger" tag="b">元</el-text>
							)
						</template>
						<el-input-number
							v-model="ruleForm.cny"
							placeholder="请输入提现金额(元)"
							clearable
							style="width: 300px"
							:precision="2"
							:min="20.0"
							:step="0.01"
							step-strictly
							controls-position="right"
						>
						</el-input-number>
					</el-form-item>
				</el-col>
				<template v-if="currentIndex === 1">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="姓名" prop="name">
							<el-input v-model="ruleForm.name" placeholder="请输入姓名" clearable style="width: 300px"> </el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb10">
						<el-form-item>
							<el-text class="mx-1" type="danger">收款微信账号请先关注并绑定创信易公众号。</el-text>
						</el-form-item>
					</el-col>
				</template>
				<template v-if="currentIndex === 2">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="姓名" prop="name">
							<el-input v-model="ruleForm.name" placeholder="请输入姓名" clearable style="width: 300px"> </el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="账号" prop="identity">
							<el-input v-model="ruleForm.identity" placeholder="请输入账号" clearable style="width: 300px"> </el-input>
						</el-form-item>
					</el-col>
				</template>
				<!--            WithdrawWay string `json:"withdrawWay" cs-max:"50"`
            Cny int64 `json:"cny" ci-min:"0"`
            Name string `json:"name" cs-max:"50"`-->
				<template v-if="currentIndex === 3">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="姓名" prop="name">
							<el-input v-model="ruleForm.name" placeholder="请输入姓名" clearable style="width: 300px"> </el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="提现方式" prop="withdrawWay">
							<el-input v-model="ruleForm.withdrawWay" placeholder="请输入提现方式" clearable style="width: 300px"> </el-input>
						</el-form-item>
					</el-col>
				</template>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="支付方式" prop="payType">
						<div class="flex-wrap tab">
							<div
								class="tab-item mt-[10px]"
								v-for="(item, index) in payDict"
								:key="index"
								:class="{ active: currentIndex === item.value }"
								@click="tabCclick(item.value)"
							>
								<img :src="weixin" v-if="item.value === 1" />
								<img :src="apliay" v-else-if="item.value === 2" />
								<img :src="apliay" v-else-if="item.value === 3" />
								{{ item.label }}
							</div>
						</div>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item>
						<el-button type="primary" class="btn" size="large" @click="onSubmit(formRef)">确认</el-button>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
	</el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, toRefs } from 'vue';
import weixin from '/@/assets/recharge/weixin.png';
import apliay from '/@/assets/recharge/apliay.png';
import { FormInstance } from 'element-plus';
import { useWalletRecordApi } from '/@/api/wallet/record';

const emit = defineEmits(['refresh']);
const payDict = [
	{ label: '微信支付', value: 1 },
	{ label: '支付宝', value: 2 },
	{ label: '人工', value: 3 },
];
const currentIndex = ref<number>(1);
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '提现',
	},
});
const formRef = ref();
const tabCclick = (index: number) => {
	currentIndex.value = index;
};
const ruleForm = reactive<any>({
	cny: 0,
	withdrawWay: '',
	name: '',
	identity: '',
	payType: 1,
});
const rules = {
	cny: [{ required: true, message: '请输入提现金额', trigger: 'blur' }],
	withdrawWay: [{ required: true, message: '请输入提现方式', trigger: 'blur' }],
	name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
	identity: [{ required: true, message: '请输入账号', trigger: 'blur' }],
	payType: [{ required: true, message: '请选择支付方式', trigger: 'blur' }],
};
const { dialog } = toRefs(state);
const openDialog = () => {
	reset();
	dialog.value.isShowDialog = true;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	formEl?.validate((valid: boolean) => {
		if (valid) {
			if (currentIndex.value === 1) {
				useWalletRecordApi()
					.wechatWithdraw({
						cny: (ruleForm.cny * 100) | 0,
						name: ruleForm.name,
					})
					.then((res: any) => {
						if (res.code === 'SUCCESS') {
							dialog.value.isShowDialog = false;
							emit('refresh');
						}
					});
			} else if (currentIndex.value === 2) {
				useWalletRecordApi()
					.alipayWithdraw({
						cny: (ruleForm.cny * 100) | 0,
						name: ruleForm.name,
						identity: ruleForm.identity,
					})
					.then((res: any) => {
						if (res.code === 'SUCCESS') {
							dialog.value.isShowDialog = false;
							emit('refresh');
						}
					});
			} else if (currentIndex.value === 3) {
				useWalletRecordApi()
					.selfpayWithdraw({
						cny: (ruleForm.cny * 100) | 0,
						name: ruleForm.name,
						withdrawWay: ruleForm.withdrawWay,
					})
					.then((res: any) => {
						if (res.code === 'SUCCESS') {
							dialog.value.isShowDialog = false;
							emit('refresh');
						}
					});
			}
		} else {
			return false;
		}
	});
};

const reset = () => {
	ruleForm.value = {
		cny: 2000,
		withdrawWay: '',
		name: '',
		identity: '',
		payType: 1,
	};
};

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
