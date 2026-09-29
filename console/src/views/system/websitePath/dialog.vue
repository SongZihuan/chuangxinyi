<template>
	<div class="system-role-dialog-container">
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="900px">
			<el-form ref="roleDialogFormRef" :model="ruleForm" size="default" label-width="170px" :rules="rules" @close="closeDialog">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="路由" prop="path">
							<el-input v-model="ruleForm.path" placeholder="请输入路由路径" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="路由模式" prop="mode">
							<el-select v-model="ruleForm.mode" placeholder="请选择" clearable class="w100">
								<el-option :label="item.label" :value="item.value" v-for="(item, index) in pathModeDict" :key="index"></el-option>
							</el-select>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="policy是or还是and关系" prop="isOr">
							<el-radio-group v-model="ruleForm.isOr">
								<el-radio :label="true">or</el-radio>
								<el-radio :label="false">and</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="路由状态" prop="status">
							<el-radio-group v-model="ruleForm.status">
								<el-radio :label="1">启用</el-radio>
								<el-radio :label="2">禁用</el-radio>
								<el-radio :label="3">不放行</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>

					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="路由方法" prop="method">
							<el-checkbox-group v-model="ruleForm.method">
								<el-checkbox v-for="childItem in methodDict" :key="childItem" :label="childItem.value">{{ childItem.label }}</el-checkbox>
							</el-checkbox-group>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="路由权限" prop="policy">
							<div style="display: flex; flex-direction: column">
								<el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAllChange">全选</el-checkbox>
								<el-checkbox-group v-model="ruleForm.policy" @change="handleCheckedCitiesChange">
									<el-checkbox v-for="childItem in permissionsSign" :key="childItem" :label="childItem.value">{{ childItem.label }}</el-checkbox>
								</el-checkbox-group>
							</div>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="路由描述" prop="describe">
							<el-input type="textarea" v-model="ruleForm.describe" placeholder="请输入路由描述" clearable :rows="4"></el-input>
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
import { pathModeDict, methodDict } from '/@/api/system/path/types';
import { message } from '/@/utils/message';
import { type FormRules, type FormInstance } from 'element-plus';
import { usePathApi } from '/@/api/system/backPath';
const checkAll = ref<boolean>(false);
const isIndeterminate = ref(true);
let allCheck = ref<any>([]);
const permissionsSign = ref<dictTypes[]>([]);
const emit = defineEmits(['refresh']);
const pathApi = usePathApi();
// 定义变量内容
const roleDialogFormRef = ref();
let ruleForm = ref<any>({
	id: '',
	path: '',
	describe: '',
	mode: '',
	status: 1,
	isOr: '',
	policy: [],
});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '',
	submitTxt: '',
});

const rules = reactive<FormRules>({
	path: [{ required: true, message: '请输入路由路径', trigger: 'blur' }],
	mode: [{ required: true, message: '请输入路由模式', trigger: 'blur' }],
	status: [{ required: true, message: '请输入路由状态', trigger: 'blur' }],
	isOr: [{ required: true, message: '请选择', trigger: 'blur' }],
	method: [{ required: true, message: '请选择', trigger: 'blur' }],
	corsMode: [{ required: true, message: '请选择跨域模式', trigger: 'blur' }],
});

const reset = () => {
	ruleForm.value = {
		id: '',
		path: '',
		describe: '',
		mode: '',
		status: 1,
		isOr: '',
		policy: [],
	};
};
// 打开弹窗
const openDialog = (type: string, row: any) => {
	dialog.isShowDialog = true;
	dialog.type = type;
	if (type === 'edit') {
		nextTick(() => {
			ruleForm.value = JSON.parse(JSON.stringify(row));
			if (ruleForm.value.policy && ruleForm.value.policy.length > 0) {
				ruleForm.value.policy = ruleForm.value.policy.map((item: any) => {
					return item.sign;
				});
			}
			const checkedCount = ruleForm.value.policy.length;
			checkAll.value = checkedCount === permissionsSign.value.length;
			isIndeterminate.value = checkedCount > 0 && checkedCount < permissionsSign.value.length;
		});
		dialog.title = '修改路由';
		dialog.submitTxt = '修 改';
	} else {
		reset();
		dialog.title = '新增路由';
		dialog.submitTxt = '新 增';
	}
};
const closeDialog = () => {
	roleDialogFormRef.value?.resetFields();
	dialog.isShowDialog = false;
};
const handleCheckedCitiesChange = (value: string[]) => {
	const checkedCount = value.length;
	checkAll.value = checkedCount === permissionsSign.value.length;
	isIndeterminate.value = checkedCount > 0 && checkedCount < permissionsSign.value.length;
};
const getPermissions = () => {
	pathApi.allPermissions().then((res: any) => {
		if (res.code === "SUCCESS") {
			permissionsSign.value = res.data.permissions;
			allCheck.value = permissionsSign.value.map((item) => {
				return item.value;
			});
		}
	});
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
				pathApi.addAdminpath(ruleForm.value).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('新增成功', { type: 'success' });
						closeDialog();
						emit('refresh');
					}
				});
			} else {
				pathApi.pathUpdate({ ...ruleForm.value, id: ruleForm.value.id }).then((res: any) => {
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
	getPermissions();
});
// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss"></style>
