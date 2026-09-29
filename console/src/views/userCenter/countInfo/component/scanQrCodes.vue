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
				<span>完成认证</span>
			</el-button>
		</div>
	</div>
</template>
<script setup lang="ts">
import { ref, nextTick, onMounted, onBeforeUnmount } from 'vue';
import QRCode from 'qrcodejs2-fixes';
import { useLoginApi } from '/@/api/login';
import { checkFaceRecognitionTypes, startFaceRecognitionTypes } from '/@/api/login/types';
import { uploadCompanyInfoJsonTypes, uploadUserInfoJsontTypes } from '/@/api/register/types';
import { ElMessage } from 'element-plus';
import {useUserInfo} from "/@/stores/userInfo";
const emit = defineEmits(['goBack', 'onSumbit', 'success']);

interface Props {
	scanQrData: uploadUserInfoJsontTypes | uploadCompanyInfoJsonTypes;
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
	if (props.currentUser === 'UpdateUserInfo') {
		const scanQrData = props.scanQrData as uploadUserInfoJsontTypes;
		data.name = scanQrData.userName;
		data.id = scanQrData.userIDCard;
	} else {
		const scanQrData = props.scanQrData as uploadCompanyInfoJsonTypes;
		data.name = scanQrData.legalPersonName;
		data.id = scanQrData.legalPersonIDCard;
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
const stores = useUserInfo();

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
        stores.setUserType({ type: res.data.type, subType: res.data.subType });
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
</style>
