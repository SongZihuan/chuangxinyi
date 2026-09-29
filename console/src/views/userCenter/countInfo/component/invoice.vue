<template>
	<!-- 更新发票抬头 -->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
		<el-form size="large" class="login-content-form" :rules="type == '1' ? rules : rules2" :model="ruleForm" ref="invoiceFormRef">
			<el-col :span="24" class="mb20">
				<el-form-item prop="type">
					<el-radio-group v-model="type">
						<el-radio label="1">个人</el-radio>
						<el-radio label="2">企业</el-radio>
					</el-radio-group>
				</el-form-item>
			</el-col>
			<template v-if="type == '1'">
				<el-col :span="24" class="mb20">
					<el-form-item prop="name">
						<el-input text placeholder="请输入用户发票抬头" v-model="ruleForm.name" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="my-people"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</el-col>
				<el-col :span="24" class="mb20">
					<el-form-item prop="taxID">
						<el-input text placeholder="请输入身份证号" v-model="ruleForm.taxID" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="my-idcard"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</el-col>
				<el-col :span="24" class="mb20">
					<el-form-item prop="bank">
						<el-input text placeholder="请输入开户银行" v-model="ruleForm.bank" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="ele-Postcard"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</el-col>
				<el-col :span="24" class="mb20">
					<el-form-item prop="bankID">
						<el-input text placeholder="请输入银行账号" v-model="ruleForm.bankID" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="my-bank-account"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</el-col>
			</template>
			<template v-else>
				<el-col :span="24" class="mb20">
					<el-form-item prop="name">
						<el-input text placeholder="请输入公司发票抬头" v-model="ruleForm.name" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="my-peoples"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</el-col>
				<el-col :span="24" class="mb20">
					<el-form-item prop="taxID">
						<el-input text placeholder="请输入税号" v-model="ruleForm.taxID" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="ele-Cellphone"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</el-col>
				<el-col :span="24" class="mb20">
					<el-form-item prop="bank">
						<el-input text placeholder="请输入开户银行" v-model="ruleForm.bank" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="ele-Postcard"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</el-col>
				<el-col :span="24" class="mb20">
					<el-form-item prop="bankID">
						<el-input text placeholder="请输入银行账号" v-model="ruleForm.bankID" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="ele-Memo"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</el-col>
			</template>
			<el-form-item>
				<el-button round type="primary" v-waves class="login-content-submit" @click="onSubmit(invoiceFormRef)">
					<span>绑定</span>
				</el-button>
			</el-form-item>
		</el-form>
	</el-dialog>
</template>
<script setup lang="ts" name="userInvoice">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { invoiceTitleTypes } from '/@/api/invoice/title/types';
import { useInvoiceTitleApi } from '/@/api/invoice/title';

const emit = defineEmits(['refresh']);

interface Props {
	invoiceForm: invoiceTitleTypes;
}

const props = withDefaults(defineProps<Props>(), {});
// 定义变量内容
const ruleForm = ref<invoiceTitleTypes>(props.invoiceForm);
// 类型
const type = ref<string>('1');
const invoiceFormRef = ref();
const rules = reactive<FormRules>({
	name: [
		{
			required: true,
			message: '请输入用户名',
			trigger: 'blur',
		},
	],
	taxID: [
		{
			required: true,
			message: '请输入身份证',
			trigger: 'blur',
		},
	],
});
const rules2 = reactive<FormRules>({
	name: [
		{
			required: true,
			message: '请输入发票抬头',
			trigger: 'blur',
		},
	],
	taxID: [
		{
			required: true,
			message: '请输入纳税人识别号',
			trigger: 'blur',
		},
	],
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '修改绑定的发票抬头',
	submitTxt: '确认',
});
const openDialog = () => {
	dialog.isShowDialog = true;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			await useInvoiceTitleApi()
				.updateInvoice(ruleForm.value)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '修改发票抬头成功',
						});
						emit('refresh');
					}
				});
		}
	});
};
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss">
.login-content-form {
	margin-top: 20px;
	@for $i from 1 through 4 {
		.login-animation#{$i} {
			opacity: 0;
			animation-name: error-num;
			animation-duration: 0.5s;
			animation-fill-mode: forwards;
			animation-delay: calc($i/10) + s;
		}
	}

	.login-content-code {
		width: 100%;
		padding: 0;
	}

	.login-content-submit {
		width: 100%;
		letter-spacing: 2px;
		font-weight: 300;
		margin-top: 15px;
	}

	.login-msg {
		color: var(--el-text-color-placeholder);
	}
}
</style>
