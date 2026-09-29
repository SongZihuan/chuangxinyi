<template>
	<div class="system-menu-dialog-container">
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" destroy-on-close @close="closeDialog">
			<el-form ref="menuDialogFormRef" :model="ruleForm" size="default" label-width="90px" :rules="rules">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="上级菜单">
							<el-cascader
								:options="props.menuList"
								:props="{ checkStrictly: true, value: 'id', label: 'title' }"
								placeholder="请选择上级菜单"
								clearable
								class="w100"
								v-model="ruleForm.menuSuperior"
							>
								<template #default="{ node, data }">
									<input hidden v-bind="ruleForm.parentID" value="{{data.menuId}}" />
									<span>{{ data.title }}</span>
									<span v-if="!node.isLeaf"> ({{ data.children.length }}) </span>
								</template>
							</el-cascader>
						</el-form-item>
					</el-col>
					<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="菜单类型">
							<el-radio-group v-model="ruleForm.menuType">
								<el-radio :label="1">菜单</el-radio>
								<el-radio :label="2">按钮</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="菜单名称" prop="title">
							<el-input v-model="ruleForm.title" placeholder="格式：首页" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="token权限" prop="subPolicy">
							<el-checkbox-group v-model="ruleForm.subPolicy">
								<el-checkbox v-for="childItem in subAll" :key="childItem" :label="childItem.value">{{ childItem.label }}</el-checkbox>
							</el-checkbox-group>
						</el-form-item>
					</el-col>

					<template v-if="ruleForm.menuType === 1">
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="路由名称" prop="name">
								<el-input v-model="ruleForm.name" placeholder="路由中的 name 值" clearable></el-input>
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="路由路径" prop="path">
								<el-input v-model="ruleForm.path" placeholder="路由中的 path 值" clearable></el-input>
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="重定向" prop="redirect">
								<el-input v-model="ruleForm.redirect" placeholder="请输入路由重定向" clearable></el-input>
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="菜单图标" prop="icon">
								<IconSelector placeholder="请输入菜单图标" v-model="ruleForm.icon" />
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="组件路径" prop="componentAlias">
								<el-input v-model="ruleForm.componentAlias" placeholder="组件路径" clearable></el-input>
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="链接地址" prop="metaIsLink">
								<el-input v-model="ruleForm.metaIsLink" placeholder="外链/内嵌时链接地址（http:xxx.com）" clearable :disabled="!ruleForm.isLink">
								</el-input>
							</el-form-item>
						</el-col>
					</template>
					<template v-if="ruleForm.menuType === 2">
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="权限标识" prop="btnPower">
								<el-input v-model="ruleForm.btnPower" placeholder="请输入权限标识" clearable></el-input>
							</el-form-item>
						</el-col>
					</template>
					<template v-if="ruleForm.menuType === 1">
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="是否隐藏" prop="isHide">
								<el-radio-group v-model="ruleForm.isHide">
									<el-radio :label="true">隐藏</el-radio>
									<el-radio :label="false">不隐藏</el-radio>
								</el-radio-group>
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="页面缓存" prop="isKeepAlive">
								<el-radio-group v-model="ruleForm.isKeepAlive">
									<el-radio :label="true">缓存</el-radio>
									<el-radio :label="false">不缓存</el-radio>
								</el-radio-group>
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="是否固定" prop="isAffix">
								<el-radio-group v-model="ruleForm.isAffix">
									<el-radio :label="true">固定</el-radio>
									<el-radio :label="false">不固定</el-radio>
								</el-radio-group>
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="是否外链">
								<el-radio-group v-model="ruleForm.isLink" :disabled="ruleForm.isIframe">
									<el-radio :label="true">是</el-radio>
									<el-radio :label="false">否</el-radio>
								</el-radio-group>
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
							<el-form-item label="是否内嵌" prop="isIframe">
								<el-radio-group v-model="ruleForm.isIframe" @change="onSelectIframeChange">
									<el-radio :label="true">是</el-radio>
									<el-radio :label="false">否</el-radio>
								</el-radio-group>
							</el-form-item>
						</el-col>
					</template>
					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="状态" prop="status">
							<el-radio-group v-model="ruleForm.status">
								<el-radio :label="1">启用</el-radio>
								<el-radio :label="2">禁用</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>

					<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
						<el-form-item label="状态选择" prop="roles">
							<el-radio-group v-model="ruleForm.isOr">
								<el-radio :label="true">任意选中权限就显示该菜单</el-radio>
								<el-radio :label="false">全部选中权限才显示</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="内置权限" prop="policy">
							<div class="flex-col">
								<div>
									<el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAllChange">全选</el-checkbox>
								</div>
								<div>
									<el-checkbox-group v-model="ruleForm.policy" @change="handleCheckedCitiesChange">
										<el-checkbox v-for="childItem in permissionsSign" :key="childItem" :label="childItem.value">{{ childItem.label }}</el-checkbox>
									</el-checkbox-group>
								</div>
							</div>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="描述" prop="describe">
							<el-input type="textarea" v-model="ruleForm.describe" placeholder="请输入描述" clearable :rows="4"></el-input>
						</el-form-item>
					</el-col>
				</el-row>
			</el-form>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit(menuDialogFormRef)" size="default">{{ dialog.submitTxt }}</el-button>
				</span>
			</template>
		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="systemMenuDialog">
import { defineAsyncComponent, reactive, ref, nextTick, onMounted } from 'vue';
import { useMenuApi } from '/@/api/system/menu/index';
import { type FormRules, type FormInstance } from 'element-plus';
import type { menuTypes } from '/@/api/system/menu/types';
import { message } from '/@/utils/message';
import { roleApi } from '/@/api/system/role/index';
const checkAll = ref<boolean>(false);
const permissionsSign = ref<dictTypes[]>([]);
const props = defineProps({
	menuList: {
		type: Array,
		default: () => [],
	},
});
const roleApiCollect = roleApi();
const isIndeterminate = ref(true);
const emit = defineEmits(['refresh']);
const IconSelector = defineAsyncComponent(() => import('/@/components/iconSelector/index.vue'));
const menuDialogFormRef = ref();
let ruleForm = ref<menuTypes>({
	menuSuperior: [],
	menuType: 1,
	name: '',
	parentID: 0,
	menuCategory: 1,
	component: '',
	componentAlias: '',
	isLink: false,
	path: '',
	redirect: '',
	title: '',
	icon: '',
	isHide: false,
	isKeepAlive: true,
	isAffix: true,
	metaIsLink: '',
	isIframe: false,
	roles: [] as number[],
	btnPower: '',
	menuId: 0,
	describe: '',
	status: 1,
	isOr: true,
	policy: [],
	subPolicy: [],
});
let allCheck = ref<any>([]);
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '',
	submitTxt: '',
});
const rules = reactive<FormRules>({
	title: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
	name: [{ required: true, message: '请输入路由名称', trigger: 'blur' }],
	path: [{ required: true, message: '请输入路由路径', trigger: 'blur' }],
	icon: [{ required: true, message: '请输入菜单图标', trigger: 'blur' }],
	componentAlias: [{ required: true, message: '请输入组件路径', trigger: 'blur' }],
	subPolicy: [{ required: true, message: '请选择token权限', trigger: 'blur' }],
	btnPower: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
	sort: [{ required: true, message: '请输入菜单排序', trigger: 'blur' }],
});
const menuApi = useMenuApi();
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
	checkAll.value = checkedCount == permissionsSign.value.length;
	isIndeterminate.value = false;
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
const openDialog = (type: string, row?: any) => {
	dialog.isShowDialog = true;
	dialog.type = type;
	if (type == 'edit') {
		nextTick(() => {
			ruleForm.value = JSON.parse(JSON.stringify(row));
			if (ruleForm.value.menuSuperior && ruleForm.value.menuSuperior.indexOf(',') != -1) {
				ruleForm.value.menuSuperior = ruleForm.value.menuSuperior.split(',');
			} else {
				let id = parseInt(ruleForm.value.menuSuperior);
				ruleForm.value.menuSuperior = [];
				ruleForm.value.menuSuperior.push(id);
			}
			if (ruleForm.value.policy && ruleForm.value.policy.length > 0) {
				ruleForm.value.policy = ruleForm.value.policy.map((item: any) => {
					return item.sign;
				});
				const checkedCount = ruleForm.value.policy.length;
				checkAll.value = checkedCount === permissionsSign.value.length;
				isIndeterminate.value = checkedCount > 0 && checkedCount < permissionsSign.value.length;
			}
		});
		dialog.title = '修改菜单';
		dialog.submitTxt = '修 改';
	} else {
		reset();
		ruleForm.value.menuSuperior = [];
		dialog.title = '新增菜单';
		dialog.submitTxt = '新 增';
	}
};
const reset = () => {
	ruleForm.value = {
		menuSuperior: [],
		menuType: 1,
		name: '',
		parentID: 0,
		menuCategory: 1,
		component: '',
		componentAlias: '',
		isLink: false,
		path: '',
		redirect: '',
		title: '',
		icon: '',
		isHide: false,
		isKeepAlive: true,
		isAffix: false,
		metaIsLink: '',
		isIframe: false,
		roles: [] as number[],
		btnPower: '',
		menuId: 0,
		describe: '',
		status: 1,
		isOr: true,
		policy: [],
		subPolicy: [],
	};
};
const subAll = ref([]);
const closeDialog = () => {
	menuDialogFormRef.value?.resetFields();
	dialog.isShowDialog = false;
};

// 是否内嵌下拉改变
const onSelectIframeChange = () => {
	if (ruleForm.value.isIframe) ruleForm.value.isLink = false;
	else ruleForm.value.isLink = true;
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
			if (ruleForm.value.menuSuperior && ruleForm.value.menuSuperior.length > 0) {
				ruleForm.value.parentID = parseInt(ruleForm.value.menuSuperior[ruleForm.value.menuSuperior.length - 1]);
			}

			if (dialog.type === 'add') {
				menuApi.addAdminMenu({ ...ruleForm.value, menuSuperior: ruleForm.value.menuSuperior?.toString() }).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('新增成功', { type: 'success' });
						closeDialog(); // 关闭弹窗
						emit('refresh');
					}
				});
			} else if (dialog.type === 'edit') {
				if (!ruleForm.value.menuSuperior) ruleForm.value.parentID = 0;
				menuApi.menuUpdate({ ...ruleForm.value, menuSuperior: ruleForm.value.menuSuperior?.toString() }).then((res: any) => {
					if (res.code === "SUCCESS") {
						message('编辑成功', { type: 'success' });
						closeDialog(); // 关闭弹窗
						emit('refresh');
					}
				});
			}
		} else {
			return false;
		}
	});
};
const getSubAll = () => {
	menuApi.subAll().then((res: any) => {
		if (res.code === "SUCCESS") {
			subAll.value = res.data.permissions;
		}
	});
};
onMounted(() => {
	getSubAll();
	getPermissions();
});
// 暴露变量
defineExpose({
	openDialog,
});
</script>
