<template>
  <el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="869px" @close="closeDialog"
             destroy-on-close>
    <el-descriptions :column="column" border>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr4">
              <ele-Avatar/>
            </el-icon>
            微信头像
          </div>
        </template>
        <el-image
            :style="{ width: `40px`, height: `40px`, borderRadius: '6px' }"
            :src="userInfo.wechatHeader"
            :zoom-rate="2"
            :preview-src-list="[userInfo.wechatHeader]"
            preview-teleported
            v-if="userInfo.wechatHeader"
            fit="cover"
            close-on-press-escape
        />
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <SvgIcon name="my-weixin" :size="14" class="mr5" color="#606266"></SvgIcon>
            微信名称
          </div>
        </template>
        {{ userInfo.wechatNickName }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr4">
              <ele-User/>
            </el-icon>
            用户ID
          </div>
        </template>
        {{ userInfo.id }}
        <el-tooltip class="box-item" effect="dark" content="点击复制" placement="top-start">
          <el-icon class="ml4 cursor-pointer" @click="copyUserId(userInfo.id)">
            <ele-CopyDocument/>
          </el-icon>
        </el-tooltip>
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr4">
              <ele-User/>
            </el-icon>
            用户名称
          </div>
        </template>
        {{ userInfo.userName }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr4">
              <ele-User/>
            </el-icon>
            用户昵称
          </div>
        </template>
        {{ userInfo.nickname }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr4">
              <ele-Iphone/>
            </el-icon>
            手机号码
          </div>
        </template>
        {{ userInfo.phone }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="el-input__icon">
              <ele-Message/>
            </el-icon>
            邮箱
          </div>
        </template>
        {{ userInfo.email }}
      </el-descriptions-item>
      <el-descriptions-item label="注册时间">
        <template #label>
          <div class="cell-item">
            <SvgIcon name="my-peoples" :size="14" color="#606266"></SvgIcon>
            注册时间
          </div>
        </template>
        {{ dayjs.unix(userInfo.createAt).format('YYYY-MM-DD HH:mm:ss') || '-' }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <SvgIcon name="my-status" :size="14" color="#606266"></SvgIcon>
            状态
          </div>
        </template>
        <el-text class="mx-1" :type="userColorEnumTypes[userInfo.status]" size="mini">
          {{ UserStatusEnumTypes[userInfo.status] }}
        </el-text>
      </el-descriptions-item>
    </el-descriptions>
  </el-dialog>
</template>
<script setup lang="ts">
import {reactive} from "vue";
import dayjs from "dayjs";
import {userColorTypes, UserStatusEnum, userTypes} from "/src/data/enum";
import {ref} from "vue-demi";
import {fatherUserType} from "/@/views/onlineUser/types";

const state = reactive({
  dialog: {
    isShowDialog: false,
    title: '父级信息查看',
    submitTxt: '',
  },
  ruleForm: {},
});
const userColorEnumTypes: userTypes = userColorTypes;
const UserStatusEnumTypes: userTypes = UserStatusEnum;
const column = 1;
const userInfo = ref<fatherUserType>({
  id: '',
  roleID: 0,
  roleName: '',
  roleSign: '',
  phone: '',
  userName: '',
  nickname: '',
  header: '',
  email: '',
  userRealName: '',
  companyName: '',
  wechatNickName: '',
  wechatHeader: '',
  unionID: '',
  signin: false,
  status: '',
  inviteCount: 0,
  tokenExpire: 0,
  createAt: 0,
})
const openDialog = (row?: any,title="父级信息查看") => {
  if (!row) return;
  userInfo.value = JSON.parse(JSON.stringify(row));
  state.dialog.title = title;
  state.dialog.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
  state.dialog.isShowDialog = false;
};
defineExpose({
  openDialog,
  closeDialog,
});
</script>
<style scoped lang="scss">

</style>