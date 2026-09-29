<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px" @close="closeDialog">
			<el-form ref="userDialogFormRef" :model="state.ruleForm" size="default" label-width="130px" :rules="rules">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="账户名称" prop="userName">
							<el-input v-model="state.ruleForm.userName" placeholder="请输入账户名称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="用户昵称" prop="nikeName">
							<el-input v-model="state.ruleForm.nikeName" placeholder="请输入用户昵称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="关联角色" prop="roleSign">
							<el-select
								v-model="state.ruleForm.roleSign"
								placeholder="请选择"
								multiple
								collapse-tags
								collapse-tags-tooltip
								:max-collapse-tags="1"
								clearable
								filterable
								class="w100"
							>
								<el-option :label="item.roleName" :value="item.roleId" v-for="(item, index) in roleDict" :key="index"></el-option>
							</el-select>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="用户类型" prop="userType">
							<el-select v-model="state.ruleForm.userType" placeholder="请选择" clearable class="w100">
								<el-option label="企业" :value="1"></el-option>
								<el-option label="个人" :value="2"></el-option>
							</el-select>
						</el-form-item>
					</el-col>

					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="公司名称" prop="socialName">
							<el-input v-model="state.ruleForm.socialName" placeholder="请输入公司名称" clearable></el-input>
						</el-form-item>
					</el-col>

					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="手机号" prop="phone">
							<el-input v-model="state.ruleForm.phone" placeholder="请输入手机号" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="邮箱" prop="email">
							<el-input v-model="state.ruleForm.email" placeholder="请输入" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="性别" prop="sex">
							<el-select v-model="state.ruleForm.sex" placeholder="请选择" clearable class="w100">
								<el-option label="男" :value="1"></el-option>
								<el-option label="女" :value="2"></el-option>
							</el-select>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="身份证号" prop="idcard">
							<el-input v-model="state.ruleForm.idcard" placeholder="请输入" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="统一社会信用代码" prop="socialCode">
							<el-input v-model="state.ruleForm.socialCode" placeholder="请输入" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="VIP等级" prop="">
							<el-input-number v-model="state.ruleForm.vipLevel" placeholder="请输入" clearable></el-input-number>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20" v-if="state.dialog.type === 'add'">
						<el-form-item label="账户密码" prop="password">
							<el-input v-model="state.ruleForm.password" placeholder="请输入" type="password" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20" v-else>
						<el-form-item label="账户密码">
							<el-input v-model="state.ruleForm.password" placeholder="请输入" type="password" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="账户过期" prop="overdueTime">
							<el-date-picker v-model="state.ruleForm.overdueTime" type="date" placeholder="请选择" class="w100"> </el-date-picker>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="用户状态" prop="status">
							<el-select v-model="state.ruleForm.userStatus" placeholder="请选择" clearable class="w100">
								<el-option label="待激活" :value="0"></el-option>
								<el-option label="激活" :value="1"></el-option>
								<el-option label="冻结" :value="2"></el-option>
								<el-option label="封禁" :value="3"></el-option>
							</el-select>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="用户描述" prop="describe">
							<el-input v-model="state.ruleForm.describe" type="textarea" placeholder="请输入用户描述" maxlength="150"></el-input>
						</el-form-item>
					</el-col>
				</el-row>
			</el-form>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit(userDialogFormRef)" size="default">{{ state.dialog.submitTxt }}</el-button>
				</span>
			</template>
		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="systemUserDialog">
import { reactive, ref, nextTick, onMounted } from 'vue';
import type { userTypes } from '../types';
import { type FormRules, type FormInstance } from 'element-plus';
import { userApi } from '/@/api/system/user';
import { message } from '/@/utils/message';
import { formatGetTime } from '/@/utils/formatTime';
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);
const userApiCollect = userApi();
const roleDict = ref<any>([]);

// 定义变量内容
const userDialogFormRef = ref<FormInstance>();
const state = reactive({
	ruleForm: {
		userName: '',
		userNickname: '',
		roleSign: null,
		phone: '',
		email: '',
		sex: 1,
		password: '',
		overdueTime: '',
		userStatus: 1,
		describe: '',
		socialCode: '',
		userType: 1,
		socialName: '',
	} as unknown as userTypes,
	dialog: {
		isShowDialog: false,
		type: '',
		title: '',
		submitTxt: '',
	} as dialogParams,
});

const rules = reactive<FormRules>({
	userName: [{ required: true, message: '请输入账户名称', trigger: 'blur' }],
	userNickname: [{ required: true, message: '请输入用户昵称', trigger: 'blur' }],
	userType: [{ required: true, message: '请选择用户类型', trigger: 'blur' }],
	phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
	email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }],
	sex: [{ required: true, message: '请输入性别', trigger: 'blur' }],
	idcard: [{ required: true, message: '请输入身份证号', trigger: 'blur' }],
	password: [{ required: true, message: '请输入账号密码', trigger: 'blur' }],
});
// 打开弹窗
const openDialog = (type: string, row?: any) => {
	state.dialog.isShowDialog = true;
	state.dialog.type = type;
	if (type === 'edit') {
		nextTick(() => {
			setTimeout(() => {
				state.ruleForm = JSON.parse(JSON.stringify(row));
			}, 500);
		});
		state.dialog.title = '修改用户';
		state.dialog.submitTxt = '修 改';
	} else {
		state.dialog.title = '新增用户';
		state.dialog.submitTxt = '新 增';
	}
};
// 关闭弹窗
const closeDialog = () => {
	userDialogFormRef.value?.resetFields();
	state.dialog.isShowDialog = false;
};
// 取消
const onCancel = () => {
	closeDialog();
};
// 提交
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			const middleData = JSON.parse(JSON.stringify(state.ruleForm));
			middleData.overdueTime = formatGetTime(middleData.overdueTime);
			if (state.dialog.type === 'add') {
				userApiCollect.userAdd(middleData).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('新增成功', { type: 'success' });
						closeDialog();
						emit('refresh');
					}
				});
			} else {
				middleData.password ? null : (middleData.password = '');
				userApiCollect.userUpdate(middleData).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('编辑成功', { type: 'success' });
						closeDialog();
						emit('refresh');
					}
				});
			}
		} else {
			return false;
		}
	});
};
onMounted(() => {
	// getRoleDict();
});
// 暴露变量
defineExpose({
	openDialog,
});
</script>
