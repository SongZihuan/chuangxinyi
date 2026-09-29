<template>
	<!-- 更新token有效时间-->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
		<el-form size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="mobileFormRef">
			<el-col :span="24" class="mb20">
				<el-form-item prop="tokenExpiration">
					<el-select v-model="ruleForm.tokenExpiration" placeholder="请选择有效时间" style="width: 100%">
						<el-option v-for="(item, index) in tokenDict" :value="item.value" :label="item.label" :key="index" />
					</el-select>
				</el-form-item>
				<el-form-item>
					<el-button round type="primary" v-waves class="login-content-submit" @click="onSubmit(mobileFormRef)">
						<span>{{ dialog.submitTxt }}</span>
					</el-button>
				</el-form-item>
			</el-col>
		</el-form>
	</el-dialog>
</template>
<script setup lang="ts" name="userName">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { useRegisterApi } from '/@/api/register/index';
import { tokenExpirationType } from '/@/api/register/types';

const emit = defineEmits(['refresh']);

interface Props {
	userNameForm: tokenExpirationType;
}
const props = withDefaults(defineProps<Props>(), {});
const useRegisterCollect = useRegisterApi();
const tokenDict = [
	{ label: '15分钟', value: 900 },
	{ label: '30分钟', value: 1800 },
	{ label: '1小时', value: 3600 },
	{ label: '2小时', value: 7200 },
	{ label: '3小时', value: 180 * 60 },
	{ label: '4小时', value: 4 * 60 * 60 },
	{ label: '5小时', value: 5 * 60 * 60 },
	{ label: '8小时', value: 8 * 60 * 60 },
	{ label: '12小时', value: 12 * 60 * 60 },
	{ label: '24小时', value: 24 * 60 * 60 },
	{ label: '48小时', value: 48 * 60 * 60 },
];
const ruleForm = ref<tokenExpirationType>(props.userNameForm);
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	tokenExpiration: [{ required: true, message: '请输入token有效期', trigger: 'blur' }],
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '更新令牌有效时间',
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
			await useRegisterCollect.tokenExpiration(ruleForm.value).then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('更新令牌有效时间成功');
					setTimeout(() => {
						emit('refresh');
					}, 500);
					closeDialog();
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
