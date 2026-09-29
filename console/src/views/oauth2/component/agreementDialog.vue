<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px"> <div v-html="contentHtml" class="dialog-content"></div></el-dialog>
</template>

<script setup lang="ts">
import { reactive, toRefs, ref } from 'vue';
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '服务协议',
	},
});

const contentHtml = ref();
const { dialog } = toRefs(state);
const openDialog = (content: string) => {
  contentHtml.value = content;
	dialog.value.isShowDialog = true;
};

defineExpose({
	openDialog,
});
</script>

<style lang="scss" scoped>
.dialog-content {
	min-height: 500px;
}
</style>
