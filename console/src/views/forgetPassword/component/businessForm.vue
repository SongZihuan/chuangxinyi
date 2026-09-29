<template>
	<div>
		<div class="businessInfoForm-content">
			<el-form size="large" class="businessInfoForm-content-form" :rules="rules" :model="ruleForm" ref="businessInfoFormRef">
				<el-form-item class="businessInfoForm-animation1 mb20" prop="id">
					<el-input text placeholder="请输入注册的企业账户、邮箱或手机号" v-model="ruleForm.id" clearable autocomplete="off">
						<template #prefix>
							<el-icon><Phone /></el-icon>
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
				<el-form-item class="businessInfoForm-animation1 mb20" prop="legalPersonID">
					<el-input text placeholder="注册法人身份证号码" v-model="ruleForm.legalPersonID" clearable autocomplete="off">
						<template #prefix>
							<el-icon><CreditCard /></el-icon>
						</template>
					</el-input>
				</el-form-item>
			</el-form>
		</div>
		<div class="businessInfoForm-animation3 businessInfoForm-bnts">
			<el-button round type="primary" v-waves class="businessInfoForm-content-next" @click="goBack">
				<span>上一步</span>
			</el-button>
			<el-button round type="primary" v-waves class="businessInfoForm-content-next" @click="onSumbit(businessInfoFormRef)">
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
import { reactive, ref } from 'vue';
import { businessInfoFormTypes } from '/@/api/register/types';
import { FormInstance, FormRules } from 'element-plus';
import { CreditCard, Phone } from '@element-plus/icons-vue';
import {useRoute, useRouter} from "vue-router";

const emit = defineEmits(['resetBackForm', 'goBack']);

interface Props {
	ruleFormBusinese: businessInfoFormTypes;
}

const router = useRouter()
const route = useRoute()

const goLogin = () => {
  router.push({ path: '/login', query: route.query});
};

const props = withDefaults(defineProps<Props>(), {});
const ruleForm = ref<businessInfoFormTypes>(props.ruleFormBusinese);
const businessInfoFormRef = ref();
const rules = reactive<FormRules>({
	id: [{ required: true, trigger: 'blur', message: '请输入注册的企业账户、邮箱或手机号' }],
	companyName: [{ required: true, message: '请输入注册企业名称', trigger: 'blur' }],
	companyID: [{ required: true, message: '请输入注册企业统社会信用代码', trigger: 'blur' }],
	legalPersonName: [{ required: true, message: '请输入法人名称', trigger: 'blur' }],
	legalPersonID: [{ required: true, message: '请输入法人身份证号码', trigger: 'blur' }],
});
const goBack = () => {
	emit('goBack', 1);
};
const onSumbit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			emit('resetBackForm', ruleForm.value);
		}
	});
};
</script>
<style scoped lang="scss">
.businessInfoForm-content {
	height: 300px;
	overflow-y: auto;
}
.businessInfoForm-bnts {
	display: flex;
	justify-content: space-between;
	.businessInfoForm-content-next {
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
/* 隐藏标准的滚动条 */
.businessInfoForm-content::-webkit-scrollbar {
	width: 0;
}

.businessInfoForm-content::-webkit-scrollbar {
	width: 0;
}

/* 隐藏 IE 和 Edge 浏览器的滚动条 */
::-ms-scrollbar {
	width: 0;
}
</style>
