<template>
	<el-image
		:style="{ width: `100%`, height: `100%`, borderRadius: '4px' }"
		:src="imagesURL"
		:zoom-rate="1.2"
		:preview-src-list="[imagesURL]"
		preview-teleported
		v-if="imagesURL"
		fit="cover"
		close-on-press-escape
	/>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Session, Local } from '/@/utils/storage';
interface Props {
	imageUrl?: string; // 回显图片地址
	user?: number | string;
}
const props = withDefaults(defineProps<Props>(), {
	imageUrl: '',
	user: '',
});

const imagesURL = ref<any>(props.imageUrl);

onMounted(() => {
	let userInfo = Session.get('userInfo') || Local.get('userInfo');
	if (userInfo && userInfo.user.header) {
		imagesURL.value = import.meta.env.VITE_API_URL + '/public/header/user?id=' + props.user;
	} else if (userInfo.user.wechatHeader) {
		imagesURL.value = userInfo.user.wechatHeader;
	} else {
		imagesURL.value = import.meta.env.VITE_API_URL + '/public/header/user?id=' + props.user;
	}
});
</script>

<style lang="scss" scoped></style>
