<template>
	<!-- 更新邮箱 -->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
		<el-form size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="mobileFormRef">
			<el-col :span="24" class="mb20">
				<el-form-item prop="webhook">
					<el-input text placeholder="请输入webhook" v-model="ruleForm.webhook" clearable autocomplete="off">
						<template #prefix>
							<SvgIcon name="my-robot" color="#61687c" :size="22"></SvgIcon>
						</template>
					</el-input>
				</el-form-item>
			</el-col>
			<el-form-item>
				<el-button round type="primary" v-waves class="login-content-submit" @click="onSubmit(mobileFormRef)">
					<span>{{ dialog.btnText }}</span>
				</el-button>
			</el-form-item>
		</el-form>
	</el-dialog>
</template>
<script setup lang="ts" name="userEmail">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { useRegisterApi } from '/@/api/register/index';
import { checkWxRobotTypes } from '/@/api/base/types';

const emit = defineEmits(['refresh']);

interface Props {
	wxRobotForm: checkWxRobotTypes;
}

const props = withDefaults(defineProps<Props>(), {});
const useRegisterCollect = useRegisterApi();
// 定义变量内容
const ruleForm = ref<checkWxRobotTypes>(props.wxRobotForm);
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	webhook: [{ required: true, message: '请输入webhook', trigger: 'blur' }],
});
const dialog = reactive({
	isShowDialog: false,
	title: '添加微信机器人',
	btnText: '绑定',
});
const openDialog = (type: string) => {
	if (type == 'add') {
		dialog.title = '添加微信机器人';
		dialog.btnText = '添加';
		ruleForm.value.webhook = '';
	} else if (type == 'edit') {
		dialog.title = '更新微信机器人';
		dialog.btnText = '更新';
	}
	dialog.isShowDialog = true;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate(async (valid) => {
		if (valid) {
			await useRegisterCollect.updateWxRobot(ruleForm.value).then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage.success('机器人绑定成功');
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
