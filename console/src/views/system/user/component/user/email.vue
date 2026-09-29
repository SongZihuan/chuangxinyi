<template>
	<!-- 更新邮箱 -->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
		<el-form size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="mobileFormRef">
			<el-col :span="24" class="mb20">
				<el-form-item prop="email">
					<el-input text placeholder="请输入邮箱" v-model="ruleForm.email" clearable autocomplete="off">
						<template #prefix>
							<SvgIcon name="iconfont icon-youxiang"></SvgIcon>
						</template>
					</el-input>
				</el-form-item>
			</el-col>
			<el-form-item>
				<el-button round type="primary" v-waves class="login-content-submit" @click="onSubmit(mobileFormRef)">
					<span>绑定</span>
				</el-button>
			</el-form-item>
		</el-form>
	</el-dialog>
</template>
<script setup lang="ts" name="userEmail">
import { reactive, ref } from 'vue';
import { verifyEmail } from '/@/utils/toolsValidate';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { userApi } from '/@/api/system/user';
import { userEmailTypes } from '/@/views/system/user/types';

const emit = defineEmits(['refresh']);

interface Props {
	emailForm: userEmailTypes;
}

const props = withDefaults(defineProps<Props>(), {});
// 定义变量内容
const ruleForm = ref<userEmailTypes>(props.emailForm);
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	email: [{ trigger: 'blur', validator: verifyEmail }],
	code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '修改绑定的邮箱',
	submitTxt: '确认',
});
const openDialog = () => {
	dialog.isShowDialog = true;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!ruleForm.value.uid) {
		ElMessage.error('用户获取失败');
		return;
	}
	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			await userApi()
				.updateUserEmail(ruleForm.value)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('修改邮箱绑定成功');
						emit('refresh');
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
