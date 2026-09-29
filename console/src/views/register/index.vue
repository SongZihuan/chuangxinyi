<template>
	<div class="register-container flex">
		<div class="register-left">
			<div class="register-left-img">
				<img :src="registerMain" />
			</div>
		</div>
		<div class="register-right flex">
			<div class="register-right-warp flex-margin">
				<span class="register-right-warp-one"></span>
				<span class="register-right-warp-two"></span>
				<div class="register-right-warp-mian">
					<div class="register-right-warp-main-title">注册</div>
					<div class="register-right-warp-main-form">
						<div>
							<el-tabs v-model="state.tabsActiveName">
								<el-tab-pane label="手机号注册" name="com" v-if="registerStep === 1">
									<Mobile @nextStep="nextStep" />
								</el-tab-pane>
							</el-tabs>
						</div>
						<div class="register">
							<div>已经有账号?</div>
							<div @click="goLogin">立即登录</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts" name="registerIndex">
import { onMounted, reactive, ref } from 'vue';
import {useRoute, useRouter} from 'vue-router';
import { NextLoading } from '/@/utils/loading';
import Mobile from '/@/views/register/component/mobile.vue';
import useFile from "/@/hooks/useFile";

const registerMain = useFile().getFile("register");
const router = useRouter();
const state = reactive({
	tabsActiveName: 'com',
});
const registerStep = ref<number>(1);
const route = useRoute()
const goLogin = () => {
	router.push({ path: '/login', query: route.query});
};

const nextStep = () => {
	registerStep.value++;
};
// 页面加载时
onMounted(() => {
	NextLoading.done();
});
</script>

<style scoped lang="scss">
:deep(.el-tabs__nav-scroll) {
	display: flex;
	justify-content: center;
}

.register-action {
	display: flex;
	flex-direction: row;
	align-items: center;
	position: relative;
	z-index: 2004;
	.admin-agreement {
		color: #002fa7;
		margin-left: 4px;
		margin-right: 30px;
		text-decoration: underline;
	}
}
:deep(.el-tabs__active-bar) {
	height: 3px;
}

:deep(input::-webkit-input-placeholder) {
	color: #9ca5ba;
	font-size: 12px;
}

:deep(.el-input__inner) {
	background-color: transparent !important;
}
:deep(.el-tabs__item.is-active) {
	color: #121212;
	font-weight: bolder;
}

.register-container {
	height: 100%;
	background: var(--el-color-white);
	.register-left {
		flex: 1;
		position: relative;
		background-color: rgba(211, 239, 255, 1);
		margin-right: 100px;
		.register-left-logo {
			display: flex;
			align-items: center;
			position: absolute;
			top: 50px;
			left: 80px;
			z-index: 1;
			animation: logoAnimation 0.3s ease;
			img {
				width: 52px;
				height: 52px;
			}
			.register-left-logo-text {
				display: flex;
				flex-direction: column;
				span {
					margin-left: 10px;
					font-size: 28px;
					color: #26a59a;
				}
				.register-left-logo-text-msg {
					font-size: 12px;
					color: #32a99e;
				}
			}
		}
		.register-left-img {
			position: absolute;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
			width: 50%;
			height: 52%;
			img {
				width: 100%;
				height: 100%;
				animation: error-num 0.6s ease;
        object-fit: contain;
			}
		}
		.register-left-waves {
			position: absolute;
			top: 0;
			right: -100px;
		}
	}
	.register-right {
		width: 700px;

		.register-right-warp {
			border: 1px solid var(--el-color-primary-light-3);
			border-radius: 3px;
			width: 500px;
			min-height: 350px;
			position: relative;
			overflow: hidden;
			background-color: var(--el-color-white);
			.back {
				padding: 5px;
				position: absolute;
				top: 0px;
				left: 0px;
				cursor: pointer;
				z-index: 2004;
			}
			.next {
				padding: 5px;
				position: absolute;
				top: 0px;
				right: 0px;
				cursor: pointer;
				z-index: 2004;
			}
			.register-right-warp-one,
			.register-right-warp-two {
				position: absolute;
				display: block;
				width: inherit;
				height: inherit;
				&::before,
				&::after {
					content: '';
					position: absolute;
					z-index: 1;
				}
			}
			.register-right-warp-one {
				&::before {
					filter: hue-rotate(0deg);
					top: 0px;
					left: 0;
					width: 100%;
					height: 3px;
					background: linear-gradient(90deg, transparent, var(--el-color-primary));
					animation: registerLeft 3s linear infinite;
				}
				&::after {
					filter: hue-rotate(60deg);
					top: -100%;
					right: 2px;
					width: 3px;
					height: 100%;
					background: linear-gradient(180deg, transparent, var(--el-color-primary));
					animation: registerTop 3s linear infinite;
					animation-delay: 0.7s;
				}
			}
			.register-right-warp-two {
				&::before {
					filter: hue-rotate(120deg);
					bottom: 2px;
					right: -100%;
					width: 100%;
					height: 3px;
					background: linear-gradient(270deg, transparent, var(--el-color-primary));
					animation: registerRight 3s linear infinite;
					animation-delay: 1.4s;
				}
				&::after {
					filter: hue-rotate(300deg);
					bottom: -100%;
					left: 0px;
					width: 3px;
					height: 100%;
					background: linear-gradient(360deg, transparent, var(--el-color-primary));
					animation: registerBottom 3s linear infinite;
					animation-delay: 2.1s;
				}
			}
			.register-right-warp-mian {
				display: flex;
				flex-direction: column;
				height: 100%;
				.register-right-warp-main-title {
					height: 80px;
					line-height: 80px;
					font-size: 27px;
					text-align: center;
					letter-spacing: 3px;
					animation: logoAnimation 0.3s ease;
					animation-delay: 0.3s;
					color: var(--el-text-color-primary);
				}
				.register-right-warp-main-form {
					flex: 1;
					padding: 0 50px 50px;
					.register-content-main-sacn {
						position: absolute;
						top: 0;
						right: 0;
						width: 50px;
						height: 50px;
						overflow: hidden;
						cursor: pointer;
						transition: all ease 0.3s;
						color: var(--el-color-primary);
						&-delta {
							position: absolute;
							width: 35px;
							height: 70px;
							z-index: 2;
							top: 2px;
							right: 21px;
							background: var(--el-color-white);
							transform: rotate(-45deg);
						}
						&:hover {
							opacity: 1;
							transition: all ease 0.3s;
							color: var(--el-color-primary) !important;
						}
						i {
							width: 47px;
							height: 50px;
							display: inline-block;
							font-size: 48px;
							position: absolute;
							right: 1px;
							top: 0px;
						}
					}
					.register {
						display: flex;
						justify-content: center;
						font-size: 14px;
						margin-top: 20px;
						cursor: pointer;
						position: relative;
						z-index: 2004;
						:nth-child(1) {
							color: #9ca5ba;
						}
						:nth-child(2) {
							color: var(--el-color-primary);
							margin-left: 6px;

							text-decoration: underline;
						}
					}
				}
			}
		}
	}
}
</style>
