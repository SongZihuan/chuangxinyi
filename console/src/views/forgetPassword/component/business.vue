<template>
	<div>
		<div class="business-content">
			<el-form size="large" class="business-content-form" :rules="rules" :model="ruleForm" ref="businessFormRef">
				<el-form-item class="business-animation1" prop="phone">
					<el-input text placeholder="请输入企业账户或手机号" v-model="ruleForm.phone" clearable autocomplete="off">
						<template #prefix>
							<el-icon><Phone /></el-icon>
						</template>
					</el-input>
				</el-form-item>
				<el-form-item class="business-animation3 business-sider">
					<el-col :span="24">
						<slider-silence @siderEmit="siderEmit" ref="sliderRef" />
					</el-col>
				</el-form-item>
				<el-form-item class="business-animation2" prop="code" v-if="isSliderCheck">
					<el-col :span="15">
						<el-input text maxlength="6" placeholder="请输入验证码" v-model="ruleForm.code" clearable autocomplete="off">
							<template #prefix>
								<el-icon class="el-input__icon">
									<ele-Position />
								</el-icon>
							</template>
						</el-input>
					</el-col>
					<el-col :span="1"></el-col>
					<el-col :span="8">
						<el-button v-waves class="business-content-code" @click="getCode" :disabled="codeState.isSend">
							{{ codeState.codeName }}
						</el-button>
					</el-col>
				</el-form-item>
			</el-form>
		</div>
		<div class="business-animation3">
			<el-button round type="primary" v-waves class="business-content-next" @click="onNext(businessFormRef)">
				<span>下一步</span>
			</el-button>
		</div>
    <div class="register">
      <div>想起密码?</div>
      <div @click="goLogin">立即登录</div>
    </div>
	</div>
</template>
<script setup lang="ts">
// 定义变量内容
import { reactive, ref } from 'vue';
import { ElMessage, FormInstance, FormRules } from 'element-plus';
import { verifyPhone } from '/@/utils/toolsValidate';
import sliderSilence from '/@/components/Slider/index.vue';
import { checkPhoneCodeTypes, codeTypes } from '/@/api/register/types';
import { useBaseApi } from '/@/api/base';
import type { sliderHeadersTypes } from '/@/api/base/types';
import { useLoginApi } from '/@/api/login';
import { useUserInfo } from '/@/stores/userInfo';
import { Phone } from '@element-plus/icons-vue';
import {useRoute, useRouter} from "vue-router";

const route = useRoute()
const emit = defineEmits(['resetBack']);
const sliderRef = ref();
const ruleForm = ref<checkPhoneCodeTypes>({
	phone: '',
	code: '',
	type: 'PhoneCheck',
});
const businessFormRef = ref();
const rules = reactive<FormRules>({
	phone: [{ trigger: 'blur', validator: verifyPhone }],
	code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
});
const isSliderCheck = ref<boolean>(false);
const siderEmit = (data: any) => {
	sliderHeaders.value = data;
	isSliderCheck.value = true;
};
//阿里云滑块返回的参数
let sliderHeaders = ref<sliderHeadersTypes>({
	sig: '',
	sessionId: '',
	token: '',
});
const useBaseApiCollect = useBaseApi();
const useLoginrCollect = useLoginApi();
const stores = useUserInfo();
const isGetCode = ref<boolean>(false);
const codeState = reactive({
	isSend: false,
	codeName: '获取验证码',
	totalTime: 60, //一般是60
	timer: undefined, //定时器
}) as codeTypes;
const getCode = async () => {
	await useBaseApiCollect.sendPhoneCode({ phone: ruleForm.value.phone }, sliderHeaders.value).then((res: any) => {
		if (res.code === "SUCCESS") {
			isGetCode.value = true;

			ElMessage.success('发送手机验证码成功');
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
		}else{
      sliderRef.value.resetSider();
    }
	}).catch(()=>{
    sliderRef.value.resetSider();
  });
};
const onNext = (formEl: FormInstance | undefined) => {
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
						stores.setUserType({ type: res.data.type, subType: res.data.subType });
						emit('resetBack', res.data.token);
					} else {
						ElMessage.warning('该账号未注册!');
					}
				}
			});
		}
	});
};

const router = useRouter()

const goLogin = () => {
  router.push({ path: '/login', query: route.query});
};

</script>
<style scoped lang="scss">
.register {
  display: flex;
  justify-content: center;
  font-size: 14px;
  margin-top: 20px;
  cursor: pointer;
  position: relative;
  z-index: 2004;

  :nth-child(1) {
    color: #9ca5ba;
  }

  :nth-child(2) {
    color: var(--el-color-primary);
    margin-left: 6px;
    text-decoration: underline;
  }
}
:deep(.nc-container #nc_1_wrapper) {
	width: 100% !important;
}
.business-content {
	height: 260px;
}
.business-content-form {
	width: 100%;
	margin-top: 20px;
	@for $i from 1 through 4 {
		.business-animation#{$i} {
			opacity: 0;
			animation-name: error-num;
			animation-duration: 0.5s;
			animation-fill-mode: forwards;
			animation-delay: calc($i/10) + s;
		}
	}

	.business-content-code {
		width: 100%;
		padding: 0;
	}

	.business-sider {
		margin-top: 25px;
	}
}
.business-content-next {
	width: 100%;
	letter-spacing: 2px;
	font-weight: 300;
	margin-top: 10px;
}
</style>
