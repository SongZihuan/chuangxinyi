<template>
	<div class="scan-qr-codes-container">
		<div class="scan-qr-codes-box">
			<div class="sacn-qr-codes-img">
				<div ref="qrcodeRef"></div>
				<div class="success" v-if="codeStatus == 2">
					<div class="success-img"></div>
					<div class="success-title">验证成功</div>
				</div>
				<div class="error" v-if="codeStatus == 3">
					<div class="error-img"></div>
					<div class="error-title">扫码超时</div>
					<div class="refresh" @click="refreshCode"><SvgIcon name="ele-RefreshRight"></SvgIcon>重新扫码</div>
				</div>
			</div>
			<div class="hint">
				<span>请使用支付宝扫码</span>
				<span>扫码后请在支付宝中完成人脸识别</span>
			</div>
		</div>
		<div class="scan-qr-codes-btns">
			<el-button round type="primary" v-waves class="personalInfoForm-content-next" @click="goBack">
				<span>上一步</span>
			</el-button>
			<el-button v-if="codeStatus == 2" round type="primary" v-waves class="personalInfoForm-content-next" @click="onSumbit">
				<span>进入系统</span>
			</el-button>
		</div>
	</div>
  <div class="login-action">
    <el-checkbox v-model="isAgree">我已阅读并同意 </el-checkbox>
    <div @click="clickAgreement" class="admin-agreement">用户协议</div>
  </div>
  <div class="register">
    <div>想起密码?</div>
    <div @click="goLogin">立即登录</div>
  </div>
  <agreementDialog ref="dialogRef" />
</template>
<script setup lang="ts">
import {ref, nextTick, onMounted, onBeforeUnmount, defineAsyncComponent} from 'vue';
import QRCode from 'qrcodejs2-fixes';
import { useLoginApi } from '/@/api/login';
import { checkFaceRecognitionTypes, startFaceRecognitionTypes } from '/@/api/login/types';
import {businessInfoFormTypes, personInfoFormTypes, uploadUserInfoJsontTypes} from '/@/api/register/types';
import { ElMessage } from 'element-plus';
import {useRoute, useRouter} from "vue-router";

const agreementDialog = defineAsyncComponent(() => import('/@/views/login/component/agreementDialog.vue')); //管理员协议组件
const emit = defineEmits(['goBack', 'onSumbit', 'success']);

const router = useRouter()
const route = useRoute()

const goLogin = () => {
  router.push({ path: '/login', query: route.query});
};

const isAgree = ref<boolean>(false);
const dialogRef = ref();
const clickAgreement = () => {
  dialogRef.value.openDialog();
};

interface Props {
	scanQrData: personInfoFormTypes | businessInfoFormTypes;
	currentUser: string;
}

const props = withDefaults(defineProps<Props>(), {
	currentUser: '',
});
const faceToken = ref<string>('');
const certifyID = ref<string>('');
const qrcodeRef = ref<HTMLElement | null>(null);
const codeUrl = ref<string>('');
// 有效时间 15分钟
const validTime = ref<number>(15 * 60 * 1000);
const codeStatus = ref<number>(0); // 0:未扫码 1:已扫码 2:扫码成功 3:扫码失败
const generateORCode = async () => {
	validTime.value = 15 * 60 * 1000;
	const data: startFaceRecognitionTypes = {
		name: '',
		id: '',
	};
	if (props.currentUser === 'personal') {
		const scanQrData = props.scanQrData as personInfoFormTypes;
		data.name = scanQrData.userName;
		data.id = scanQrData.userIdCard;
	}else if(props.currentUser === 'UpdateUserInfo'){
    const scanQrData = props.scanQrData as uploadUserInfoJsontTypes;
    data.name = scanQrData.userName;
    data.id = scanQrData.userIDCard;
  } else {
		const scanQrData = props.scanQrData as businessInfoFormTypes;
		data.name = scanQrData.legalPersonName;
		data.id = scanQrData.legalPersonID;
	}
	if (!data.name) {
		ElMessage.warning('用户名获取失败');
		return;
	}
	if (!data.id) {
		ElMessage.warning('身份证号码获取失败');
		return;
	}
	await useLoginApi()
		.startFaceRecognition(data)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				certifyID.value = res.data.certifyID;
				getCodeUrl(certifyID.value);
			}
		});
};
const goBack = () => {
	emit('goBack', 2);
};
const onSumbit = () => {
  if (!isAgree.value) {
    ElMessage.warning('请先同意用户协议');
    return;
  }
	if (codeStatus.value != 2) {
		ElMessage.warning('请先扫码');
		return;
	}
	if (!faceToken.value) {
		codeStatus.value = 3;
		ElMessage.warning('人脸识别失败');
		return;
	}
	emit('success', faceToken.value);
};

// 获取验证码url
const getCodeUrl = async (certifyid: string) => {
	let url = import.meta.env.VITE_FACE_CALLBACK+`?certifyid=${certifyid}`;
	codeUrl.value = url;
	initQrcode();
};
// 生成二维码
const initQrcode = () => {
	nextTick(() => {
		(<HTMLElement>qrcodeRef.value).innerHTML = '';
		new QRCode(qrcodeRef.value, {
			text: codeUrl.value,
			width: 200,
			height: 200,
			colorDark: '#000000',
			colorLight: '#ffffff',
		});
	});
	timer.value = setInterval(() => {
		validTime.value -= 3000;
		if (validTime.value <= 0) {
			clearInterval(timer.value);
			codeStatus.value = 3;
		}
		checkQRcodeStatus(certifyID.value);
	}, 3000);
};
// 定时请求
const timer = ref<any>(null);
// 轮询二维码状态
const checkQRcodeStatus = async (certifyid: string) => {
	let data: checkFaceRecognitionTypes = {
		certifyID: certifyid,
	};
	await useLoginApi()
		.checkFaceRecognition(data)
		.then((res: any) => {
			if (res && res.code === "SUCCESS") {
				clearInterval(timer.value);
				if (res.data.type == 'Face') {
					faceToken.value = res.data.token;
				}
				codeStatus.value = 2;
			}
		});
};
onBeforeUnmount(() => {
	clearInterval(timer.value);
});
// 刷新二维码
const refreshCode = () => {
	codeStatus.value = 1;
	generateORCode();
};
//
onMounted(() => {
	generateORCode();
});
</script>

<style scoped lang="scss">
.scan-qr-codes-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	width: 100%;
	height: 100%;
	.sacn-qr-codes-img {
		position: relative;
		width: 210px;
		height: 200px;
		margin-bottom: 10px;
		.success,
		.error {
			position: absolute;
			left: 0;
			top: -2px;
			display: block;
			width: 100%;
			height: 105%;
			background: rgba(255, 255, 255, 0.98);
			background-size: 100% 100%;
			.success-img {
				width: 100px;
				height: 100px;
				margin: 40px auto 0;
				background: url('/@/assets/check-circle.png') no-repeat center;
				background-size: 100% 100%;
			}
			.refresh {
				display: flex;
				justify-content: center;
				align-items: center;
				font-size: 14px;
				color: #999999;
				margin-top: 10px;
				cursor: pointer;
				user-select: none;
			}
			.error-img {
				width: 100px;
				height: 100px;
				margin: 40px auto 0;
				background: url('/@/assets/x-circle.png') no-repeat center;
				background-size: 100% 100%;
			}
			.success-title,
			.error-title {
				width: 100%;
				text-align: center;
				font-size: 16px;
				color: #333333;
				margin-top: 10px;
			}
		}
	}
	.scan-qr-codes-box {
		max-width: 300px;
		max-height: 300px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		-webkit-box-pack: center;
		padding: 30px;
		margin-top: 6px;
		border: 4px solid rgb(216, 216, 216);
		.hint {
			display: flex;
			flex-direction: column;
			align-items: center;
			margin-top: 10px;
			font-size: 14px;
			color: #999999;
		}
	}
	.scan-qr-codes-btns {
		display: flex;
		justify-content: space-between;
		width: 100%;
		margin-top: 10px;
		.personalInfoForm-content-next {
			width: 100%;
			letter-spacing: 2px;
			font-weight: 300;
			margin-top: 10px;
		}
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
.login-action {
  display: flex;
  flex-direction: row;
  align-items: center;
  position: relative;
  z-index: 2002;
  .admin-agreement {
    color: var(--el-color-primary);
    margin-left: 4px;
    margin-right: 30px;
    text-decoration: underline;
    cursor: pointer;
    margin-top: -2px;
  }
}
</style>
