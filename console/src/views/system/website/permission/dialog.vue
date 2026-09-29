<template>
  <div class="system-role-dialog-container">
    <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px">
      <el-form ref="roleDialogFormRef" :model="ruleForm" size="default" label-width="90px" :rules="rules"
               @close="closeDialog">
        <el-row :gutter="35">
<!-- Name        string       `json:"name"`
    Sign        string       `json:"sign"`
    Describe    string       `json:"describe"`
    Sort        int64        `json:"sort"`
    Status      int64        `json:"status"`-->
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="权限名称" prop="name">
              <el-input v-model="ruleForm.name" placeholder="请输入角色名称" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="权限标识" prop="sign">
              <el-input v-model="ruleForm.sign" placeholder="请输入角色标识" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="排序" prop="sort">
              <el-input v-model="ruleForm.sort" placeholder="请输入排序" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="权限状态" prop="status">
              <el-radio-group v-model="ruleForm.status">
                <el-radio :label="2">启用</el-radio>
                <el-radio :label="1">禁用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="权限描述" prop="describe">
              <el-input v-model="ruleForm.describe" placeholder="请输入角色描述" clearable type="textarea" :rows="4"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit(roleDialogFormRef)" size="default">{{
              dialog.submitTxt
            }}</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="systemRoleDialog">
import {reactive, ref, nextTick, onMounted} from 'vue';

import {message} from '/@/utils/message';
import {type FormRules, type FormInstance} from 'element-plus';
import {useSiteApi} from "/@/api/website/site";
import {usePermissionApi} from "/@/api/system/permission";

const emit = defineEmits(['refresh']);
const permissionApi = usePermissionApi();
// 定义变量内容
const roleDialogFormRef = ref();
let ruleForm = ref<any>({
  name: '',
  sign: '',
  sort: '',
  status: 1,
  describe: '',
});

const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: '',
  title: '',
  submitTxt: '',
});
const siteList = ref<any[]>([]);
const rules = reactive<FormRules>({
  name: [{required: true, message: '请输入权限名称', trigger: 'blur'}],
  sign: [{required: true, message: '请输入权限标识', trigger: 'blur'}],
  sort: [{required: true, message: '请输入排序', trigger: 'blur'}],
  status: [{required: true, message: '请选择权限状态', trigger: 'change'}],
});
// 获取站点列表
const useSiteApiCollect = useSiteApi();
const getSiteList = () => {
  useSiteApiCollect.siteList({
    page: 1,
    pagesize: 9999,
    domain: '',
  }).then((res: any) => {
    if (res.code === "SUCCESS") {
      siteList.value = res.data.website;
    }
  });
};
const reset = () => {
  ruleForm.value = {
    name: '', // 角色名称
    sign: '', // 角色标识
    status: 2, // 角色状态
    describe: '', // 角色描述
    sort: '', // 排序
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
        permissionApi.addPermissionMenu(ruleForm.value).then((res: any) => {
          if (res.code === "SUCCESS") {
            message('新增成功', {type: 'success'});
            closeDialog();
            emit('refresh');
          }
        });
      } else {
        permissionApi.permissionUpdate({...ruleForm.value, id: ruleForm.value.roleID}).then((res: any) => {
          if (res.code === "SUCCESS") {
            message('编辑成功', {type: 'success'});
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
});
// 暴露变量
defineExpose({
  openDialog,
});
</script>

<style scoped lang="scss"></style>
