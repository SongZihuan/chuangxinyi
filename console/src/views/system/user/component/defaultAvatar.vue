<template>
	<el-dialog :title="props.title" v-model="dialog.isShowDialog" width="400px" height="500px">
		<div class="dialog-content">
			<div style="width: 150px; height: 150px">
				<AvatarUpload @imgSuccess="imgSuccess" :key="avatarUploadKey" :type="defaultAvatar"></AvatarUpload>
			</div>
		</div>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, toRefs, defineAsyncComponent, ref } from 'vue';
import { useBaseApi } from '/@/api/base';
import { ElMessage } from 'element-plus';
const AvatarUpload = defineAsyncComponent(() => import('/@/components/avatarUpload/index.vue'));
interface Props {
	title: string;
	contentHtml: HtmlType;
	subtitle: string;
}
const emit = defineEmits(['refresh']);
const defaultAvatar = 'defaultAvatar';
const props = withDefaults(defineProps<Props>(), {
	title: '点击上传',
	contentHtml: '',
	subtitle: '',
});
const avatarUploadKey = ref(1);
const state = reactive({
	dialog: {
		isShowDialog: false,
	},
});
const { dialog } = toRefs(state);
const openDialog = () => {
	dialog.value.isShowDialog = true;
};

const imgSuccess = async (base64: string) => {
	await useBaseApi()
		.uploadHeader({ header: base64, isDelete: false })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				ElMessage({
					type: 'success',
					message: '默认头像设置成功',
				});
        emit('refresh');
			}
		});

	avatarUploadKey.value++;
};

defineExpose({
	openDialog,
});
</script>

<style lang="scss" scoped>
.subtitle {
	color: var(--el-text-color-secondary);
	text-align: center;
	padding-bottom: 20px;
}
.dialog-content {
	min-height: 200px;
	display: flex;
	align-items: center;
	justify-content: center;
}
</style>
