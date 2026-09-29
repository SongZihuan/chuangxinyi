<template>
	<div class="system-role-dialog-container">
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px">
			<el-form ref="roleDialogFormRef" :model="ruleForm" size="default" label-width="90px" :rules="rules" @close="closeDialog">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="角色名称" prop="name">
							<el-input v-model="ruleForm.name" placeholder="请输入角色名称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="角色标识" prop="sign">
							<template #label>
								<el-tooltip effect="dark" content="用于 `router/route.ts` meta.roles" placement="top-start">
									<span>角色标识</span>
								</el-tooltip>
							</template>
							<el-input v-model="ruleForm.sign" placeholder="请输入角色标识" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="角色状态">
							<el-switch
								v-model="ruleForm.status"
								inline-prompt
								active-text="启用"
								inactive-text="禁用"
								:active-value="1"
								:inactive-value="2"
							></el-switch>
						</el-form-item>
					</el-col>

					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="删除权限" prop="roles">
							<el-radio-group v-model="ruleForm.notDelete" :disabled="dialog.type === 'edit'">
								<el-radio :label="true">是</el-radio>
								<el-radio :label="false">否</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="编辑sign" prop="roles">
							<el-radio-group v-model="ruleForm.notChangeSign" :disabled="dialog.type === 'edit'">
								<el-radio :label="true">是</el-radio>
								<el-radio :label="false">否</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="修改权限" prop="roles">
							<el-radio-group v-model="ruleForm.notChangePermissions" :disabled="dialog.type === 'edit'">
								<el-radio :label="true">是</el-radio>
								<el-radio :label="false">否</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>
					<!--          站点增加-->
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="站点" prop="belong">
							<el-select v-model="ruleForm.belong" placeholder="请选择站点" clearable>
								<el-option v-for="item in siteList" :key="item.id" :label="item.name" :value="item.id"></el-option>
							</el-select>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="角色描述">
							<el-input v-model="ruleForm.describe" type="textarea" placeholder="请输入角色描述" maxlength="150"></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="内置权限" prop="roles">
							<el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAllChange">全选 </el-checkbox>
							<el-checkbox-group v-model="ruleForm.permission" @change="handleCheckedCitiesChange">
								<el-checkbox v-for="childItem in permissionsSign" :key="childItem" :label="childItem.value">
									{{ childItem.label }}
								</el-checkbox>
							</el-checkbox-group>
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
	</div>
</template>

<script setup lang="ts" name="systemRoleDialog">
import { reactive, ref, nextTick, onMounted } from 'vue';
import { roleApi } from '/@/api/system/role/index';
import { message } from '/@/utils/message';
import { type FormRules, type FormInstance } from 'element-plus';
import type { roleTypes } from '/@/api/system/role/types';
import { useSiteApi } from '/@/api/website/site';

const emit = defineEmits(['refresh']);
const checkAll = ref<boolean>(false);
const isIndeterminate = ref(true);
const roleApiCollect = roleApi();
let allCheck = ref<any>([]);
// 定义变量内容
const roleDialogFormRef = ref();
const permissionsSign = ref<dictTypes[]>([]);
let ruleForm = ref<roleTypes>({
	name: '', // 角色名称
	sign: '', // 角色标识
	status: 2, // 角色状态
	describe: '', // 角色描述
	policy: [], //菜单权限
	belong: '', // 站点
	notDelete: false,
	notChangeSign: false,
	notChangePermissions: false,
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '',
	submitTxt: '',
});
const siteList = ref<any[]>([]);
const rules = reactive<FormRules>({
	name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
	sign: [{ required: true, message: '请输入角色标识', trigger: 'blur' }],
	sort: [{ required: true, message: '请输入排序', trigger: 'blur' }],
	belong: [{ required: true, message: '请选择站点', trigger: 'blur' }],
});
// 获取站点列表
const useSiteApiCollect = useSiteApi();
const getSiteList = () => {
	useSiteApiCollect
		.siteList({
			page: 1,
			pagesize: 9999,
			domain: '',
		})
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				siteList.value = [{id: 0, name: "用户中心全局"}, ...res.data.website];
			}
		});
};
const reset = () => {
	ruleForm.value = {
		name: '', // 角色名称
		sign: '', // 角色标识
		status: 2, // 角色状态
		describe: '', // 角色描述
		policy: [], //菜单权限
		notDelete: false,
		notChangeSign: false,
		notChangePermissions: false,
	};
};
// 打开弹窗
const openDialog = (type: string, row: any) => {
	dialog.isShowDialog = true;
	dialog.type = type;
	getSiteList();
	if (type === 'edit') {
		nextTick(() => {
			ruleForm.value = JSON.parse(JSON.stringify(row));
      if (ruleForm.value.permission && ruleForm.value.permission.length > 0) {
        ruleForm.value.permission = ruleForm.value.permission.map((item: any) => {
          return item.sign;
        });
        const checkedCount = ruleForm.value.permission.length;
        checkAll.value = checkedCount === permissionsSign.value.length;
        isIndeterminate.value = checkedCount > 0 && checkedCount < permissionsSign.value.length;
      }
		});
		dialog.title = '修改角色';
		dialog.submitTxt = '修 改';
	} else {
		reset();
		dialog.title = '新增角色';
		dialog.submitTxt = '新 增';
	}
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
			if (dialog.type === 'add') {
				roleApiCollect.roleAdd(ruleForm.value).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('新增成功', { type: 'success' });
						closeDialog();
						emit('refresh');
					}
				});
			} else {
				roleApiCollect.roleUpdate({ ...ruleForm.value, id: ruleForm.value.id }).then((res: any) => {
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
const getPermissions = () => {
	roleApiCollect.allPermissions().then((res: any) => {
		if (res.code === "SUCCESS") {
			permissionsSign.value = res.data.permissions;
			allCheck.value = permissionsSign.value.map((item) => {
				return item.value;
			});
		}
	});
};
const handleCheckedCitiesChange = (value: string[]) => {
	const checkedCount = value.length;
	checkAll.value = checkedCount === permissionsSign.value.length;
	isIndeterminate.value = checkedCount > 0 && checkedCount < permissionsSign.value.length;
};
const handleCheckAllChange = (val: boolean) => {
	ruleForm.value.policy = val ? allCheck : [];
	if (val) {
		allCheck.value = permissionsSign.value.map((item) => {
			return item.value;
		});
	}
	isIndeterminate.value = false;
};
onMounted(() => {
	getPermissions();
});
// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss"></style>
