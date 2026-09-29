<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="800px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="120px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="优惠包名称" prop="name">
						<el-input v-model="ruleForm.name" placeholder="请输入优惠包名称" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="类型" prop="type">
						<el-select v-model="ruleForm.type" placeholder="请选择类型" style="width: 100%">
							<el-option v-for="(item, index) in typeDict" :value="item.value" :label="item.label" :key="index" />
						</el-select>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" v-if="ruleForm.type === 1">
					<el-form-item label="赠送额度（元）" required label-width="150px">
						<template #label>
							<el-tooltip effect="dark" content="直接赠送的额度" placement="top-start">
								<el-icon style="height: 100%"><ele-QuestionFilled /></el-icon>
							</el-tooltip>
							<span>赠送额度（元）</span>
						</template>
						<el-input-number
							v-model="ruleForm.quota.amount"
							placeholder="请输入赠送额度"
							controls-position="right"
							:min="0.01"
							:step="0.01"
							step-strictly
							style="width: 100%"
						></el-input-number>
					</el-form-item>
				</el-col>
				<template v-else>
					<!--  充值送额度-->
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="优惠券类型" required>
							<el-select v-model="ruleForm.quota.type" placeholder="请选择类型" style="width: 100%">
								<el-option v-for="(item, index) in quotaTypeDict" :value="item.value" :label="item.label" :key="index" />
							</el-select>
						</el-form-item>
					</el-col>

					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20" v-if="ruleForm.quota.type == 1">
						<el-form-item label="充值最低限额（元）" required  label-width="200px">
							<template #label>
								<el-tooltip effect="dark" content="充值最低限额" placement="top-start">
									<el-icon style="height: 100%"><ele-QuestionFilled /></el-icon>
								</el-tooltip>
								<span>充值最低限额（元）</span>
							</template>
							<el-input-number
								v-model="ruleForm.quota.bottom"
								placeholder="请输入充值最低限额"
								controls-position="right"
								:min="0"
								:step="0.01"
								step-strictly
								style="width: 100%"
							></el-input-number>
						</el-form-item>
					</el-col>

					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" v-if="ruleForm.quota.type == 1">
						<el-form-item label="赠送额度（元）"  label-width="150px">
							<template #label>
								<el-tooltip effect="dark" content="赠送账户余额" placement="top-start">
									<el-icon style="height: 100%"><ele-QuestionFilled /></el-icon>
								</el-tooltip>
								<span>赠送额度（元）</span>
							</template>
							<el-input-number
								v-model="ruleForm.quota.send"
								placeholder="请输入赠送额度"
								controls-position="right"
								:min="0.01"
								:step="0.01"
								step-strictly
								style="width: 100%"
							></el-input-number>
						</el-form-item>
					</el-col>
					<!-- 购物满减 -->
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20" v-if="ruleForm.quota.type == 2">
						<el-form-item label="购物满减（元）" required  label-width="150px">
							<template #label>
								<el-tooltip effect="dark" content="购物达到一定的金额优惠一定量的额度" placement="top-start">
									<el-icon style="height: 100%"><ele-QuestionFilled /></el-icon>
								</el-tooltip>
								<span>满减线（元）</span>
							</template>
							<el-input-number
								v-model="ruleForm.quota.bottom"
								placeholder="请输入购物满减"
								controls-position="right"
								:min="0.01"
								:step="0.01"
								step-strictly
								style="width: 100%"
							></el-input-number>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" v-if="ruleForm.quota.type == 2">
						<el-form-item label="满减额度（元）"  label-width="150px">
							<template #label>
								<el-tooltip effect="dark" content="满减额度" placement="top-start">
									<el-icon style="height: 100%"><ele-QuestionFilled /></el-icon>
								</el-tooltip>
								<span>满减额度（元）</span>
							</template>
							<el-input-number
								v-model="ruleForm.quota.discount"
								placeholder="请输入满减额度"
								controls-position="right"
								:min="0.01"
								:step="0.01"
								step-strictly
								style="width: 100%"
							></el-input-number>
						</el-form-item>
					</el-col>
					<!-- 购物打折 -->
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20" v-if="ruleForm.quota.type == 3">
						<el-form-item label="打折线（元）" required  label-width="150px">
							<template #label>
								<el-tooltip effect="dark" content="打折线" placement="top-start">
									<el-icon style="height: 100%"><ele-QuestionFilled /></el-icon>
								</el-tooltip>
								<span>打折线（元）</span>
							</template>
							<el-input-number
								v-model="ruleForm.quota.bottom"
								placeholder="请输入打折线"
								controls-position="right"
								:min="0.01"
								:step="0.01"
								step-strictly
								style="width: 100%"
							></el-input-number>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" v-if="ruleForm.quota.type == 3">
						<el-form-item label="打折百分比">
							<template #label>
								<el-tooltip effect="dark" content="打折百分比" placement="top-start">
									<el-icon style="height: 100%"><ele-QuestionFilled /></el-icon>
								</el-tooltip>
								<span>打折百分比</span>
							</template>
							<el-input-number
								v-model="ruleForm.quota.pre"
								placeholder="请输入打折百分比"
								controls-position="right"
								:min="0"
								:max="100"
								:step="1"
								step-strictly
								style="width: 100%"
							></el-input-number>
						</el-form-item>
					</el-col>
				</template>

				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="日限额" prop="dayLimit">
						<el-input-number
							v-model="ruleForm.dayLimit"
							placeholder="请输入日限额"
							controls-position="right"
							:min="0"
							:step="1"
							step-strictly
							style="width: 100%"
						></el-input-number>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="月限额" prop="monthLimit">
						<el-input-number
							v-model="ruleForm.monthLimit"
							placeholder="请输入月限额"
							controls-position="right"
							:min="0"
							:step="1"
							step-strictly
							style="width: 100%"
						>
						</el-input-number>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="年限额" prop="yearLimit">
						<el-input-number
							v-model="ruleForm.yearLimit"
							placeholder="请输入年限额"
							controls-position="right"
							:min="0"
							:step="1"
							step-strictly
							style="width: 100%"
						></el-input-number>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="总限额" prop="limit">
						<el-input-number
							v-model="ruleForm.limit"
							placeholder="请输入总限额"
							controls-position="right"
							:min="0"
							:step="1"
							step-strictly
							style="width: 100%"
						></el-input-number>
					</el-form-item>
				</el-col>
				<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="实名">
						<el-radio-group v-model="ruleForm.needVerify">
							<el-radio :label="true">是</el-radio>
							<el-radio :label="false">否</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="仅限企业">
						<el-radio-group v-model="ruleForm.needCompany">
							<el-radio :label="true">是</el-radio>
							<el-radio :label="false">否</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="身份证实名">
						<el-radio-group v-model="ruleForm.needUserOrigin">
							<el-radio :label="true">是</el-radio>
							<el-radio :label="false">否</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="企业证件实名">
						<el-radio-group v-model="ruleForm.needCompanyOrigin">
							<el-radio :label="true">是</el-radio>
							<el-radio :label="false">否</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="使用人人脸核验">
						<el-radio-group v-model="ruleForm.needUserFace">
							<el-radio :label="true">是</el-radio>
							<el-radio :label="false">否</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="法人人脸核验">
						<el-radio-group v-model="ruleForm.needCompanyFace">
							<el-radio :label="true">是</el-radio>
							<el-radio :label="false">否</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="是否显示">
						<el-radio-group v-model="ruleForm.show">
							<el-radio :label="true">是</el-radio>
							<el-radio :label="false">否</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="简短描述" prop="shortDescribe">
						<el-input v-model="ruleForm.shortDescribe" type="textarea" placeholder="请输入简短描述"></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="描述" prop="describe"> <WangEditor v-model:get-html="ruleForm.describe" /></el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue';
import WangEditor from '/@/components/Editor/index.vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { useCouponApi } from '/@/api/cardRoll/discountPackage';
import type { formTypes } from '/@/api/cardRoll/discountPackage/types';
const useCouponApiCollect = useCouponApi();
const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<formTypes>({
	name: '',
	describe: '',
	shortDescribe: '',
	type: 1,
	quota: {
		type: 1,
		bottom: 0,
		send: 0,
		discount: 0,
		pre: 0,
		amount: 0,
	},
	dayLimit: 0,
	monthLimit: 0,
	yearLimit: 0,
	limit: 0,
	needVerify: false,
	needCompany: false,
	needUserOrigin: false,
	needCompanyOrigin: false,
	needUserFace: false,
	needCompanyFace: false,
	show: false,
});
const typeDict = [
	{ label: '赠送额度', value: 1 },
	{ label: '赠送优惠券', value: 2 },
] as dictTypes[];
const quotaTypeDict = [
	{ label: '充值送额度', value: 1 },
	{ label: '购物满减', value: 2 },
	{ label: '购物打折', value: 3 },
];
const rules = reactive<FormRules>({
	name: [{ required: true, message: '请输入优惠包名称', trigger: 'blur' }],
	type: [{ required: true, message: '请选择类型', trigger: 'blur' }],
	bottom: [{ required: true, message: '请输入满足条件', trigger: 'blur' }],
	dayLimit: [{ required: true, message: '请输入日限额', trigger: 'blur' }],
	monthLimit: [{ required: true, message: '请输入月限额', trigger: 'blur' }],
	yearLimit: [{ required: true, message: '请选择年限额', trigger: 'blur' }],
	limit: [{ required: true, message: '请选择总限额', trigger: 'blur' }],
	shortDescribe: [{ required: true, message: '请输入描述', trigger: 'blur' }],
	describe: [{ required: true, message: '请输入简短描述', trigger: 'blur' }],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '优惠包新增',
	submitTxt: '新增',
});

const openDialog = (type: string, row: any) => {
  reset();
	if (type === 'add') {
		dialog.title = '优惠包新增';
		dialog.submitTxt = '新增';
	} else {
		ruleForm.value = JSON.parse(JSON.stringify(row));

    ruleForm.value.quota.bottom = ruleForm.value.quota.bottom / 100
    ruleForm.value.quota.send = ruleForm.value.quota.send / 100
    ruleForm.value.quota.discount = ruleForm.value.quota.discount / 100
    ruleForm.value.quota.amount = ruleForm.value.quota.amount / 100


		dialog.title = '优惠包编辑';
		dialog.submitTxt = '编辑';
	}
	dialog.type = type;
	dialog.isShowDialog = true;
};

//重置
const reset = () => {
	ruleForm.value = {
		name: '',
		describe: '',
		shortDescribe: '',
		type: 1,
		quota: {
			type: 1,
			bottom: 0,
			send: 0,
			discount: 0,
			pre: 0,
		},
		dayLimit: 0,
		monthLimit: 0,
		yearLimit: 0,
		limit: 0,
		needVerify: false,
		needCompany: false,
		needUserOrigin: false,
		needCompanyOrigin: false,
		show: false,
		needUserFace: false,
		needCompanyFace: false,
	};
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			if (!ruleForm.value.quota.type) {
				ElMessage({
					type: 'warning',
					message: '请选择优惠券类型!',
				});
			}

      ruleForm.value.quota.bottom = (ruleForm.value.quota.bottom || 0) * 100 | 0
      ruleForm.value.quota.send = (ruleForm.value.quota.send || 0) * 100 | 0
      ruleForm.value.quota.discount = (ruleForm.value.quota.discount || 0) * 100 | 0
      ruleForm.value.quota.amount = (ruleForm.value.quota.amount || 0) * 100 | 0

			if (dialog.type === 'add') {
				const middleData = JSON.parse(JSON.stringify(ruleForm.value));
				useCouponApiCollect.createCoupon(middleData).then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '新增优惠包成功',
						});
						emit('refresh');
					}
				});
			} else {
				const middleData = JSON.parse(JSON.stringify(ruleForm.value));
				useCouponApiCollect.updateCoupon(middleData).then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '编辑优惠包成功',
						});
						emit('refresh');
					}
				});
			}
		} else {
			return false;
		}
	});
};
onMounted(() => {});
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style lang="scss" scoped>
.ifr {
	width: 100%;
	height: 560px;
}
.tip {
	display: flex;
	flex-direction: row;
	justify-content: center;
}
</style>
