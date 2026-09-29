<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<div class="login-scan-container">
			<div ref="qrcodeRef" id="weixin-qrcode"></div>
		</div>
	</el-dialog>
</template>

<script setup lang="ts" name="userWeixin">
import { ref, onMounted, onUnmounted, reactive, nextTick } from 'vue';
import { useWeixinApi } from '/@/api/weixin/index';
const qrcodeRef = ref<HTMLElement | null>(null);
const emit = defineEmits(['refresh']);
const useWeixinApiCollect = useWeixinApi();
const getWeixinConfig = () => {
	useWeixinApiCollect.weixinConfig().then((res: any) => {
		if (res.code === "SUCCESS") {
			let appid = res.data.appID;
			// @ts-ignore
			new WxLogin({
				self_redirect: true,
				id: 'weixin-qrcode',
				appid: appid,
				scope: 'snsapi_login',
				redirect_uri: encodeURIComponent(import.meta.env.VITE_WECHAT_BIND_SUCCESS_CALLBACK),
				state: '',
				style: '',
				href: '',
			});
		}
	});
};
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '绑定微信',
	submitTxt: '确认',
});
const openDialog = () => {
	dialog.isShowDialog = true;
	getWeixinConfig();
};

const closeDialog = () => {
	dialog.isShowDialog = false;
};
const intervalId = ref();
onMounted(() => {
	nextTick(() => {
		intervalId.value = setInterval(() => {
			if (sessionStorage.getItem('weixinStatus') == '2') {
				closeDialog();
				emit('refresh');
				clearInterval(intervalId.value);
			}
		}, 500);
	});
});
onUnmounted(() => {
	clearInterval(intervalId.value);
});
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss">
.login-scan-animation {
	opacity: 0;
	animation-name: error-num;
	animation-duration: 0.5s;
	animation-fill-mode: forwards;
}
.login-scan-container {
	padding: 0 20px 20px;
	display: flex;
	flex-direction: column;
	text-align: center;
	@extend .login-scan-animation;
	animation-delay: 0.1s;
	:deep(img) {
		margin: auto;
	}
	.login-msg {
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--el-text-color-placeholder);
		@extend .login-scan-animation;
		animation-delay: 0.2s;
	}
}
</style>
