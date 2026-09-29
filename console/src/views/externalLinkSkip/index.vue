<template>
	<div class="external-link-skip">
		<div class="external-link-skip-box">
			<div class="content">
				<div class="logo">
					<img :src="logo" />
				</div>
				<div class="title"><el-text tag="b" size="large">即将离开{{projectName}}，请注意账号财产安全</el-text></div>
				<div class="link">{{ router.currentRoute.value.query.link }}</div>
				<div class="btns">
					<el-button round type="primary" v-waves class="login-content-submit" @click="onSkip">
						<span>点击前往</span>
					</el-button>
				</div>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router';
import useFile from "/@/hooks/useFile";
const logo = useFile().getFile("logo");
const router = useRouter();
const projectName = import.meta.env.VITE_PROJECT_NAME
const onSkip = () => {
	let link = router.currentRoute.value.query.link as srting;
	let url;
	if (link.startsWith('https://') || link.startsWith('https://')) {
		url = link;
	} else {
		url = 'https://' + link;
	}
	(window as any).open(url,'_self');
};
</script>

<style scoped lang="scss">
.external-link-skip-box {
	position: absolute;
	left: 50%;
	top: 30%;
	max-width: 624px;
	width: 86%;
	background-color: #fff;
	transform: translateX(-50%);
	padding: 30px 40px;
	box-sizing: border-box;
	border: 1px solid #e5e6eb;
	border-radius: 2px;
	.content {
		.logo {
			width: 120px;
			height: 40px;
			margin-bottom: 20px;
			img {
				width: 100%;
				height: 100%;
				object-fit: cover;
			}
		}
		.title {
			margin: 0;
			font-size: 18px;
			line-height: 24px;
			margin-bottom: 20px;
		}
		.link {
			font-size: 16px;
			padding-bottom: 20px;
			border-bottom: 1px solid #e5e6eb;
		}
		.btns {
			display: flex;
			justify-content: flex-end;
			margin-top: 20px;
		}
	}
}
</style>
