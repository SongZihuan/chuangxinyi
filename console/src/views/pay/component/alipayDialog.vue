<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="650px" destroy-on-close @close="closeDialog">
		<iframe :src="codeUrl" frameborder="0" class="ifr"></iframe>
		<div class="tip" @click="handleQuery"><el-link type="primary">支付成功,点击这里</el-link></div>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
const emit = defineEmits(['queryOrderBack']);
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '支付宝充值',
	submitTxt: '确认',
});
const codeUrl = ref<string>('');
const openDialog = (url: string) => {
	codeUrl.value = url;
	dialog.isShowDialog = true;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
//查询订单是否成功
const handleQuery = () => {
	emit('queryOrderBack');
};
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style lang="scss" scoped>
.ifr {
	width: 100%;
	height: 560px;
}
.tip {
	display: flex;
	flex-direction: row;
	justify-content: center;
}
</style>
