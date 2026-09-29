<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" height="500px" destroy-on-close @close="closeDialog">
		<div v-loading="dialog.loading">
			<el-steps :active="stepActive" align-center class="mb20">
				<el-step title="个人信息填写" />
				<el-step title="身份证上传" />
				<el-step title="个人认证" />
			</el-steps>
			<el-form v-if="stepActive == 1" size="large" :rules="rules" :model="ruleForm" ref="mobileFormRef">
				<el-form-item prop="userName">
					<el-input text placeholder="请输入您的真实姓名" v-model="ruleForm.userName" clearable autocomplete="off">
						<template #prefix>
							<SvgIcon name="my-people"></SvgIcon>
						</template>
					</el-input>
				</el-form-item>
				<el-form-item prop="userIDCard">
					<el-input text placeholder="请输入您的身份证号码" v-model="ruleForm.userIDCard" clearable autocomplete="off">
						<template #prefix>
							<SvgIcon name="my-idcard"></SvgIcon>
						</template>
					</el-input>
				</el-form-item>
				<el-form-item>
					<div class="btns">
						<el-button type="primary" v-waves @click="onNext(mobileFormRef)" block class="btn mt20" round> 绑定</el-button>
            <el-button v-if="hasUserOriginal" type="primary" v-waves @click="onNext(mobileFormRef)" block class="btn mt20" round> 下一步</el-button>
					</div>
				</el-form-item>
			</el-form>
			<el-form v-show="stepActive == 2" size="large" :rules="rules" :model="ruleForm">
				<el-form-item prop="phone">
					<imageUpload
						:upImgBoxCustomStyle="upImgBoxCustomStyle"
						@imgSuccess="imgFrontSuccess"
						imgText="支持jpg/png/bmp；文件大小不能超过2M"
						imgUpText="身份证正面照"
					/>
				</el-form-item>
				<el-form-item prop="phone">
					<imageUpload
						:upImgBoxCustomStyle="upImgBoxCustomStyle"
						@imgSuccess="imgReverseSideSuccess"
						imgText="支持jpg/png/bmp；文件大小不能超过2M"
						imgUpText="身份证反面照"
					/>
				</el-form-item>
				<el-form-item>
					<div class="btns">
						<el-button round type="primary" v-waves class="btn" @click="goBackUserInfo">
							<span>上一步</span>
						</el-button>
            <el-button v-if="isBack" type="primary" v-waves @click="onSubmit()" block class="btn" round> 下一步</el-button>
            <el-button v-else type="primary" v-waves @click="onSubmit()" block class="btn" round> 上传</el-button>
					</div>
				</el-form-item>
			</el-form>
			<scan-qr-codes
				v-if="stepActive == 3"
				current-user="UpdateUserInfo"
				:scan-qr-data="ruleForm"
				@go-back="goBack"
				@success="scanSuccess"
			></scan-qr-codes>
      <el-empty v-if="stepActive == 4" description="您已完成个人信息上传" :image="checkCircle">
        <!--  修改信息      -->
        <el-button type="primary" v-waves @click="editInfo" block class="btn" round> 修改信息</el-button>
      </el-empty>
		</div>
	</el-dialog>
</template>
<script setup lang="ts" name="registerUploadUserInfo">
import { reactive, ref } from 'vue';
import imageUpload from '/@/components/imageUpload/index.vue';
import { useRegisterApi } from '/@/api/register/index';
import { type FormRules, type FormInstance } from 'element-plus';
import type { uploadUserInfoJsontTypes } from '/@/api/register/types';
import { ElMessage, ElLoading } from 'element-plus';
import { ElNotification } from 'element-plus';
import ScanQrCodes from '/@/views/userCenter/countInfo/component/scanQrCodes.vue';
import checkCircle from '/@/assets/check-circle.png'

const useRegisterCollect = useRegisterApi();
const emit = defineEmits(['refresh']);

interface Props {
	userForm: uploadUserInfoJsontTypes;
	hasUserOriginal: boolean;
	hasUserInfo: boolean;
  hasFace: boolean;
}

const props = withDefaults(defineProps<Props>(), {});
// 定义变量内容
const oldUserInfo = ref<uploadUserInfoJsontTypes>(props.userForm);
const ruleForm = ref<uploadUserInfoJsontTypes>({
  authenticationMethod: 1,
  userName: '',
  userIDCard: '',
});
const stepActive = ref(1);
const upImgBoxCustomStyle = ref({
	width: '100%',
	height: '150px',
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	loading: false,
	title: '修改绑定的个人信息',
	submitTxt: '确认',
});
const openDialog = (form:any) => {
  if(!form) return;
  ruleForm.value = JSON.parse(JSON.stringify(form));
	dialog.isShowDialog = true;
  if (props.hasUserInfo && props.hasUserOriginal&&props.hasFace) {
    stepActive.value = 4;
    ElNotification({
      title: '提示',
      message: '您已完成个人信息上传 !',
      type: 'warning',
    });
  }else if(props.hasUserInfo && props.hasUserOriginal){
    stepActive.value = 3;
    ElNotification({
      title: '提示',
      message: '不上传身份证信息和不扫脸认证部分功能不能使用 !',
      type: 'warning',
    });
  }else if (props.hasUserInfo) {
		stepActive.value = 2;
		ElNotification({
			title: '提示',
			message: '不上传身份证信息部分功能不能使用 !',
			type: 'warning',
		});
	} else if (props.hasUserOriginal) {
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
	// 清除ElNotification
	ElNotification.closeAll();
  clearImg();
};
const mobileFormRef = ref();
const rules = reactive<FormRules>({
	userName: [
		{
			required: true,
			message: '请输入您的真实姓名',
			trigger: 'blur',
		},
	],
	userIDCard: [
		{
			required: true,
			message: '请输入您的身份证号码',
			trigger: 'blur',
		},
	],
});
const imgFrontFile = ref();
const imgReverseFile = ref();
//图片组件传过来的文件数据
const imgFrontSuccess = async (file: any) => {
	imgFrontFile.value = file;
};
const imgReverseSideSuccess = async (file: any) => {
	imgReverseFile.value = file;
};
const onNext = async (formEl: FormInstance | undefined) => {
  if(props.hasUserOriginal){
    isBack.value = true;
  }else{
    isBack.value = false;
  }
  if (oldUserInfo.value.userName == ruleForm.value.userName && oldUserInfo.value.userIDCard == ruleForm.value.userIDCard && props.hasUserOriginal) {
    stepActive.value = 2;
    return;
  }
	if (!formEl) return;
	await formEl.validate(async (valid) => {
		if (valid) {
			dialog.loading = true;
			let middleData = JSON.parse(JSON.stringify(ruleForm.value));
			delete middleData.authenticationMethod;
			await useRegisterCollect.uploadUserInfoJson(ruleForm.value).then((res: any) => {
				if (res.code === "SUCCESS") {
					dialog.loading = false;
					ElMessage.success('身份信息添加成功');
					ElNotification.closeAll();
					ElNotification({
						title: '提示',
						message: '不上传身份证信息部分功能不能使用 !',
						type: 'warning',
					});
          emit('refresh');
					stepActive.value = 2;
				}
			});
		}
	});
};
const onSubmit = () => {
  if(isBack.value){
    if(!imgFrontFile.value||!imgReverseFile.value){
      stepActive.value = 3;
      return;
    }else {
      uploadImage();
    }
  }
  uploadImage();
};
const uploadImage =async ()=>{
  if (!imgFrontFile.value) {
    ElMessage.warning('请上传身份证正面照');
    return;
  }
  if (!imgReverseFile.value) {
    ElMessage.warning('请上传身份证反面照');
    return;
  }
  const loadingInstance = ElLoading.service({
    text: '正在上传',
    background: 'rgba(0,0,0,.2)',
  });
  let formData = new FormData();
  formData.append('idcard', imgFrontFile.value);
  formData.append('idcardback', imgReverseFile.value);
  await useRegisterCollect.uploadUserInfoBack(formData).then((res: any) => {
    if (res.code === "SUCCESS" && res.data.type == 'IDCard') {
      setUserInfoUpload(res.data.token, loadingInstance);
    }
  });
  // useRegisterCollect.userinfoUpload(formData)
}
// 图片清除
const clearImg = () => {
  imgFrontFile.value = null;
  imgReverseFile.value = null;
};
//
const setUserInfoUpload = async (idcardToken: string, loadingInstance: any) => {
	await useRegisterCollect.userinfoUpload({ idcardToken }).then((res: any) => {
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
const goBackUserInfo = () => {
	stepActive.value = 1;
  ElNotification.closeAll();
};
const isBack = ref(false);
const goBack = () => {
	stepActive.value = 2;
  // 返回上一步状态
  isBack.value = true;
  ElNotification.closeAll();
};
const scanSuccess = (faceToken: string) => {
	dialog.loading = true;
	useRegisterApi()
		.uploadUserInfoFace({ faceToken })
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
</style>
