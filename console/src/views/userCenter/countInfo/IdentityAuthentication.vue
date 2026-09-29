<template>
  <div class="container layout-padding">
    <el-card shadow="hover">
      <div class="main">
        <div class="user-info">
          <div class="user-info-left">
            <div class="main-title">实名认证</div>
            <div class="user-head" v-loading="avatarLoading">
              <el-tooltip v-if="authUser.updateHeader" class="box-item" effect="dark" content="更换头像" placement="top-start">
                <div style="width: 48px; height: 48px">
                  <AvatarUpload
                      ref="avatarUploadRef"
                      v-if="userInfo"
                      @imgSuccess="imgSuccess"
                      :userId="userInfo.user.id"
                      :key="avatarUploadKey"
                  ></AvatarUpload>
                </div>
              </el-tooltip>
              <div v-else style="width: 48px; height: 48px">
                <AvatarUpload
                    ref="avatarUploadRef"
                    v-if="userInfo"
                    :disabledType="!authUser.updateHeader"
                    @imgSuccess="imgSuccess"
                    :userId="userInfo.user.id"
                    :key="avatarUploadKey"
                ></AvatarUpload>
              </div>
              <div class="user-head-right">
                <div>登录账号</div>
                <div v-if="userInfo && userInfo.user.phone">
                  {{ userInfo.user.nickname || userInfo.user.userName || userInfo.user.phone || '异常用户'
                  }}<span v-if="userInfo.user.nickname">({{ userInfo.user.userName || userInfo.user.phone }})</span>-<span class="roleName">{{
                    userInfo.role.name
                  }}</span>
                </div>
              </div>
            </div>
          </div>
          <div class="user-info-right" style="width: 120px; margin-right: -30px">
            <div class="item-action" v-if="authUser.updateUserName" @click.stop="updateUserName('userName')" style="width: 60px">用户名</div>
            <div class="item-action mr10" v-if="authUser.updateUserName" @click.stop="updateUserName('nickName')" style="width: 60px">昵称</div>
            <div class="item-action" v-else>
              <el-text class="mx-1" type="info" size="mini">抱歉，暂无权限</el-text>
            </div>
          </div>
        </div>
        <div class="info-list">
          <div class="info-item">
            <el-image style="width: 24px; height: 24px" :src="user" fit="cover" lazy />
            <div class="item-main">
              <div class="item-title">个人认证</div>
              <div class="item-text" v-if="userInfo && userInfo.data.hasVerified">
                您已上传身份信息 {{ userInfo && userInfo.data.hasUserOriginal ? '已上传身份证' : '您暂未上传身份证' }}
                {{ userInfo && userInfo.data.hasUserFaceCheck ? '已扫脸' : '您暂未扫脸' }}
              </div>
              <div class="item-text" v-else>您暂未进行实名认证</div>
            </div>
            <div class="item-action" v-if="authUser.updateRealName" @click="updateUser">修改</div>
            <div class="item-action" v-else>
              <el-text type="info">抱歉，暂无权限</el-text>
            </div>
          </div>
          <div class="info-item">
            <SvgIcon name="my-peoples" color="#61687c" :size="22"></SvgIcon>
            <div class="item-main">
              <div class="item-title">企业认证</div>
              <div class="item-text" v-if="userInfo && userInfo.data.isCompany">
                您已上传企业信息 {{ userInfo && userInfo.data.hasCompanyOriginal ? '已上传身份证' : '您暂未企业营业执照 法人身份证照' }}
                {{ userInfo && userInfo.data.hasLegalPersonFaceCheck ? '已扫脸' : '您暂未扫脸' }}
              </div>
              <div class="item-text" v-else>您暂未进行实名认证</div>
            </div>
            <template v-if="authUser.getinfo">
              <div class="item-action" @click="updateCompany">修改</div>
            </template>
            <template v-else>
              <div class="item-action">
                <el-text type="info">抱歉，暂无权限</el-text>
              </div>
            </template>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 更新用户名 -->
    <user-name ref="userNameRef" :user-name-form="userNameForm" @refresh="getUserInfo" />
    <!-- 更新手机号绑定 -->
    <Mobile ref="mobileRef" :phone-form="phoneForm" @refresh="getUserInfo" />
    <!-- 邮箱绑定 -->
    <Email ref="emailRef" :email-form="emailForm" @refresh="getUserInfo" />
    <!-- 微信绑定 -->
    <Weixin ref="weixinRef" @refresh="getUserInfo" />
    <!-- 微信机器人绑定 -->
    <WxRobot ref="wxRobotRef" :wx-robot-form="wxRobotForm" @refresh="getUserInfo" />
    <!-- 2fa绑定 -->
    <TotpCom ref="totpComRef" :totp-form="totpForm" @refresh="getUserInfo"></TotpCom>
    <!-- 个人信息修改 -->
    <UploadUserInfo
        ref="userRef"
        :hasUserInfo="userInfo && userInfo.data.hasVerified"
        :hasUserOriginal="userInfo && userInfo.data.hasUserOriginal"
        :hasFace="userInfo && userInfo.data.hasUserFaceCheck"
        :user-form="userForm"
        @refresh="getUserInfo"
    ></UploadUserInfo>
    <!-- 企业信息修改 -->
    <UploadCompany
        ref="companyRef"
        :hasUserInfo="userInfo && userInfo.data.isCompany"
        :hasCompanyOriginal="userInfo && userInfo.data.hasCompanyOriginal"
        :hasFace="userInfo && userInfo.data.hasLegalPersonFaceCheck"
        :company-form="companyForm"
        @refresh="getUserInfo"
    ></UploadCompany>
    <!-- 手机验证注销弹窗 -->
    <LogOffMobile ref="logOffMobileRef" @reset="reset" />
    <!-- 服务号绑定弹窗   -->
    <weixin-service ref="weixinServiceRef" :wx-status="wxStatus" @refresh="getUserInfo"></weixin-service>
    <!-- 发票抬头 -->
    <invoice ref="invoiceRef" :invoice-form="invoiceForm" @refresh="getUserInfo"></invoice>
    <!-- 公告 -->
    <ShowHtml ref="showHtmlRef" :title="currentInfo.title" :contentHtml="currentInfo.content" :subtitle="currentInfo.subtitle" />
    <!--    登录控制-->
    <updateLoginCtr ref="updateLoginCtrRef" @refresh="getUserInfo"></updateLoginCtr>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent, watch } from 'vue';
import dayjs from 'dayjs';
import phone from '/@/assets/userCenter/phone.png';
import email from '/@/assets/userCenter/email.png';
import resgister from '/@/assets/userCenter/resgister.png';
import user from '/@/assets/userCenter/user.png';
import weixin from '/@/assets/userCenter/weixin.png';
import { Local, Session } from '/@/utils/storage';
import { useUserApi } from '/@/api/user/user';
import ShowHtml from '/@/components/ShowHtml/index.vue';
import LogOffMobile from '/@/views/userCenter/countInfo/component/logOffMobile.vue';
import Mobile from '/@/views/userCenter/countInfo/component/mobile.vue';
import Email from '/@/views/userCenter/countInfo/component/email.vue';
import TotpCom from '/@/views/userCenter/countInfo/component/totpCom.vue';
import UploadUserInfo from '/@/views/userCenter/countInfo/component/uploadUserInfo.vue';
import UploadCompany from '/@/views/userCenter/countInfo/component/uploadCompany.vue';
import Weixin from '/@/views/userCenter/countInfo/component/weixin.vue';
import invoice from '/@/views/userCenter/countInfo/component/invoice.vue';

import { ElMessage, ElMessageBox } from 'element-plus';
import { checkEmailTypes, checkPhoneTypes, checkWxRobotTypes } from '/@/api/base/types';
import { updateUserNameTypes, uploadCompanyInfoJsonTypes, uploadUserInfoJsontTypes } from '/@/api/register/types';
import WeixinService from '/@/views/userCenter/countInfo/component/weixinService.vue';
import { totpConfigTypes } from '/@/api/Totp/types';
import { useWeixinApi } from '/@/api/weixin/index';
import { invoiceTitleTypes } from '/@/api/invoice/title/types';
import useSubAuth from '/@/hooks/useSubAuth';
import UserName from '/@/views/userCenter/countInfo/component/userName.vue';
import WxRobot from '/@/views/userCenter/countInfo/component/wxRobot.vue';
import { useRegisterApi } from '/@/api/register';
import { useUserInfo } from '/@/stores/userInfo';
import { useLoginSignIn } from '/@/hooks/useLoginSignIn';
import { useSocketListInfo } from '/@/stores/socketListInfo';
import { checkPhoneRes } from '/@/views/userCenter/countInfo/types';

const AvatarUpload = defineAsyncComponent(() => import('/@/components/avatarUpload/index.vue'));
const updateLoginCtr = defineAsyncComponent(() => import('/@/views/userCenter/countInfo/component/updateLoginCtr.vue'));

// 引入组件
const logOffMobileRef = ref<any>(false);
const stores = useUserInfo();
const useWeixinApiCollect = useWeixinApi();
const useUserApiCollect = useUserApi();
const showHtmlRef = ref();
const mobileRef = ref();
const emailRef = ref();
const userInfo = ref();
const totpComRef = ref();
const userRef = ref();
const companyRef = ref();
const weixinRef = ref();
const weixinServiceRef = ref();
const invoiceRef = ref();
const userNameRef = ref();
const authUser = useSubAuth();
const wxRobotRef = ref();
const updateLoginCtrRef = ref();
const avatarUploadKey = ref(1); //上传头像组件key
const avatarLoading = ref(false); //上传头像loading
const { resetSignIn } = useLoginSignIn();
const userInfoStore = useUserInfo();
const avatarUploadRef = ref();
watch(
    () => userInfoStore.userInfos,
    (newVal) => {
      setUserInfo(newVal);
    },
    { deep: true }
);
// 发票抬头
const invoiceForm = ref<invoiceTitleTypes>({
  name: '',
  taxID: '',
  bankID: '',
  bank: '',
});
// 登录控制
let updateLoginCtrForm = ref<checkPhoneRes>({
  allowPhone: false,
  allowEmail: false,
  allowPassword: false,
  allowWeChat: false,
  allowSecondFA: false,
});
// 单点登录
const signin = ref<boolean>(false);
// 手机号
const phoneForm = ref<checkPhoneTypes>({
  phone: '',
  code: '',
  type: 'PhoneCheck',
});
const tokenForm = ref({
  tokenExpiration: '',
});
// 用户名
const userNameForm = ref<updateUserNameTypes>({
  userName: '',
  nickname: '',
});
// 邮箱
const emailForm = ref<checkEmailTypes>({
  email: '',
  code: '',
  type: 'EmailCheck',
});
// 个人认证
const userForm = ref<uploadUserInfoJsontTypes>({
  authenticationMethod: 1,
  userName: '',
  userIDCard: '',
});
// 企业认证
const companyForm = ref<uploadCompanyInfoJsonTypes>({
  authenticationMethod: 1,
  legalPersonName: '',
  legalPersonIDCard: '',
  companyName: '',
  companyID: '',
});
// 微信绑定状态 0未绑定 1已绑定 2服务号绑定 3服务号未绑定
const wxStatus = ref<number>(0);
// 2fa基本数据
const totpForm = ref<totpConfigTypes>({
  phone: '',
  secret: '',
});
// 微信机器人
const wxRobotForm = ref<checkWxRobotTypes>({
  webhook: '',
  isDelete: false,
});
// 财务信息
const financeInfo = ref<any>({
  walletID: '',
  balance: '',
  notBilled: '',
  billed: '',
  hasBilled: '',
  cny: '',
  waitBalance: '',
  withdraw: '',
  waitWithdraw: '',
  notWithdraw: '',
  hasWithdraw: '',
});
const currentInfo = ref({ title: '', content: '', subtitle: '' });
//跟新手机号
const updatePhone = () => {
  mobileRef.value.openDialog();
};
// 微信机器人
const updateWxRobot = () => {
  wxRobotRef.value.openDialog('add');
};
const useRegisterCollect = useRegisterApi();
// 删除微信机器人
const deleteWxRobot = () => {
  ElMessageBox.confirm('是否删除微信机器人?', '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(async () => {
        await useRegisterCollect.updateWxRobot({ isDelete: true }).then((res: any) => {
          if (res.code === 'SUCCESS') {
            ElMessage.success('删除微信机器人成功');
            getUserInfo();
          }
        });
      })
      .catch(() => {});
};
//更新用户名
const updateUserName = (type: string) => {
  userNameRef.value.openDialog(type);
};
//更新邮箱
const updateEmail = () => {
  emailRef.value.openDialog();
};
//更新微信
const updateWeixin = () => {
  weixinRef.value.openDialog();
};
// 更新发票抬头
const updateInvoice = () => {
  invoiceRef.value.openDialog();
};
//解绑微信
const delWeixin = () => {
  ElMessageBox.confirm('是否解绑微信?', '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        useWeixinApiCollect.bindWeixin({ isDelete: true }).then((res: any) => {
          if (res.code === 'SUCCESS') {
            ElMessage.success('解绑微信成功');
            getUserInfo();
          }
        });
      })
      .catch(() => {});
};
// 打开服务号
const openWxService = () => {
  weixinServiceRef.value.openDialog();
};
const imgSuccess = async (base64: string) => {
  avatarLoading.value = true;
  await useUserApiCollect.avatarUpdate({ header: base64, isDelete: false }).then((res: any) => {
    avatarLoading.value = false;
    if (res.code === 'SUCCESS') {
      ElMessage({
        type: 'success',
        message: '头像上传成功!',
      });
    }
  });
  await getUserInfo();
  avatarUploadKey.value++;
};


//更新用户信息
const updateUser = () => {
  userRef.value.openDialog(userForm.value);
};
//更新企业'
const updateCompany = () => {
  companyRef.value.openDialog(companyForm.value);
};
//重置
const reset = async () => {
  let res = await resetSignIn();
  if (!res) {
    Session.clear();
    Local.clear();
    window.location.reload();
  }
};
const setUserInfo = (data: any) => {
  userInfo.value = data;
  if (userInfo.value?.user) {
    phoneForm.value.phone = userInfo.value.user.phone;
    totpForm.value.phone = userInfo.value.user.phone;
    emailForm.value.email = userInfo.value.user.email;
    signin.value = userInfo.value.user.signin;
  }
  if (userInfo.value?.info) {
    userForm.value.userName = userInfo.value.info.userName;
    userForm.value.userIDCard = userInfo.value.info.userIDCard;
    if (userInfo.value?.info.isCompany) {
      companyForm.value.legalPersonName = userInfo.value.info.legalPersonName;
      companyForm.value.legalPersonIDCard = userInfo.value.info.legalPersonIdCard;
      companyForm.value.companyName = userInfo.value.info.companyName;
      companyForm.value.companyID = userInfo.value.info.companyID;
    }
  }
  if (userInfo.value?.balance) {
    financeInfo.value = userInfo.value.balance;
  }
  if (userInfo.value?.title) {
    invoiceForm.value.name = userInfo.value.title.name;
    invoiceForm.value.taxID = userInfo.value.title.taxID;
    invoiceForm.value.bankID = userInfo.value.title.bandID;
    invoiceForm.value.bank = userInfo.value.title.band;
  }
  if (userInfo.value?.data) {
    if (userInfo.value.data.hasWeChat) {
      wxStatus.value = 1;
    } else {
      wxStatus.value = 0;
    }
    if (userInfo.value.data.hasWeChat && userInfo.value.data.hasFuwuhao) {
      wxStatus.value = 2;
    }
    if (userInfo.value.data.hasWeChat && !userInfo.value.data.hasFuwuhao) {
      wxStatus.value = 3;
    }
    updateLoginCtrForm.value.allowPhone = userInfo.value.data.allowPhone;
    updateLoginCtrForm.value.allowEmail = userInfo.value.data.allowEmail;
    updateLoginCtrForm.value.allowPassword = userInfo.value.data.allowPassword;
    updateLoginCtrForm.value.allowWeChat = userInfo.value.data.allowWeChat;
    updateLoginCtrForm.value.allowSecondFA = userInfo.value.data.allowSecondFA;
  }
  if (userInfo.value?.user?.phone) {
    if (userInfo.value?.user.userName) {
      userNameForm.value.userName = userInfo.value.user.userName;
    }

    userNameForm.value.nickname = userInfo.value.user.nickname;
    tokenForm.value.tokenExpiration = userInfo.value.user.tokenExpire;
  }
  avatarUploadRef.value && avatarUploadRef.value.onUpdate();
};
const getUserInfo = async () => {
  await useUserApiCollect.userInfo().then((res: any) => {
    if (res.code === 'SUCCESS') {
      userInfo.value = res.data;
      setUserInfo(res.data);
      stores.setUserData(res.data);
      Session.set('userInfo', res.data);
      Session.set('userData', res.data);
      stores.setUserInfos();
    }
  });
};
const socketInfo = useSocketListInfo();
onMounted(async () => {
  await getUserInfo();
});
</script>

<style lang="scss" scoped>
.el-input__wrapper {
  margin-bottom: 10px;
}

.main-title {
  font-size: 16px;
  font-weight: 600;
  color: #000000;
  line-height: 24px;
}

:deep(.el-card) {
  overflow-y: auto;
}

.card {
  margin-top: 16px;
}

.notice-bar {
  display: flex;
  flex-direction: row;

  .notice-date {
    width: 42px;
    height: 54px;
    background: #f1f2f5;
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    .day {
      color: #61687c;
      height: 24px;
      font-weight: bold;
      font-size: 16px;
      line-height: 24px;
    }

    .mouth {
      height: 20px;
      font-size: 12px;
      line-height: 20px;
      color: #61687c;
    }
  }

  .notice-main {
    flex: 1;
    height: 54px;
    margin: 0px 15px;

    .content-box-item {
      height: 54px;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      cursor: pointer;
    }

    .content-box-title {
      height: 22px;
      font-size: 14px;
      font-weight: 600;
      color: #414960;
      line-height: 22px;
    }
  }

  .notice-arrow {
    width: 24px;
    height: 54px;
    margin-left: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.main {
  .user-info {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .user-info-right {
      display: flex;
      flex-direction: row;
    }

    .user-head {
      padding: 16px 0px 20px 0px;
      display: flex;
      border-bottom: 1px solid #f1f2f5;

      .user-head-right {
        margin-left: 16px;

        :nth-child(1) {
          height: 22px;
          font-size: 14px;
          font-weight: 400;
          color: #61687c;
          line-height: 22px;
        }

        :nth-child(2) {
          height: 22px;
          font-size: 14px;
          font-weight: 400;
          color: #414960;
          line-height: 22px;
          margin-top: 8px;
        }

        .roleName {
          color: var(--el-text-color-secondary);
        }
      }
    }

    .item-action {
      color: var(--el-color-primary);
      line-height: 40px;
      padding-top: 10px;
      padding-left: 8px;
      cursor: pointer;
    }
  }

  .info-list {
    .info-item {
      padding: 16px 0px 20px 0px;
      display: flex;
      flex-direction: row;
      border-bottom: 1px solid #f1f2f5;
      max-height: 200px;

      .item-main {
        flex: 1;
        margin-left: 8px;

        .item-title {
          height: 22px;
          font-size: 14px;
          font-weight: 400;
          color: #61687c;
          line-height: 22px;
        }

        .item-text {
          font-size: 14px;
          font-weight: 400;
          color: #414960;
          min-height: 22px;
          margin-top: 8px;
          display: flex;
          align-items: center;

          span {
            display: flex;
            align-items: center;
          }

          .binding {
            color: var(--el-color-primary);
            cursor: pointer;
          }
        }
      }

      .item-action {
        color: var(--el-color-primary);
        line-height: 40px;
        padding-top: 10px;
        padding-left: 8px;
        cursor: pointer;
      }
    }
  }
}
</style>
