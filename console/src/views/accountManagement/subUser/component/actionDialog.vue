<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="440px" destroy-on-close @close="closeDialog">
    <el-form ref="formRef" :model="ruleForm" size="default" :rules="rules">
      <el-row>
        <el-col :span="24" class="mb10">
          <el-form-item prop="type">
            <el-radio-group v-model="ruleForm.newWallet">
              <el-radio :label="true">新建钱包</el-radio>
              <el-radio :label="false">钱包共享</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="24" class="mb15">
          <el-form-item prop="phone">
            <el-input text placeholder="请输入手机号" v-model="ruleForm.phone" clearable autocomplete="off">
              <template #prefix>
                <el-icon>
                  <ele-Phone/>
                </el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row class="mb10">
        <el-form-item prop="code" v-if="isSliderCheck" class="code-box mb15">
          <el-col :span="15">
            <el-input text maxlength="6" placeholder="请输入验证码" v-model="ruleForm.code" clearable
                      autocomplete="off">
              <template #prefix>
                <el-icon class="el-input__icon">
                  <ele-Position/>
                </el-icon>
              </template>
            </el-input>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="8">
            <el-button v-waves @click="getCode" :disabled="codeState.isSend" class="code"> {{
                codeState.codeName
              }}
            </el-button>
          </el-col>
        </el-form-item>
      </el-row>
      <el-col :span="24" class="mb10">
        <el-form-item>
          <slider-silence @siderEmit="siderEmit" ref="sliderRef"/>
        </el-form-item>
      </el-col>
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
import {onMounted, reactive, ref} from 'vue';
import {type FormRules, type FormInstance} from 'element-plus';
import {ElMessage} from 'element-plus';
import {verifyPhone} from '/@/utils/toolsValidate';
import {sliderHeadersTypes} from '/@/api/base/types';
import sliderSilence from '/@/components/Slider/index.vue';
import {codeTypes} from '/@/api/register/types';
import {checkPhoneTypes} from '/@/views/accountManagement/subUser/types';
import {useBaseApi} from '/@/api/base';
import {useLoginApi} from '/@/api/login';
import {registerSubUserTypes} from '/@/api/subUser/types';
import {useSubUserCenterApi} from '/@/api/subUser';

const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<checkPhoneTypes>({
  phone: '',
  code: '',
  newWallet: true,
  type: 'PhoneCheck',
});
const rules = reactive<FormRules>({
  phone: [{required: true, trigger: 'blur', validator: verifyPhone}],
  code: [{required: true, message: '请输入验证码', trigger: 'blur'}],
});
let registerSubUserData = ref<registerSubUserTypes>({
  phoneToken: '',
});
// let userData:any = Session.get('userData');
const sliderRef = ref();
//阿里云滑块返回的参数
let sliderHeaders = ref<sliderHeadersTypes>({
  sig: '',
  sessionId: '',
  token: '',
});
const isSliderCheck = ref<boolean>(false);
const siderEmit = (data: any) => {
  sliderHeaders.value = data;
  isSliderCheck.value = true;
};
const useBaseApiCollect = useBaseApi();
const useLoginrCollect = useLoginApi();
const isGetCode = ref<boolean>(false);
const codeState = reactive({
  isSend: false,
  codeName: '获取验证码',
  totalTime: 60, //一般是60
  timer: undefined, //定时器
}) as codeTypes;
const getCode = async () => {
  await useBaseApiCollect.sendPhoneCode({phone: ruleForm.value.phone}, sliderHeaders.value).then((res: any) => {
    if (res.code === "SUCCESS") {
      isGetCode.value = true;
      if (codeState.isSend) return;
      // getCode() // 获取验证码的接口
      codeState.isSend = true;
      codeState.codeName = codeState.totalTime + 's后重新发送';
      codeState.timer = setInterval(() => {
        codeState.totalTime--;
        codeState.codeName = codeState.totalTime + 's后重新发送';
        if (codeState.totalTime < 0) {
          clearInterval(codeState.timer);
          codeState.codeName = '重新发送';
          codeState.totalTime = 60;
          codeState.isSend = false;
          sliderRef.value.resetSider();
        }
      }, 1000);

      ElMessage.success('发送手机验证码成功');
    }else{
      sliderRef.value.resetSider();
    }
  }).catch(()=>{
    sliderRef.value.resetSider();
  });
};
const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: 'add',
  title: '子账号新增',
  submitTxt: '新增',
});
const resetForm = () => {
  ruleForm.value = {
    phone: '',
    code: '',
    newWallet: true,
    type: 'PhoneCheck',
  };
};
const openDialog = () => {
  resetForm()
  dialog.isShowDialog = true;
};
const closeDialog = () => {
  dialog.isShowDialog = false;
};
// 新增子账号
const onSubmit = (formEl: FormInstance | undefined) => {
  if (!isSliderCheck.value) {
    ElMessage.warning('请滑动滑块移动最右侧!');
    return;
  }

  if (!formEl) return;
  formEl.validate(async (valid) => {
    if (valid) {
      await useLoginrCollect.phoneLogin(ruleForm.value, sliderHeaders.value).then((res: any) => {
        isSliderCheck.value = false;
        sliderRef.value.resetSider();
        if (res.code === "SUCCESS") {
          if (res.data.type == 'PhoneCheck') {
            registerSubUserData.value.phoneToken = res.data.token;
            useSubUserCenterApi()
                .registerSubUser(registerSubUserData.value)
                .then((res: any) => {
                  if (res.code === "SUCCESS") {
                    ElMessage.success('子账号新增成功');
                    closeDialog();
                    emit('refresh');
                  }
                });
          } else {
            ElMessage.warning('该账号未注册!');
          }
        }
      });
    }
  });
};
onMounted(() => {
});
defineExpose({
  openDialog,
  closeDialog,
});
</script>

<style lang="scss" scoped>
.code-box {
  width: 500px;

  .code {
    width: 100%;
  }
}
</style>
