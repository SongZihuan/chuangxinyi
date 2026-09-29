<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="800px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="90px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="站点名称" prop="name">
						<el-input v-model="ruleForm.name" placeholder="请输入站点名称" clearable></el-input>
					</el-form-item>
				</el-col>
<!--				<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="域名" prop="domain">
						<el-input v-model="ruleForm.domain" placeholder="请输入域名" clearable></el-input>
					</el-form-item>
				</el-col>-->
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" v-if="dialog.type === 'add'">
					<el-form-item label="公钥" prop="pubkey">
						<el-input v-model="ruleForm.pubkey" type="textarea" placeholder="请输入公钥" :autosize="{ minRows: 2, maxRows: 4 }"></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="标签列表" prop="keyMap" v-for="(item, index) in ruleForm.keyMap" :key="index">
						<el-input v-model="item.label" placeholder="标签" style="width: 35%"></el-input>
						<el-input v-model="item.value" placeholder="标签值" style="width: 35%" class="ml10"></el-input>
						<el-button type="primary" @click="handleAdd" class="ml5" text v-if="index === 0">新增</el-button>
						<el-button type="danger" @click="handleDel(index)" text v-else>删除</el-button>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="内置权限" prop="policy">
            <div style="display: flex; flex-direction: column">
              <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAllChange" :key="allKey">全选</el-checkbox>
              <el-checkbox-group v-model="ruleForm.policy" @change="handleCheckedCitiesChange">
                <el-checkbox v-for="childItem in permissionsSign" :key="childItem" :label="childItem.value">{{ childItem.label }}</el-checkbox>
              </el-checkbox-group>
            </div>
          </el-form-item>
				</el-col>
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="启用" prop="status">
<!--              WebsiteStatusOK     = 1
              WebsiteStatusBanned = 2-->
                <el-radio-group v-model="ruleForm.status">
                    <el-radio :label="1">启用</el-radio>
                    <el-radio :label="2">禁用</el-radio>
                </el-radio-group>
            </el-form-item>
        </el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="描述" prop="describe">
						<el-input v-model="ruleForm.describe" type="textarea" placeholder="请输入描述" maxlength="150"></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="协议" prop="name"> <WangEditor v-model:get-html="ruleForm.agreement" /></el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue';
import WangEditor from '/@/components/Editor/index.vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { useSiteApi } from '/@/api/website/site/index';
import { Base64 } from 'js-base64';
import type { websiteFromType } from '/@/api/website/site/types';
const useSiteApiCollect = useSiteApi();
const isIndeterminate = ref<boolean>(false);
const checkAll = ref<boolean>(false);
const allKey = ref(0);
const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<websiteFromType>({
	name: '',
	describe: '',
	domain: '',
	pubkey: '',
	agreement: '',
	policy: [],
  status: 1,
	keyMap: [{ value: '', label: '' }],
});
let allCheck = ref<any>([]);
const rules = reactive<FormRules>({
	name: [{ required: true, message: '请输入站点名称', trigger: 'blur' }],
	domain: [{ required: true, message: '请输入域名', trigger: 'blur' }],
	describe: [{ required: true, message: '请输入描述', trigger: 'blur' }],
	pubkey: [{ required: true, message: '请输入公钥', trigger: 'blur' }],
	agreement: [{ required: true, message: '请输入协议', trigger: 'blur' }],
	policy: [{ required: true, message: '请选择权限', trigger: 'blur' }],
});

const permissionsSign = ref<dictTypes[]>([]);
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '站点新增',
	submitTxt: '新增',
});

const openDialog = (type: string, row: any) => {
	if (type === 'add') {
		reset();
		dialog.title = '站点新增';
		dialog.submitTxt = '新增';
	} else {
    ruleForm.value = JSON.parse(JSON.stringify(row));
    if (ruleForm.value.policy && ruleForm.value.policy.length > 0) {
      ruleForm.value.policy = ruleForm.value.policy.map((item: any) => {
        return item.sign;
      });
      const checkedCount = ruleForm.value.policy.length;
      checkAll.value = checkedCount === permissionsSign.value.length;
      isIndeterminate.value = checkedCount > 0 && checkedCount < permissionsSign.value.length;
    }

		dialog.title = '站点编辑';
		dialog.submitTxt = '编辑';
	}
	dialog.type = type;
	dialog.isShowDialog = true;
};
const handleAdd = () => {
	ruleForm.value.keyMap.push({ value: '', label: '' });
};
const handleDel = (index: number) => {
	ruleForm.value.keyMap.splice(index, 1);
};
//重置
const reset = () => {
	allCheck.value = false;
	ruleForm.value = {
		name: '',
		describe: '',
		domain: '',
		pubkey: '',
		agreement: '',
    status: 1,
		policy: [],
		keyMap: [{ value: '', label: '' }],
	};
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			if (dialog.type === 'add') {
				const middleData = JSON.parse(JSON.stringify(ruleForm.value));
				middleData.pubkey = "base64:" + Base64.encode(middleData.pubkey);
				useSiteApiCollect.createSite(middleData).then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '新增站点成功',
						});
						emit('refresh');
					}
				});
			} else {
				const middleData = JSON.parse(JSON.stringify(ruleForm.value));
				middleData.pubkey = "base64:" + Base64.encode(middleData.pubkey);
				useSiteApiCollect.updateSite(middleData).then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '编辑站点成功',
						});
						emit('refresh');
					}
				});
			}
		} else {
			return false;
		}
	});
};
const handleCheckAllChange = (val: boolean) => {
	ruleForm.value.policy = val ? permissionsSign.value : [];
	if (val) {
		ruleForm.value.policy = permissionsSign.value.map((item) => {
			return item.value;
		});
	}
	isIndeterminate.value = false;
};
const handleCheckedCitiesChange = (value: string[]) => {
	const checkedCount = value.length;
	checkAll.value = checkedCount === permissionsSign.value.length;
	isIndeterminate.value = checkedCount > 0 && checkedCount < permissionsSign.value.length;
};
const getPermissions = () => {
	useSiteApiCollect.allPermissions().then((res: any) => {
		if (res.code === "SUCCESS") {
			permissionsSign.value = res.data.permissions;
			ruleForm.value.policy = permissionsSign.value.map((item:any) => {
				return item.value?.sign;
			});
			isIndeterminate.value = false;
		}
	});
};
onMounted(() => {
	getPermissions();
});
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
	flex-direction: row;
	justify-content: center;
}
</style>
