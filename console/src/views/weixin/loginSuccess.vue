<template>
	<div class="warp">
		<el-result title="成功" sub-title="微信登录成功" icon="success" v-if="isBind">
			<template #extra> </template>
		</el-result>
	</div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { useWeixinApi } from '/@/api/weixin/index';
import { ElLoading } from 'element-plus';
import { useUserInfo } from '/@/stores/userInfo';
import { ref } from 'vue';
const route = useRoute();

const useWeixinApiCollect = useWeixinApi();
const loadingInstance = ElLoading.service({
	text: '正在登录...',
});
const stores = useUserInfo();
const isBind = ref<boolean>(false);
const setWeixinCode = () => {
	sessionStorage.setItem('weixinLoginStatus', '1');
	sessionStorage.setItem('weixinToken', '');
	if (route.query.code) {
		sessionStorage.setItem('weixincode', route.query.code + '');
		useWeixinApiCollect.weixinAccessToken({ code: route.query.code, type: 'Auto' }).then((res: any) => {
			if (res.code === 'SUCCESS') {
				loadingInstance.close();
				if (res.data.type === 'Wechat') {
					isBind.value = true;
					sessionStorage.setItem('weixinLoginStatus', '3');
					sessionStorage.setItem('weixinToken', res.data.token);
					sessionStorage.setItem('weixinTokenType', res.data.type);
				} else if (res.data.type === 'Login2FA') {
					sessionStorage.setItem('weixinLoginStatus', '4');
					sessionStorage.setItem('weixinToken', res.data.token);
					sessionStorage.setItem('weixinTokenType', res.data.type);
				} else {
					sessionStorage.setItem('weixinLoginStatus', '2');
					sessionStorage.setItem('weixinToken', res.data.token);
					sessionStorage.setItem('weixinTokenType', res.data.type);
					sessionStorage.setItem('weixinTokenSubType', res.data.subType);
					stores.setUserType({ type: res.data.type, subType: res.data.subType });
				}
			}
		});
	}
};
setWeixinCode();
</script>

<style lang="scss" scoped>
.warp {
	width: 100%;
	height: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
}
</style>
