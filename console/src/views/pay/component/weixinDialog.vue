<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="650px" destroy-on-close @close="closeDialog">
		<div class="center">
			<div ref="qrcodeRef"></div>
			<div class="font12 mt20 login-msg">
				<i class="iconfont icon-saoyisao mr5"></i>
				<span>打开微信扫一扫</span>
			</div>
			<div class="tip" @click="handleQuery"><el-link type="primary">支付成功,点击这里</el-link></div>
		</div>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, nextTick } from 'vue';
import QRCode from 'qrcodejs2-fixes';
const emit = defineEmits(['queryOrderBack']);
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '微信充值',
	submitTxt: '确认',
});
const qrcodeRef = ref<HTMLElement | null>(null);
const codeUrl = ref<string>('');
const initQrcode = () => {
	nextTick(() => {
		(<HTMLElement>qrcodeRef.value).innerHTML = '';
		new QRCode(qrcodeRef.value, {
			text: codeUrl.value,
			width: 260,
			height: 260,
			colorDark: '#000000',
			colorLight: '#ffffff',
		});
	});
};
const openDialog = (url: string) => {
	codeUrl.value = url;
	dialog.isShowDialog = true;
	initQrcode();
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
.center {
	display: flex;
	flex-direction: column;
	align-items: center;
}
</style>
