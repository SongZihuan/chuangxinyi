<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<el-form size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="mobileFormRef">
			<el-form-item prop="password">
				<el-input text placeholder="请设置密码" v-model="ruleForm.password" clearable autocomplete="off" show-password>
					<template #prefix>
						<el-icon class="el-input__icon"><ele-Unlock /></el-icon>
					</template>
				</el-input>
			</el-form-item>
			<el-form-item prop="checkPass">
				<el-input text placeholder="请再次输入密码" v-model="ruleForm.checkPass" clearable autocomplete="off" show-password>
					<template #prefix>
						<el-icon class="el-input__icon"><ele-Unlock /></el-icon>
					</template>
				</el-input>
			</el-form-item>
			<el-form-item>
				<el-button round type="primary" v-waves class="login-content-submit" @click="onSubmit(mobileFormRef)">
					<span>设置</span>
				</el-button>
			</el-form-item>
		</el-form>
	</el-dialog>
</template>
<script setup lang="ts" name="registerPassword">
import { reactive, ref } from 'vue';
import { useRegisterApi } from '/@/api/register/index';
import { useBaseApi } from '/@/api/base/index';
import { type FormRules, type FormInstance } from 'element-plus';
import sha256 from 'sha256';
import type { passwordTypes } from '/@/api/register/types';
import { ElMessage } from 'element-plus';
const useRegisterCollect = useRegisterApi();
const emit = defineEmits(['refresh']);
const salt = ref('');
// 定义变量内容
const ruleForm = ref<passwordTypes>({
	password: '',
	checkPass: '',
});
const validatePass2 = (rule: any, value: any, callback: any) => {
	if (ruleForm.value.checkPass === '') {
		callback(new Error('请再次输入新密码'));
	} else if (ruleForm.value.checkPass !== ruleForm.value.password) {
		callback(new Error('两次密码不一致!'));
	} else {
		callback();
	}
};
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
	checkPass: [{ trigger: 'blur', validator: validatePass2 }],
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '设置密码',
	submitTxt: '确认',
});
const openDialog = () => {
	dialog.isShowDialog = true;
	reset();
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const getslat = () => {
	useBaseApi()
		.salt()
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				salt.value = res.data.salt;
			}
		});
};
const reset = () => {
	ruleForm.value.password = '';
	ruleForm.value.checkPass = '';
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			await useRegisterCollect.updatePassword({ newPasswordHash: sha256(salt.value + ":" + ruleForm.value.password), isDelete: false }).then((res: any) => {
				if (res.code === "SUCCESS") {
					closeDialog();
					ElMessage.success('设置密码成功');
					emit('refresh');
				}
			});
		}
	});
};
getslat();
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
