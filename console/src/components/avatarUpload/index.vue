<template>
	<!-- 头像上传组件 -->
	<el-upload
		action="https://jsonplaceholder.typicode.com/posts/"
		:show-file-list="false"
		:disabled="disabledType"
		:auto-upload="false"
		accept=".png,.jepg,.jpg"
		ref="excelUploadRef"
		:limit="1"
		:on-change="handleChange"
		class="avatar-box"
	>
		<img v-if="imagesURL || props.imageUrl" :src="imagesURL || props.imageUrl" class="avatar" />
	</el-upload>
	<!-- 头像裁切 -->
	<CropperDialog ref="cropperDialogRef" @cropperBack="cropperBack" />
</template>
<script setup lang="ts">
import { ref, defineAsyncComponent, onMounted } from 'vue';
import { compressionFile } from '/@/utils/compressionImage';
import { Session, Local } from '/@/utils/storage';
import { ElMessage } from 'element-plus';
const CropperDialog = defineAsyncComponent(() => import('/@/components/cropper/index.vue'));

interface Props {
	imageUrl?: string; // 回显图片地址
	action?: string; //   上传地址
	imgText?: string; //   文字可以不传
	imgUpText?: string; // 上传按钮的文字
	disabledType?: boolean; // 是否禁用上传
	userId?: string; //用户id
	type?: string;
}
const props = withDefaults(defineProps<Props>(), {
	imageUrl: '',
	action: '/activity/resource/uploadFile',
	imgText: '支持jpg/jpeg/bmp；文件大小不能超过2M',
	imgUpText: '上传图片',
	disabledType: false,
	type: '',
});
const imagesURL = ref<any>(props.imageUrl);
const emits = defineEmits(['imgSuccess']);
const cropperDialogRef = ref();
const onCropperDialogOpen = (imgbase64: string) => {
	cropperDialogRef.value.openDialog(imgbase64);
};
const middleImagesURL = ref();
const cropperBack = (file: string) => {
	// imagesURL.value = middleImagesURL.value;
	emits('imgSuccess', file);
};
const handleChange = async (file: any) => {
	let rawFile = file.raw;
	if (rawFile.size / 1024 / 1024 > 2) {
		await compressionFile(rawFile).then((res: any) => {
			rawFile = res;
		});
	}

	if (rawFile.type !== 'image/jpeg' && rawFile.type !== 'image/png' && rawFile.type !== 'image/jpg' && rawFile.type !== 'image/bmp') {
		ElMessage.error('仅支持格式为 jpeg, png,bmp,文件大小不能超过2M的图片');
		return false;
	} else if (rawFile.size / 1024 / 1024 > 2) {
		ElMessage.error('图片文件大小不能超过2MB!');
		return false;
	} else {
		new Promise((resolve, reject) => {
			//读取文件内容转成base64格式回显图片
			const reader = new FileReader();
			reader.readAsDataURL(rawFile);
			reader.onload = () => {
				middleImagesURL.value = reader.result;
				onCropperDialogOpen(middleImagesURL.value);
				resolve(file);
			};
			reader.onerror = (error) => reject(error);
		});
	}
	return true;
};
const onUpdate = ()=>{
  if (props.type === 'defaultAvatar') {
    imagesURL.value = import.meta.env.VITE_API_URL + '/public/header/user?id=' + '&date=' + new Date();
  } else {
    let userInfo = Session.get('userInfo') || Local.get('userInfo');
    imagesURL.value = import.meta.env.VITE_API_URL + '/public/header/user?id=' + userInfo.user.id + '&date=' + new Date();
  }
}
onMounted(() => {
  onUpdate()
});
defineExpose({
  onUpdate,
});
</script>

<style lang="scss" scoped>
:deep(.el-upload) {
	width: 100%;
	height: 100%;
}
.avatar-box {
	width: 100%;
	height: 100%;
	.avatar {
		width: 100%;
		height: 100%;
		border-radius: 4px;
	}
}
</style>
