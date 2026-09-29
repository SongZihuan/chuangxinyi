<template>
	<!-- 更新用户名 昵称 -->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
		<el-form v-loading="dialog.loading" size="large" class="login-content-form" :rules="rules" :model="ruleForm" ref="mobileFormRef">
			<el-col :span="24" class="mb20">
				<el-form-item prop="userName" v-if="dialog.title == '更新用户名'">
					<el-input text placeholder="请输入用户名" v-model="ruleForm.userName" clearable autocomplete="off">
						<template #prefix>
							<SvgIcon name="ele-User"></SvgIcon>
						</template>
					</el-input>
				</el-form-item>
				<el-form-item prop="nickname" v-else>
					<el-input text placeholder="请输入用户昵称" v-model="ruleForm.nickname" clearable autocomplete="off">
						<template #prefix>
							<SvgIcon name="ele-User"></SvgIcon>
						</template>
					</el-input>
				</el-form-item>
			</el-col>
			<el-form-item>
				<el-button round type="primary" v-waves class="login-content-submit" @click="onSubmit(mobileFormRef)">
					<span>{{ dialog.submitTxt }}</span>
				</el-button>
			</el-form-item>
		</el-form>
	</el-dialog>
</template>
<script setup lang="ts" name="userName">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { useRegisterApi } from '/@/api/register/index';
import { updateUserNameTypes } from '/@/api/register/types';
import { Session, Local } from '/@/utils/storage';

const emit = defineEmits(['refresh']);

interface Props {
	userNameForm: updateUserNameTypes;
}
const props = withDefaults(defineProps<Props>(), {});
const useRegisterCollect = useRegisterApi();
// 定义变量内容
const ruleForm = ref<updateUserNameTypes>(props.userNameForm);
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	userName: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
	nickname: [{ required: true, message: '请输入用户昵称', trigger: 'blur' }],
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'nickName',
	title: '更新用户名',
	submitTxt: '确认',
  loading: false,
});
const openDialog = (type: string) => {
	if (type == 'userName') {
		dialog.title = '更新用户名';
		dialog.type == 'userName';
	} else {
		dialog.title = '更新昵称';
		dialog.type == 'nickName';
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
			let userInfos = Session.get('userInfo') || Local.get('userInfo');
			if (!userInfos) return ElMessage.error('请先登录');
			let data: updateUserNameTypes = {
				userName: ruleForm.value.userName,
				isDelete: false,
				uid: userInfos.user.id,
				nickname: ruleForm.value.nickname,
			};
      dialog.loading = true;
			if (dialog.title == '更新用户名') {
				await useRegisterCollect.updateUserName(data).then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('更新用户名成功');
						setTimeout(() => {
							emit('refresh');
						}, 500);

						closeDialog();
					}
          dialog.loading = false;
				});
			} else {
				await useRegisterCollect.updateNickname(data).then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage.success('更新昵称成功');
						setTimeout(() => {
							emit('refresh');
						}, 500);
						closeDialog();
					}
				});
        dialog.loading = false;
			}
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
