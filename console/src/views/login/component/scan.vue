<template>
	<div class="login-scan-container">
		<div id="weixin-qrcode"></div>
	</div>
	<!-- 手机注册绑定 -->
	<el-dialog v-model="showWeixin" title="绑定手机号" width="500px" :modal="false">
		<Mobile :isBindWeixin="isBindWeixin" @success="onSuccess" :withLogin="true"></Mobile>
	</el-dialog>
</template>

<script setup lang="ts" name="loginWeixin">
import { onMounted, onUnmounted, defineAsyncComponent, ref } from 'vue';
import { useWeixinApi } from '/@/api/weixin';
import { ElMessageBox } from 'element-plus';
const emit = defineEmits(['signBack', 'totpBack']);
const Mobile = defineAsyncComponent(() => import('/@/views/register/component/mobile.vue'));
const useWeixinApiCollect = useWeixinApi();
const isBindWeixin=ref(true)
const showWeixin = ref<boolean>(false);
sessionStorage.setItem('weixinLoginStatus', '1');
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
				redirect_uri: encodeURIComponent(import.meta.env.VITE_WECHAT_LOGIN_SUCCESS_CALLBACK),
				state: '',
				style: '',
				href: '',
			});
		}
	});
};

const onSuccess = ()=>{
  showWeixin.value=false
}

const intervalId = setInterval(() => {
  const status = sessionStorage.getItem('weixinLoginStatus')
  if (status === '4') {
    //2FA
    clearInterval(intervalId);
    const token = sessionStorage.getItem('weixinToken')
    emit('totpBack', token);
  } else if (status == '2') {
    //绑定时扫码
    clearInterval(intervalId);
    const token = sessionStorage.getItem('weixinToken')
    const tokenType = sessionStorage.getItem("weixinTokenType")
    const tokenSubType = sessionStorage.getItem("weixinTokenSubType")
    emit('signBack', token, tokenType, tokenSubType);
	} else if (status == '3') {
    //未绑定手机时扫码
		clearInterval(intervalId);
		ElMessageBox.confirm('您的微信暂未绑定手机号,是否现在绑定?', '提示', {
			confirmButtonText: '去绑定',
			cancelButtonText: '取消',
			type: 'warning',
			center: true,
		})
			.then(() => {
				showWeixin.value = true;
			})
			.catch(() => {});
	}
}, 500);

onMounted(() => {
	getWeixinConfig();
});

onUnmounted(() => {
	clearInterval(intervalId);
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
	// padding: 0px 20px;
	display: flex;
	flex-direction: column;
	min-height: 300px;
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
