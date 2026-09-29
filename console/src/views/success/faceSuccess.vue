<template>
	<div class="container">
		<div class="main">
			<SvgIcon name="iconfont icon--wanchengchenggong" :size="60" class="icon left-item-animation"> </SvgIcon>
			<div class="tip left-item-animation">人脸识别成功</div>
			<div class="left-item-animation pay-text">您已完成人脸识别</div>
			<div class="left-item-animation left-item-btn mt10">
				<el-button type="primary" size="default" round @click="onGoHome">返回首页</el-button>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { NextLoading } from '/@/utils/loading';
import { useRouter } from 'vue-router';
// 定义变量内容
const router = useRouter();
onMounted(() => {
	NextLoading.done();
});
// 返回首页
const onGoHome = () => {
	router.push('/');
};
</script>

<style lang="scss" scoped>
.main {
	width: 100%;
	height: calc(100vh + 40px);
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	background: var(--next-bg-color);
	margin-top: -40px;
	overflow: hidden;
	.icon {
		color: var(--el-color-primary);
	}
	.tip {
		color: var(--el-color-primary);
		font-size: 20px;
		font-weight: 600;
		margin-top: 10px;
	}
	.pay-text {
		color: var(--el-text-color-primary);
		margin-top: 4px;
	}
	.left-item-animation {
		opacity: 0;
		animation-name: error-num;
		animation-duration: 0.5s;
		animation-fill-mode: forwards;
	}
	.left-item-btn {
		animation-delay: 0.2s;
	}
}
</style>
