<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" @close="closeDialog">
    <el-form v-loading="dialog.loading" ref="roleDialogFormRef" :model="ruleForm" size="default" label-width="120px" :rules="rules">
      <el-row :gutter="35">
        <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="公司" prop="company">
            <el-input v-model="ruleForm.company" placeholder="请输入公司名称" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="地址" prop="address">
            <el-input v-model="ruleForm.address" placeholder="请输入地址" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="电话" prop="phone">
            <el-input v-model="ruleForm.phone" placeholder="请输入电话" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="ruleForm.email" placeholder="请输入邮箱" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="微信" prop="wechat">
            <el-input v-model="ruleForm.wechat" placeholder="请输入微信" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="QQ" prop="qq">
            <el-input v-model="ruleForm.qq" placeholder="请输入QQ" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="性别" prop="sex">
            <el-radio-group v-model="ruleForm.sex">
              <el-radio label="男">男</el-radio>
              <el-radio label="女">女</el-radio>
              <el-radio label="保密">保密</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="行业" prop="industry">
            <el-input v-model="ruleForm.industry" placeholder="请输入行业" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="职位" prop="position">
            <el-input v-model="ruleForm.position" placeholder="请输入职位" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
          <el-form-item label="官网链接" prop="link">
            <el-input v-model="ruleForm.link" placeholder="请输入官网链接" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
          <el-form-item label="个人简介" prop="introduction">
            <el-input v-model="ruleForm.introduction" type="textarea" placeholder="请输入个人简介" clearable></el-input>
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
import {reactive, ref} from 'vue';
import {message} from '/@/utils/message';
import {type FormRules, type FormInstance} from 'element-plus';
import {homePageTypes} from '/@/api/homePage/types';
import {useHomePageApi} from '/@/api/homePage';
import {Session} from "/@/utils/storage";

const emit = defineEmits(['refresh']);
const props = defineProps({
  homePageData: {
    type: Object,
    default: () => {
    },
  },
});
// 定义变量内容
const roleDialogFormRef = ref();
let ruleForm = ref<homePageTypes>({
  company: '',
  introduction: '',
  address: '',
  phone: '',
  email: '',
  wechat: '',
  qq: '',
  sex: '保密',
  industry: '',
  position: '',
  link: '',
});
const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: '',
  title: '',
  loading: false,
  submitTxt: '修改',
});
const rules = reactive<FormRules>({

});

// 打开弹窗
const openDialog = () => {
  ruleForm.value = JSON.parse(JSON.stringify(props.homePageData));
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
  formEl.validate(async (valid) => {
    if (valid) {
      dialog.loading = true;
      await useHomePageApi()
          .updateHomePage({...ruleForm.value, isDelete: false})
          .then((res: any) => {
            dialog.loading = false;
            if (res.code === "SUCCESS") {
              message('提交成功', {type: 'success'});
              emit('refresh', Session.get('userInfo').user.id);
              closeDialog();
            } else {
              return false;
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
