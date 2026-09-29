<template>
  <el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px" @close="closeDialog">
  <div>
    <el-card shadow="never">
      <el-descriptions  direction="vertical" :column="column">
        <el-descriptions-item label="用户实名">{{ userInfo.userName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="用户身份证">{{ userInfo.userIDCard || '-' }}</el-descriptions-item>
<!--        身份证图片-->
        <el-descriptions-item label="身份证正面">
          <el-image v-if="userInfo.userIdCardUrl" style="width: 300px; height: 150px" :src="userInfo.userIdCardUrl"  :preview-src-list="[userInfo.userIdCardUrl]" fit="contain"></el-image>
          <div v-else>-</div>
        </el-descriptions-item>
        <el-descriptions-item label="身份证反面">
          <el-image v-if="userInfo.userIdCardBackUrl" style="width: 300px; height: 150px" :src="userInfo.userIdCardBackUrl" :preview-src-list="[userInfo.userIdCardBackUrl]" fit="contain"></el-image>
          <div v-else>-</div>
        </el-descriptions-item>
        <el-descriptions-item label="企业实名">{{ userInfo.companyName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="企业统一社会信用代码">{{ userInfo.companyID || '-' }}</el-descriptions-item>
        <el-descriptions-item label="法人姓名">{{ userInfo.legalPersonName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="法人身份证">{{ userInfo.legalPersonIDCard || '-' }}</el-descriptions-item>
<!--        法人身份证图片-->
        <el-descriptions-item label="法人身份证正面">
          <el-image v-if="userInfo.legalPersonIdCardUrl" style="width: 300px; height: 150px" :src="userInfo.legalPersonIdCardUrl" :preview-src-list="[userInfo.legalPersonIdCardUrl]" fit="contain"></el-image>
          <div v-else>-</div>
        </el-descriptions-item>
        <el-descriptions-item label="法人身份证反面">
          <el-image v-if="userInfo.legalPersonIdCardBackUrl" style="width: 300px; height: 150px" :src="userInfo.legalPersonIdCardBackUrl" :preview-src-list="[userInfo.legalPersonIdCardBackUrl]" fit="contain"></el-image>
          <div v-else>-</div>
        </el-descriptions-item>
         <el-descriptions-item label="营业执照">
          <el-image v-if="userInfo.licenseUrl" style="width: 300px; height: 150px" :src="userInfo.licenseUrl" :preview-src-list="[userInfo.licenseUrl]" fit="contain"></el-image>
          <div v-else>-</div>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
  </el-dialog>
</template>

<script setup lang="ts">
import {ref,  reactive} from 'vue';
import {userRealNameTypes} from "/@/views/system/user/types";
import {userApi} from "/@/api/system/user";
import {ElMessage} from "element-plus";
import useFile from "/@/hooks/useFile";
const column = ref(1);
const userInfo = ref<userRealNameTypes>({
  userName: '',
  userIDCard: '',
  companyName: '',
  companyID: '',
  legalPersonName: '',
  legalPersonIDCard: '',
});
const {getFile} = useFile()
const state = reactive({
  dialog: {
    isShowDialog: false,
    title: '用户实名信息',
    submitTxt: '',
  },
  ruleForm: {} as any,
});
const openDialog = (row?: any) => {
  state.dialog.isShowDialog = true;
  if(!row.id){
    ElMessage.error('用户ID不能为空');
    return;
  }
  state.ruleForm = JSON.parse(JSON.stringify(row));
  getRealNameInfo();
};
// 关闭弹窗
const closeDialog = () => {
  state.dialog.isShowDialog = false;
};
const getRealNameInfo = () => {
  userApi().getUserRealName({uid:state.ruleForm.id}).then((res: any) => {
    if (res.code === "SUCCESS") {
      userInfo.value = res.data;
    }
  });
};
defineExpose({
  openDialog,
  closeDialog,
})
</script>

<style lang="scss" scoped></style>
