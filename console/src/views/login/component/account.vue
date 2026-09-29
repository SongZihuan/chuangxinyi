<template>
  <el-form size="large" class="login-content-form" ref="formRef" :rules="rules" :model="state.ruleForm">
    <el-form-item class="login-animation1 mb20" prop="userID">
      <el-input text placeholder="手机号、邮箱、用户名" v-model="state.ruleForm.userID" clearable autocomplete="off">
        <template #prefix>
          <el-icon class="el-input__icon"><ele-User /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item class="login-animation2" prop="passwordHash">
      <el-input :type="state.isShowPassword ? 'text' : 'password'" placeholder="密码" v-model="state.ruleForm.password" autocomplete="off">
        <template #prefix>
          <el-icon class="el-input__icon"><ele-Unlock /></el-icon>
        </template>
        <template #suffix>
          <i
              class="iconfont el-input__icon login-content-password"
              :class="state.isShowPassword ? 'icon-yincangmima' : 'icon-xianshimima'"
              @click="state.isShowPassword = !state.isShowPassword"
          >
          </i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item class="login-animation3">
      <el-col :span="24">
        <slider @siderEmit="siderEmit" ref="sliderRef" />
      </el-col>
    </el-form-item>
    <!-- <el-form-item>
      <el-col :span="24">
        <silenceSider @siderEmit="siderEmit" ref="sliderRef" @nvcValEmit="nvcValEmit" />
      </el-col>
    </el-form-item> -->
    <div class="forgetPassword login-animation2" @click="goPassword">
      <div>忘记密码?</div>
    </div>
    <el-form-item class="login-animation4">
      <el-button type="primary" class="login-content-submit" round v-waves @click="onLogin(formRef)" :loading="state.loading.signIn">
        <span>登 录</span>
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts" name="loginAccount">
import { reactive, ref, onMounted, onBeforeUnmount } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import sha256 from 'sha256';
import { useLoginApi } from '/@/api/login/index';
import type { sliderHeadersTypes } from '/@/api/base/types';
import { ElMessage } from 'element-plus';
import { useBaseApi } from '/@/api/base/index';
import { useRouter, useRoute } from 'vue-router';
import { useUserInfo } from '/@/stores/userInfo';
import slider from '/@/components/Slider/index.vue';
const router = useRouter();
const route = useRoute();
interface Props {
  isAgree: boolean; // 回显图片地址
}
const stores = useUserInfo();
const salt = ref('');
const props = withDefaults(defineProps<Props>(), {
  isAgree: false,
});
// 定义变量内容
const formRef = ref();
const emit = defineEmits(['signBack', 'totpBack']);
const useLoginrCollect = useLoginApi();
const state = reactive({
  isShowPassword: false,
  ruleForm: {
    userID: '',
    password: '',
    type: 'UserToken',
  },
  loading: {
    signIn: false,
  },
});
const rules = reactive<FormRules>({
  userID: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
});
//阿里云滑块返回的参数
let sliderHeaders = ref<sliderHeadersTypes>({
  sig: '',
  sessionId: '',
  token: '',
});
const sliderRef = ref();
const isSliderCheck = ref<boolean>(false);
const siderEmit = (data: any) => {
  sliderHeaders.value = data;
  isSliderCheck.value = true;
};
const onLogin = (formEl: FormInstance | undefined) => {
  if (!isSliderCheck.value) {
    ElMessage.warning('请滑动滑块移动最右侧!');
    return;
  }

  if (!props.isAgree) {
    ElMessage.warning('请勾选用户协议');
    return;
  }
  if (!formEl) return;
  formEl.validate(async (valid) => {
    if (valid) {
      await useLoginrCollect
          .accountLogin({ ...state.ruleForm, passwordHash: sha256(salt.value + ':' + state.ruleForm.password), headers: sliderHeaders.value })
          .then(async (res: any) => {
            isSliderCheck.value = false;
            sliderRef.value.resetSider();
            if (res.code === 'SUCCESS') {
              await stores.setUserType({ type: res.data.type, subType: res.data.subType });
              if (res.data.type === 'Login2FA') {
                emit('totpBack', res.data.token);
              } else {
                emit('signBack', res.data.token, res.data.type, res.data.subType);
              }
            }
          });
    }
  });
};
const siderCheck = () => {
  sliderRef.value.siderCheck();
};

const getslat = () => {
  useBaseApi()
      .salt()
      .then((res: any) => {
        if (res.code === 'SUCCESS') {
          salt.value = res.data.salt;
        }
      });
};
const goPassword = () => {
	router.push({ path: '/forgetPassword', query: route.query});
};
onMounted(() => {
  getslat();
});
onBeforeUnmount(() => {});
defineExpose({
  siderCheck,
});
</script>

<style scoped lang="scss">
.forgetPassword {
  margin-top: 10px;
  text-align: right;
  color: var(--el-color-primary);
}
.login-content-form {
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

  .login-content-password {
    display: inline-block;
    width: 20px;
    cursor: pointer;
    &:hover {
      color: #909399;
    }
  }
  .login-content-code {
    width: 100%;
    padding: 0;
    font-weight: bold;
    letter-spacing: 5px;
  }
  .login-content-submit {
    width: 100%;
    letter-spacing: 2px;
    font-weight: 300;
    margin-top: 15px;
  }
}
</style>
