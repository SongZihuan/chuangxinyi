<template>
	<div>
		<div class="personalInfoForm-content">
			<el-form size="large" class="personalInfoForm-content-form" :rules="rules" :model="ruleForm" ref="personalInfoFormRef">
				<el-form-item class="personalInfoForm-animation1 mb20" prop="id">
					<el-input text placeholder="请输入注册的个人账户名称、邮箱或手机号" v-model="ruleForm.id" clearable autocomplete="off">
						<template #prefix>
              <el-icon><Phone /></el-icon>
						</template>
					</el-input>
				</el-form-item>
				<el-form-item class="personalInfoForm-animation1 mb20" prop="userName">
					<el-input text placeholder="请输入注册用户名字" v-model="ruleForm.userName" clearable autocomplete="off">
						<template #prefix>
              <el-icon><User /></el-icon>
						</template>
					</el-input>
				</el-form-item>
				<el-form-item class="personalInfoForm-animation1 mb20" prop="userIdCard">
					<el-input text placeholder="注册用户身份证号码" v-model="ruleForm.userIdCard" clearable autocomplete="off">
						<template #prefix>
              <el-icon><CreditCard /></el-icon>
						</template>
					</el-input>
				</el-form-item>
				<el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
					<el-form-item label="类别">
						<el-radio-group v-model="ruleForm.isCompany">
							<el-radio :label="false">个人</el-radio>
							<el-radio :label="true">企业</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
				<template v-if="ruleForm.isCompany">
					<el-form-item class="businessInfoForm-animation1 mb20" prop="companyID">
						<el-input text placeholder="请输入注册企业统社会信用代码" v-model="ruleForm.companyID" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="ele-Stamp"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
					<el-form-item class="businessInfoForm-animation1 mb20" prop="legalPersonName">
						<el-input text placeholder="请输入法人名称" v-model="ruleForm.legalPersonName" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="ele-Avatar"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
					<el-form-item class="businessInfoForm-animation1 mb20" prop="companyName">
						<el-input text placeholder="请输入注册企业名称" v-model="ruleForm.companyName" clearable autocomplete="off">
							<template #prefix>
								<SvgIcon name="ele-OfficeBuilding"></SvgIcon>
							</template>
						</el-input>
					</el-form-item>
				</template>
			</el-form>
		</div>
		<div class="personalInfoForm-animation3 personalInfoForm-bnts">
			<el-button round type="primary" v-waves class="personalInfoForm-content-next" @click="goBack">
				<span>上一步</span>
			</el-button>
			<el-button round type="primary" v-waves class="personalInfoForm-content-next" @click="onNext(personalInfoFormRef)">
				<span>下一步</span>
			</el-button>
		</div>
	</div>
  <div class="register">
    <div>想起密码?</div>
    <div @click="goLogin">立即登录</div>
  </div>
</template>
<script setup lang="ts">
import { reactive, ref } from 'vue';
import { personInfoFormTypes } from '/@/api/register/types';
import { FormInstance, FormRules } from 'element-plus';
import {CreditCard, Phone} from "@element-plus/icons-vue";
import {useRoute, useRouter} from "vue-router";

const emit = defineEmits(['resetBackForm', 'goBack']);
interface Props {
	ruleFormPersonal: personInfoFormTypes;
}
const props = withDefaults(defineProps<Props>(), {});
const ruleForm = ref<personInfoFormTypes>(props.ruleFormPersonal);
const personalInfoFormRef = ref();
const rules = reactive<FormRules>({
	id: [{ required: true, trigger: 'blur', message: '请输入注册的个人账户名称、邮箱或手机号' }],
	userName: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
	userIdCard: [{ required: true, message: '请输入用户身份证号码', trigger: 'blur' }],
});

const router = useRouter()
const route = useRoute()

const goLogin = () => {
  router.push({ path: '/login', query: route.query});
};

const goBack = () => {
	emit('goBack', 1);
};
const onNext = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			emit('resetBackForm', ruleForm.value);
		}
	});
};
</script>
<style scoped lang="scss">
.personalInfoForm-bnts {
	display: flex;
	justify-content: space-between;
	.personalInfoForm-content-next {
		width: 100%;
		letter-spacing: 2px;
		font-weight: 300;
		margin-top: 10px;
	}
}
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
.personalInfoForm-content {
	height: 300px;
	overflow-y: auto;
}

/* 隐藏标准的滚动条 */
.personalInfoForm-content::-webkit-scrollbar {
	width: 0;
}

.personalInfoForm-content::-webkit-scrollbar {
	width: 0;
}

/* 隐藏 IE 和 Edge 浏览器的滚动条 */
::-ms-scrollbar {
	width: 0;
}
</style>
