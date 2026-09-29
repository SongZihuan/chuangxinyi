<template>
	<!-- 阿里云静默认证 -->
	<div id="slider-con">
		<div :id="sliderId"></div>
	</div>
</template>
<script setup lang="ts">
import { nextTick, onMounted } from 'vue';
import { Session } from '/@/utils/storage';
import { useBaseApi } from '/@/api/base';
const emit = defineEmits(['siderEmit', 'nvcValEmit']);
const sliderId = 'sliderId' + new Date() + 20;
// 实例化nvc 对无痕验证进行初始化操作
const initSide = () => {
	// @ts-ignore
	const siderWidth = document.getElementById('slider-con').offsetWidth;
	// @ts-ignore
	AWSC.use('nvc', function (state?: any, module?: any) {
		// 初始化 调用module.init进行初始化
		// @ts-ignore
		window.nvc = module.init({
			// 应用类型标识。它和使用场景标识（scene字段）一起决定了无痕验证的业务场景与后端对应使用的策略模型。您可以在阿里云验证码控制台的配置管理页签找到对应的appkey字段值，请务必正确填写。
			appkey: Session.get('appkey'),
			//使用场景标识。它和应用类型标识（appkey字段）一起决定了无痕验证的业务场景与后端对应使用的策略模型。您可以在阿里云验证码控制台的配置管理页签找到对应的scene值，请务必正确填写。
			scene: 'nvc_other',
			width: siderWidth,
			test: module.TEST_NC_PASS,
			// 二次验证获取人机信息串，跟随业务请求一起上传至业务服务器，由业务服务器进行验签。
			success: function (data: any) {
				emit('siderEmit', data);
			},

			fail: function (failCode: any) {
				// @ts-ignore
				console.log(failCode);
			},
			// 前端二次验证加载异常时触发该回调参数。
			error: function (errorCode: any) {
				// @ts-ignore
				console.log(errorCode);
			},
		});
	});
};
// 发送业务请求：点击按钮时触发，主动获取人机信息串，并发送给业务服务端
function registerClick() {
	// @ts-ignore
	window.nvc.getNVCValAsync(function (nvcVal) {
		emit('nvcValEmit', nvcVal);
	});
}
const resetSider = () => {
	// @ts-ignore
	window.nvc.reset();
};
function siderCheck() {
	var ncoption = {
		// 声明滑动验证需要渲染的目标ID。
		renderTo: sliderId,
	};
	// 唤醒二次验证（滑动验证码）
	window.nvc.getNC(ncoption);
}
const getAppkey = async () => {
	await useBaseApi()
		.afsGet()
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				Session.set('appkey', res.data.hAppKey);
			}
		});
};
onMounted(async () => {
	if (!Session.get('appkey')) {
		await getAppkey();
	}
	nextTick(() => {
		initSide();
	});
});

defineExpose({
	registerClick,
	siderCheck,
	resetSider,
});
</script>
<style lang="scss" scoped>
#slider-con {
	width: 100%;
	height: 20px;
	padding-bottom: 30px;
}
.nc-container :deep(#nc_2_wrapper) {
	width: 100% !important;
	.nc-lang-cnt {
		width: 100%;
	}
}
.nc-container :deep(#nc_3_wrapper) {
	width: 100% !important;
	.nc-lang-cnt {
		width: 100%;
	}
}
.nc-container :deep(#nc_1_wrapper) {
	width: 100% !important;
	.nc-lang-cnt {
		width: 100%;
	}
}
</style>
