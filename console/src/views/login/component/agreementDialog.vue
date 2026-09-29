<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px"> <div v-html="contentHtml" class="dialog-content"></div></el-dialog>
</template>

<script setup lang="ts">
import { reactive, toRefs, onMounted, ref } from 'vue';
import { useAgreementApi } from '/@/api/agreement/index';
const useAgreementApiCollect = useAgreementApi();
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '用户协议',
	},
});
const contentHtml = ref();
const { dialog } = toRefs(state);
const openDialog = () => {
	dialog.value.isShowDialog = true;
};
const getAgreement = () => {
	useAgreementApiCollect.getAgreement({ aid: '用户协议' }).then((res: any) => {
		contentHtml.value = res;
	});
};
onMounted(() => {
	getAgreement();
});
defineExpose({
	openDialog,
});
</script>

<style lang="scss" scoped>
.dialog-content {
	min-height: 500px;
}
</style>
