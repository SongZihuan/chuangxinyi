<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" height="500px" center>
    <el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
      <el-row :gutter="35">
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
          <el-form-item label="主题" prop="subject">
            <el-input v-model="ruleForm.subject" placeholder="请输入主题" clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
          <el-form-item label="内容" prop="content">
            <el-input v-model="ruleForm.content" placeholder="请输入内容" clearable type="textarea"
                      :autosize="{ minRows: 2, maxRows: 6 }"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">发送</el-button>
			</span>
    </template>
  </el-dialog
  >
</template>

<script setup lang="ts">
import {reactive, toRefs, ref} from 'vue';
import {type FormRules, type FormInstance, ElMessage} from 'element-plus';
import {useMessageAdminApi} from '/@/api/message/admin';

const state = reactive({
  dialog: {
    isShowDialog: false,
    type: '',
    title: '发送邮件',
    submitTxt: '',
  },
});
let ruleForm = ref({
  subject: '',
  content: '',
  id: null,
});
const formRef = ref();
const rules = reactive<FormRules>({
  subject: [
    {
      required: true,
      message: '请输入主题',
      trigger: 'blur',
    },
  ],
  id: [
    {
      required: true,
      message: '请输入用户ID',
      trigger: 'blur',
    },
  ],
  content: [
    {
      required: true,
      message: '请输入内容',
      trigger: 'blur',
    },
  ],
});
const closeDialog = () => {
  dialog.value.isShowDialog = false;
};

const onSubmit = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate((valid) => {
    if (valid) {
      useMessageAdminApi()
          .sendEmail({...ruleForm.value, id: Number(ruleForm.value.id)})
          .then((res: any) => {
            if (res.code === "SUCCESS" && !res.data.have) {
              ElMessage({
                type: 'error',
                message: '当前用户暂未绑定邮箱',
              });
              return;
            }
            if (res.code === "SUCCESS" && res.data.success) {
              ElMessage({
                type: 'success',
                message: '发送邮件成功!',
              });
              dialog.value.isShowDialog = false;
            } else {
              ElMessage({
                type: 'error',
                message: '发送失败',
              });
            }
          });
    } else {
      return false;
    }
  });
};
const {dialog} = toRefs(state);
const openDialog = (row: any) => {
  ruleForm.value.id = row.id;
  ruleForm.value = {
    subject: '',
    content: '',
    id: row.numberID || row.userID || row.id,
  };
  dialog.value.isShowDialog = true;
};

defineExpose({
  openDialog,
});
</script>

<style lang="scss" scoped>
.subtitle {
  color: var(--el-text-color-secondary);
  text-align: center;
  padding-bottom: 20px;
}
</style>
