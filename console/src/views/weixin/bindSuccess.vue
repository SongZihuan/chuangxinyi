<template>
	<div class="warp">
		<el-result title="成功" sub-title="微信绑定成功" icon="success" v-if="isBind">
			<template #extra> </template>
		</el-result>
	</div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { useWeixinApi } from '/@/api/weixin/index';
import { ElMessage, ElLoading } from 'element-plus';

import { ref } from 'vue';
const route = useRoute();
const useWeixinApiCollect = useWeixinApi();
const loadingInstance = ElLoading.service({
	text: '正在绑定微信...',
});
const isBind = ref<boolean>(false);
const setWeixinCode = () => {
	sessionStorage.setItem('weixinStatus', '1');
	if (route.query.code) {
		sessionStorage.setItem('weixincode', route.query.code as string);
		useWeixinApiCollect.weixinAccessToken({ code: route.query.code, type: 'Wechat' }).then((res: any) => {
			if (res.code === "SUCCESS") {
				userBindWeixin(res.data.token);
			}
		});
	}
};
const userBindWeixin = (token: string) => {
	useWeixinApiCollect.bindWeixin({ wechatToken: token, isDelete: false }).then((res: any) => {
		if (res.code === "SUCCESS") {
			loadingInstance.close();
			isBind.value = true;
			ElMessage.success('绑定微信成功');
			sessionStorage.setItem('weixinStatus', '2');
		}
	});
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
