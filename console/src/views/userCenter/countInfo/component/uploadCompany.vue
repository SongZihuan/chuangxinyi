<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
    <div v-loading="dialog.loading">
      <el-steps :active="stepActive" align-center class="mb20">
        <el-step title="企业信息填写"/>
        <el-step title="企业信息上传"/>
        <el-step title="企业认证"/>
      </el-steps>
      <el-form v-if="stepActive == 1" size="large" class="login-content-form" :rules="rules" :model="ruleForm"
               ref="mobileFormRef">
        <el-form-item prop="legalPersonName">
          <el-input text placeholder="请输入企业法人真实姓名" v-model="ruleForm.legalPersonName" clearable
                    autocomplete="off">
            <template #prefix>
              <SvgIcon name="my-people"></SvgIcon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="legalPersonIDCard">
          <el-input text placeholder="请输入企业法人身份证号码" v-model="ruleForm.legalPersonIDCard" clearable
                    autocomplete="off">
            <template #prefix>
              <SvgIcon name="my-idcard"></SvgIcon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="companyName">
          <el-input text placeholder="请输入公司名称" v-model="ruleForm.companyName" clearable autocomplete="off">
            <template #prefix>
              <SvgIcon name="my-peoples"></SvgIcon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="companyID">
          <el-input text placeholder="请输入公司统一社会信用代码" v-model="ruleForm.companyID" clearable
                    autocomplete="off">
            <template #prefix>
              <SvgIcon name="my-passport"></SvgIcon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
         <div class="btns">
           <el-button round type="primary" v-waves class="login-content-submit" @click="onNext(mobileFormRef)">
             <span>绑定</span>
           </el-button>
           <el-button  v-if="hasCompanyOriginal" round type="primary" v-waves class="login-content-submit" @click="onNext(mobileFormRef)">
             <span>下一步</span>
           </el-button>
         </div>
        </el-form-item>
      </el-form>
      <el-form v-if="stepActive == 2" size="large" class="login-content-form" :rules="rules" :model="ruleForm">
        <el-form-item prop="phone">
          <imageUpload
              :upImgBoxCustomStyle="upImgBoxCustomStyle"
              @imgSuccess="idcardImgSuccess"
              imgText="支持jpg/png/bmp；文件大小不能超过2M"
              imgUpText="身份证正面照"
          />
        </el-form-item>
        <el-form-item prop="phone">
          <imageUpload
              :upImgBoxCustomStyle="upImgBoxCustomStyle"
              @imgSuccess="idcardNationalEmblemImgSuccess"
              imgText="支持jpg/png/bmp；文件大小不能超过2M"
              imgUpText="身份证国徽照"
          />
        </el-form-item>
        <el-form-item prop="phone">
          <imageUpload
              :upImgBoxCustomStyle="upImgBoxCustomStyle"
              @imgSuccess="licenseImgSuccess"
              imgText="支持jpg/png/bmp；文件大小不能超过2M"
              imgUpText="营业执照"
          />
        </el-form-item>
        <el-form-item>
          <div class="btns">
            <el-button round type="primary" v-waves class="btn" @click="goBackCompanyInfo">
              <span>上一步</span>
            </el-button>
            <el-button v-if="isBack" type="primary" v-waves @click="onSubmit()" block class="btn" round> 下一步
            </el-button>
            <el-button v-else type="primary" v-waves @click="onSubmit()" block class="btn" round> 上传</el-button>
          </div>
        </el-form-item>
      </el-form>
      <scan-qr-codes
          v-if="stepActive == 3"
          current-user="UpdateCompanyInfo"
          :scan-qr-data="ruleForm"
          @go-back="goBack"
          @success="scanSuccess"
      ></scan-qr-codes>
      <el-empty v-if="stepActive == 4" description="您已完成企业信息上传" :image="checkCircle">
        <!--  修改信息      -->
        <el-button type="primary" v-waves @click="editInfo" block class="btn" round> 修改信息</el-button>
      </el-empty>
    </div>
  </el-dialog>
</template>
<script setup lang="ts" name="userUploadCompany">
import {reactive, ref} from 'vue';
import imageUpload from '/@/components/imageUpload/index.vue';
import {useRegisterApi} from '/@/api/register/index';
import {verifyPhone} from '/@/utils/toolsValidate';
import {type FormRules, type FormInstance} from 'element-plus';
import type {uploadCompanyInfoJsonTypes} from '/@/api/register/types';
import {ElMessage, ElLoading} from 'element-plus';
import {ElNotification} from 'element-plus';
import ScanQrCodes from '/@/views/userCenter/countInfo/component/scanQrCodes.vue';
import checkCircle from "/@/assets/check-circle.png";

const useRegisterCollect = useRegisterApi();
const emit = defineEmits(['refresh']);

interface Props {
  companyForm: uploadCompanyInfoJsonTypes;
  hasUserInfo: boolean;
  hasCompanyOriginal: boolean;
  hasFace: boolean;
}

const props = withDefaults(defineProps<Props>(), {});
const stepActive = ref(1);
// 定义变量内容
const oldCompanyForm = ref<uploadCompanyInfoJsonTypes>(props.companyForm);
const ruleForm = ref<uploadCompanyInfoJsonTypes>({
  authenticationMethod: 1,
  legalPersonName: '',
  legalPersonIDCard: '',
  companyName: '',
  companyID: '',
});
const mobileFormRef = ref();
const rules = reactive<FormRules>({
  phone: [{trigger: 'blur', validator: verifyPhone}],
  code: [{required: true, message: '请输入验证码', trigger: 'blur'}],
});
const upImgBoxCustomStyle = {
  width: '100%',
  height: '150px',
};
const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: '',
  loading: false,
  title: '修改绑定的企业信息',
  submitTxt: '确认',
});
const openDialog = (form:any) => {
  if(!form) return;
  ruleForm.value = JSON.parse(JSON.stringify(form));
  dialog.isShowDialog = true;
  if (props.hasUserInfo && props.hasCompanyOriginal && props.hasFace) {
    stepActive.value = 4;
    ElNotification({
      title: '提示',
      message: '您已完成企业信息上传 !',
      type: 'warning',
    });
  } else if (props.hasUserInfo && props.hasCompanyOriginal) {
    stepActive.value = 3;
    ElNotification({
      title: '提示',
      message: '不上传身份证信息和不扫脸认证部分功能不能使用 !',
      type: 'warning',
    });
  } else if (props.hasUserInfo) {
    stepActive.value = 2;
    ElNotification({
      title: '提示',
      message: '不上传身份证信息部分功能不能使用 !',
      type: 'warning',
    });
  } else if (props.hasCompanyOriginal) {
    stepActive.value = 3;
    ElNotification({
      title: '提示',
      message: '不扫脸认证部分功能不能使用 !',
      type: 'warning',
    });
  }
};
const editInfo = () => {
  stepActive.value = 1;
  ElNotification.closeAll();
};
const closeDialog = () => {
  dialog.isShowDialog = false;
  ElNotification.closeAll();
  clearImg();
};
const idcardImgFile = ref();
//图片组件传过来的文件数据
const idcardImgSuccess = (file: any) => {
  idcardImgFile.value = file;
};
const idcardNationalEmblemImgFile = ref();
const idcardNationalEmblemImgSuccess = (file: any) => {
  idcardNationalEmblemImgFile.value = file;
};
const licenseImgFile = ref();
const licenseImgSuccess = (file: any) => {
  licenseImgFile.value = file;
};
const goBackCompanyInfo = () => {
  stepActive.value = 1;
  ElNotification.closeAll();
};
const isBack = ref(false);
const goBack = () => {
  stepActive.value = 2;
  isBack.value = true;
  ElNotification.closeAll();
};
const onNext = async (formEl: FormInstance | undefined) => {
  if(props.hasCompanyOriginal){
    isBack.value = true;
  }else{
    isBack.value = false;
  }
  if (oldCompanyForm.value.legalPersonName == ruleForm.value.legalPersonName && oldCompanyForm.value.legalPersonIDCard == ruleForm.value.legalPersonIDCard && oldCompanyForm.value.companyName == ruleForm.value.companyName && oldCompanyForm.value.companyID == ruleForm.value.companyID && props.hasCompanyOriginal) {
    stepActive.value = 2;
    return;
  }
  if (!formEl) return;
  await formEl.validate(async (valid) => {
    if (valid) {
      dialog.loading = true;
      let middleData = JSON.parse(JSON.stringify(ruleForm.value));
      delete middleData.authenticationMethod;
      await useRegisterCollect.uploadCompanyInfoJson(ruleForm.value).then((res: any) => {
        if (res.code === "SUCCESS") {
          dialog.loading = false;
          ElMessage.success('企业信息修改成功');
          stepActive.value = 2;
          ElNotification.closeAll();
          ElNotification({
            title: '提示',
            message: '不上传身份证信息部分功能不能使用 !',
            type: 'warning',
          });
          emit('refresh');
        } else {
          dialog.loading = false;
          ElMessage.error('企业信息修改失败');
        }
      });
    }
  });
};
const onSubmit = () => {
  if (isBack.value) {
    if (!idcardImgFile.value || !idcardNationalEmblemImgFile.value || !licenseImgFile.value) {
      stepActive.value = 3;
      return;
    } else {
      uploadImage();
    }
  }
  uploadImage();
};
const uploadImage = async () => {
  if (!idcardImgFile.value) {
    ElMessage.error('请上传身份证正面照');
    return;
  }
  if (!idcardNationalEmblemImgFile.value) {
    ElMessage.error('请上传身份证国徽照');
    return;
  }
  if (!licenseImgFile.value) {
    ElMessage.error('请上传营业执照');
    return;
  }
  const loadingInstance = ElLoading.service({
    text: '正在上传',
    background: 'rgba(0,0,0,.2)',
  });
  let formData = new FormData();
  formData.append('idcard', idcardImgFile.value);
  formData.append('license', licenseImgFile.value);
  formData.append('idcardback', idcardNationalEmblemImgFile.value);
  await useRegisterCollect.uploadCompanyInfoFront(formData).then((res: any) => {
    if (res.code === "SUCCESS" && res.data.type == 'Company') {
      setCompanyInfoUpload(res.data.token, loadingInstance);
    }
  });
  return;
}
// 清除图片
const clearImg = () => {
  idcardImgFile.value = null;
  idcardNationalEmblemImgFile.value = null;
  licenseImgFile.value = null;
};
const setCompanyInfoUpload = async (companyToken: string, loadingInstance: any) => {
  await useRegisterApi()
      .uploadCompanyInfo({companyToken})
      .then((res: any) => {
        loadingInstance.close();
        if (res.code === "SUCCESS") {
          ElMessage.success('上传成功');
          stepActive.value = 3;
          ElNotification.closeAll();
          ElNotification({
            title: '提示',
            message: '不扫脸认证部分功能不能使用 !',
            type: 'warning',
          });
          emit('refresh');
        } else {
          ElMessage.error('上传失败');
        }
      });
};
const scanSuccess = (faceToken: string) => {
  dialog.loading = true;
  useRegisterApi()
      .uploadCompanyInfoFace({faceToken})
      .then((res: any) => {
        if (res.code === "SUCCESS") {
          dialog.loading = false;
          ElMessage.success('上传成功');
          dialog.isShowDialog = false;
          emit('refresh');
        } else {
          ElMessage.error('上传失败');
        }
      });
};
defineExpose({
  openDialog,
});
</script>

<style scoped lang="scss">
:deep(.el-form-item__content) {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  margin-bottom: 10px;
}

:deep(.el-tabs__active-bar) {
  height: 3px;
}

:deep(input::-webkit-input-placeholder) {
  color: #9ca5ba;
  font-size: 12px;
}

:deep(.el-input__inner) {
  background-color: transparent !important;
}

:deep(.el-tabs__item.is-active) {
  color: #121212;
  font-weight: bolder;
}

:deep(.el-step__icon-inner) {
  font-size: 16px !important;
}

:deep(.el-form-item__content) {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  margin-bottom: 10px;
}

.btns {
  display: flex;
  justify-content: space-between;
  width: 100%;
  margin-top: 10px;

  .btn {
    width: 100%;
    letter-spacing: 2px;
    font-weight: 300;
    margin-top: 10px;
  }
}

.login-content-form {
  .login-animation1 {
    width: 100%;
  }

  margin-top: 20px;
  @for $i from 1 through 4 {
    .login-animation#{$i} {
      opacity: 0;
      animation-name: error-num;
      animation-duration: 0.5s;
      animation-fill-mode: forwards;
      animation-delay: calc($i/10) + s;
    }
  }

  .login-content-code {
    width: 100%;
    padding: 0;
  }

  .login-content-submit {
    width: 100%;
    letter-spacing: 2px;
    font-weight: 300;
    margin-top: 15px;
  }

  .login-msg {
    color: var(--el-text-color-placeholder);
  }
}
</style>
