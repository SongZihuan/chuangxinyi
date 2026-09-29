<template>
	<div id="slider-con">
		<div :id="sliderId"></div>
	</div>
</template>
<script setup lang="ts">
import { onMounted, nextTick } from 'vue';
import { useBaseApi } from '/@/api/base';
import { Session } from '/@/utils/storage';
const emit = defineEmits(['siderEmit']);
const sliderId = 'sliderId' + new Date();
const initSide = async () => {
	// @ts-ignore
	const siderWidth = document.getElementById('slider-con').offsetWidth;

  await async function () {
    return new Promise((resolve) => {
      if (Session.get('appkey')) {
        resolve(0);
      } else {
        // 设置一个定时器来检查条件
        const interval = setInterval(() => {
          if (Session.get('appkey')) {
            clearInterval(interval);
            resolve(0);
          }
        }, 100);
      }
    });
  }()
  
	// @ts-ignore
	AWSC.use('nc', function (state?: any, module?: any) {
		// 初始化
		// @ts-ignore
		window.nc = module.init({
			// 应用类型标识。它和使用场景标识（scene字段）一起决定了滑动验证的业务场景与后端对应使用的策略模型。您可以在阿里云验证码控制台的配置管理页签找到对应的appkey字段值，请务必正确填写。
			appkey: Session.get('appkey'),
			//使用场景标识。它和应用类型标识（appkey字段）一起决定了滑动验证的业务场景与后端对应使用的策略模型。您可以在阿里云验证码控制台的配置管理页签找到对应的scene值，请务必正确填写。
			scene: 'nc_other',
			// 声明滑动验证需要渲染的目标ID。
			renderTo: sliderId,
			width: siderWidth,
			//前端滑动验证通过时会触发该回调参数。您可以在该回调参数中将会话ID（sessionId）、签名串（sig）、请求唯一标识（token）字段记录下来，随业务请求一同发送至您的服务端调用验签。
			success: function (data: any) {
				emit('siderEmit', data);
			},
			// 前端二次验证失败时触发该回调参数
			fail: function (failCode) {
				console.log(failCode);
			},
			// 前端二次验证加载异常时触发该回调参数。
			error: function (errorCode) {
				console.log(errorCode);
			},
		});
	});
};
const getAppkey = async () => {
	await useBaseApi()
		.afsGet()
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				Session.set('appkey', res.data.hAppKey);
			}
		});
};
const resetSider = () => {
	window.nc.reset();
};
onMounted(async () => {
	if (!Session.get('appkey')) {
		await getAppkey();
	}
	nextTick(() => {
		initSide();
		setTimeout(() => {
			resetSider();
		}, 50000);
	});
});
defineExpose({
	resetSider,
});
</script>
<style lang="scss" scoped>
#slider-con {
	width: 100%;
	padding-bottom: 20px;
	margin-bottom: 20px;
}
.nc-container #nc_1_wrapper {
	width: 100% !important;
	.nc-lang-cnt {
		width: 100%;
	}
}
</style>
