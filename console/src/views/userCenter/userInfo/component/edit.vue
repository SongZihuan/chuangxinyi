<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px">
		<el-form ref="roleDialogFormRef" :model="ruleForm" size="default" label-width="90px" :rules="rules" @close="closeDialog">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="姓名" prop="name">
						<el-input v-model="ruleForm.name" placeholder="请输入姓名" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="手机" prop="phone">
						<el-input v-model="ruleForm.phone" placeholder="请输入手机" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="邮箱" prop="email">
						<el-input v-model="ruleForm.email" placeholder="请输入邮箱" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="国家" prop="country">
						<el-select v-model="ruleForm.country" placeholder="请选择国家" style="width: 100%" filterable>
							<el-option v-for="(item, index) in country" :value="item.codeNumber" :label="item.country" :key="index" />
						</el-select>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20" v-if="ruleForm.country == 156">
					<el-form-item label="地区" prop="area">
						<el-cascader v-model="ruleForm.area" :options="areas" :props="areasProps" style="width: 100%" ref="treeRef" />
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="地址" prop="address">
						<el-input v-model="ruleForm.address" type="textarea" placeholder="请输入详细地址" maxlength="150"></el-input>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="onCancel" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(roleDialogFormRef)" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script setup lang="ts" name="systemRoleDialog">
import { reactive, ref } from 'vue';
import { message } from '/@/utils/message';
import { areas } from '/@/data/areas';
import { country } from '/@/data/country';
import { type FormRules, type FormInstance } from 'element-plus';
import { useUserApi } from '/@/api/user/user';
import type { addressUpdateypes } from '/@/api/user/user/types';
const useUserApiCollect = useUserApi();
const emit = defineEmits(['refresh']);
const areasProps = {
	expandTrigger: 'hover' as const,
	value: 'code',
	label: 'value',
};
const props = defineProps({
	userInfo: {
		type: Object,
		default: () => {},
	},
});
// 定义变量内容
const roleDialogFormRef = ref();
let ruleForm = ref<addressUpdateypes>({
	name: '',
	phone: '',
	email: '',
	province: '',
	city: '',
	district: '',
	address: '',
	area: [],
	country: '',
});
const treeRef = ref();
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '',
	submitTxt: '修改',
});
const rules = reactive<FormRules>({
	name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
	phone: [{ required: false, message: '请输入手机', trigger: 'blur' }],
	email: [{ required: false, message: '请输入邮箱', trigger: 'blur' }],
	area: [{ required: false, message: '请选择地区', trigger: 'blur' }],
	address: [{ required: false, message: '请输入地址', trigger: 'blur' }],
	country: [{ required: false, message: '请选择国家', trigger: 'blur' }],
});

// 打开弹窗
const openDialog = () => {
	ruleForm.value = JSON.parse(JSON.stringify(props.userInfo));
	if (ruleForm.value.area instanceof Array && ruleForm.value.area.length >= 4) {
		ruleForm.value.country = ruleForm.value.area[0];
	}
	dialog.isShowDialog = true;
};

const closeDialog = () => {
	roleDialogFormRef.value?.resetFields();
	dialog.isShowDialog = false;
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
      let arr = []

      if (!ruleForm.value.area) {
        ruleForm.value.area = []
      }

      if (treeRef.value && ruleForm.value.country) {
        arr = treeRef.value.cascaderPanelRef.checkedNodes[0].text.split('/');
        ruleForm.value.area?.unshift(ruleForm.value.country + '');
      }

			useUserApiCollect
				.addressUpdate({
					...ruleForm.value,
					province: arr[0] ? arr[0].replace(/\s*/g, '') : '',
					city: arr[1] ? arr[1].replace(/\s*/g, '') : '',
					district: arr[2] ? arr[2].replace(/\s*/g, '') : '',
				})
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						message('编辑成功', { type: 'success' });
						emit('refresh');
						closeDialog();
					}
				});
		} else {
			return false;
		}
	});
};
// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss"></style>
