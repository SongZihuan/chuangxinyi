<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<!-- 剩余余额  -->
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="78px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="用户状态" prop="status">
						<el-select v-model="ruleForm.status" placeholder="请选择用户状态类型">
<!--              注销和封禁只能在注销和封禁之间转换，不能转换为其他状态
              <el-option label="注册" value="REGISTER"></el-option>
							<el-option label="正常" value="NORMAL"></el-option>
							<el-option label="禁用" value="BANNED"></el-option>
							<el-option label="删除" value="DELETE"></el-option>
							<el-option label="冻结" value="FREEZE"></el-option>-->
              <template v-if="status === 'REGISTER' || status === 'NORMAL' || status === 'FREEZE'">
                <el-option label="注册" value="REGISTER"></el-option>
                <el-option label="正常" value="NORMAL"></el-option>
                <el-option label="禁用" value="BANNED"></el-option>
                <el-option label="注销" value="DELETE"></el-option>
                <el-option label="冻结" value="FREEZE"></el-option>
              </template>
              <template v-else-if="status === 'DELETE'">
                <el-option label="禁用" value="BANNED"></el-option>
                <el-option label="注销" value="DELETE"></el-option>
              </template>
              <template v-else-if="status === 'BANNED'">
                <el-option label="禁用" value="BANNED"></el-option>
                <el-option label="注销" value="DELETE"></el-option>
              </template>

						</el-select>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
		</template>
		<checkPhone ref="checkPhoneRef" @checkPhoneSuccess="checkPhoneSuccess" />
	</el-dialog>
</template>

<script setup lang="ts" name="invoiceDialog">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import checkPhone from '/@/components/checkPhone/index.vue';
import { ElMessage } from 'element-plus';
import { userApi } from '/@/api/system/user';
import { userStatusTypes } from '/@/views/system/user/types';

const formRef = ref();
const emit = defineEmits(['refresh']);
const status = ref<string>('');
const checkPhoneRef = ref();
let ruleForm = ref<userStatusTypes>({
	uid: '',
	status: '',
});
const rules = reactive<FormRules>({
	status: [
		{
			required: true,
			message: '请选择用户状态类型',
			trigger: 'blur',
		},
	],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '设置用户状态',
	submitTxt: '确定',
});

const openDialog = (row: { uid: string; status: string }) => {
	if (!row.uid) {
		ElMessage({
			type: 'error',
			message: '用户ID获取失败',
		});
		return;
	}
	status.value = row.status;
	ruleForm.value.uid = row.uid;
	ruleForm.value.status = row.status;
	dialog.isShowDialog = true;
};
//重置
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	if (!ruleForm.value.uid) {
		ElMessage.error('用户获取失败');
		return;
	}
	formEl.validate((valid) => {
		if (valid) {
			checkPhoneRef.value.openDialog({}, '用户状态确认');
		} else {
			return false;
		}
	});
};
const checkPhoneSuccess = (phoneToken: string) => {
	userApi()
		.updateUserStatus({ ...ruleForm.value, phoneToken: phoneToken })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				closeDialog();
				ElMessage({
					type: 'success',
					message: '用户状态更新成功',
				});
				emit('refresh');
			}
		});
};
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
	padding: 10px 10px;
}
</style>
