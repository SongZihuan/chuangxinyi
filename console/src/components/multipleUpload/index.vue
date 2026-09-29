<template>
	<el-upload
		class="aupload-demo"
		:show-file-list="true"
		:auto-upload="false"
		ref="uploadRef"
		multiple
		:on-change="handleChange"
    :on-remove="handleRemove"
		v-model:file-list="fileList"
		:limit="5"
	>
		<el-button type="primary">{{ fileUpText }}</el-button>
		<template #tip>
			<div class="el-upload__tip">{{ fileText }}</div>
		</template>
	</el-upload>
</template>

<script lang="ts" setup>
import { ElMessage } from 'element-plus';
import { ref } from 'vue';
interface Props {
	fileText?: string; //   文字可以不
	fileUpText?: string; // 上传按钮的文字
  fileList?: any[]; // 文件列表
}

const props = withDefaults(defineProps<Props>(), {
	fileText: '最多可以上传5个附件',
	fileUpText: '上传',
  fileList: []
});
const fileList = ref(props.fileList);
const emits = defineEmits(['fileSuccess']);
const handleChange = async (file: any, uploadFiles: any) => {
	let rawFile = file.raw;
	if (rawFile.size / 1024 / 1024 > 2) {
		ElMessage.error('图片文件大小不能超过5MB!');
		return false;
	} else {
		fileList.value = uploadFiles;
		emits('fileSuccess', fileList.value);
	}
	return true;
};
const handleRemove = () => {
  emits('fileSuccess', fileList.value);
};
</script>

<style lang="scss" scoped></style>
