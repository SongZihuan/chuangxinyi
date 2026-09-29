<template>
	<div class="uploader" :style="customStyle">
		<el-upload
			class="avatar-uploader"
			:show-file-list="false"
			:disabled="disabledType"
			:auto-upload="false"
			accept=".png,.jepg,.jpg"
			ref="uploadRef"
			:on-change="handleChange"
			:style="upImgBoxCustomStyle"
		>
			<img v-if="imagesURL || props.imageUrl" :src="props.imageUrl ? props.imageUrl : imagesURL" class="avatar" />
			<div class="upImgBox" :style="upImgBoxCustomStyle" v-else>
				<el-icon class="avatar-uploader-icon">
					<ele-Plus />
				</el-icon>
				<div>{{ imgUpText }}</div>
			</div>
		</el-upload>
		<div class="upImgText">
			{{ imgText }}
		</div>
	</div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { compressionFile } from '/@/utils/compressionImage';

interface Props {
	imageUrl?: string; // 回显图片地址
	action?: string; //   上传地址
	imgText?: string; //   文字可以不传
	imgUpText?: string; // 上传按钮的文字
	disabledType?: boolean; // 是否禁用上传
	customStyle?: { maxWidth: string }; // 自定义样式
	upImgBoxCustomStyle?: { width: string; height: string }; // 自定义样式
}

const props = withDefaults(defineProps<Props>(), {
	imageUrl: '',
	action: '',
	imgText: '支持jpg/jpeg/bmp；文件大小不能超过2M',
	imgUpText: '上传图片',
	disabledType: false,
	customStyle: () => ({
		maxWidth: '398px',
	}),
	upImgBoxCustomStyle: () => ({
		width: '150px',
		height: '150px',
	}),
});

const imagesURL = ref<any>(props.imageUrl);
const emits = defineEmits(['imgSuccess']);
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
		emits('imgSuccess', rawFile);
		new Promise((resolve, reject) => {
			//读取文件内容转成base64格式回显图片
			const reader = new FileReader();
			reader.readAsDataURL(rawFile);
			reader.onload = () => {
				imagesURL.value = reader.result;
				resolve(file);
			};
			reader.onerror = (error) => reject(error);
		});
	}
	return true;
};
watch(
	() => props.imageUrl,
	() => {
		imagesURL.value = props.imageUrl;
	}
);
</script>

<style lang="scss" scoped>
.uploader {
	margin-top: 10px;
	display: flex;
	flex-direction: column;
	align-items: center;
}

:deep().avatar-uploader {
	.avatar {
		width: 100%;
		height: 100%;
		display: block;
		object-fit: cover;
	}

	.el-upload {
		width: 100%;
		height: 100%;
		border: 1px dashed #dcdfe6;
		border-radius: 6px;
		cursor: pointer;
		position: relative;
		overflow: hidden;
		transition: 0.2s;
		background: rgba(0, 0, 0, 0.04) !important;
	}

	.el-upload:hover {
		border-color: #14b194;
	}
}

.el-icon.avatar-uploader-icon {
	width: 100%;
	font-size: 16px;
	color: rgba(0, 0, 0, 0.45);
	text-align: center;
}

.upImgBox {
	width: 150px;
	height: 150px;
	font-size: 14px;
	font-weight: 400;
	color: rgba(0, 0, 0, 0.65);
	text-align: center;
	padding-top: 24px;
	box-sizing: border-box;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-content: center;
}

.upImgText {
	font-size: 14px;
	text-align: center;
	color: rgba(0, 0, 0, 0.6);
	margin-top: 4px;
}
</style>
