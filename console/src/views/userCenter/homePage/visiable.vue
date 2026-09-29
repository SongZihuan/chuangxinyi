<template>
	<div class="visiable">
		<template v-if="visiable">
			<div class="home-page-container container layout-padding" v-if="isOpen">
				<div class="left-container">
					<!--      显示用户信息logo-->
					<img class="img" :src="registerMain" />
					<div class="title">用户信息</div>
				</div>
				<div class="right-container">
					<div class="card-container">
						<el-card shadow="hover">
              <el-descriptions title="简介" direction="vertical" :column="column">
                <el-descriptions-item label="">{{ homePageData.introduction ? homePageData.introduction : "-" }}</el-descriptions-item>
              </el-descriptions>
              <el-descriptions title="联系信息" direction="vertical" :column="column">
                <el-descriptions-item label="地址">
                  <div class="cell-item">
                    <SvgIcon name="ele-MapLocation" :size="20" class="mr5"></SvgIcon>
                    {{ homePageData.address ? homePageData.address : "-" }}
                  </div>
                </el-descriptions-item>
                <el-descriptions-item label="电话">
                  <div class="cell-item">
                    <el-icon class="mr5"><ele-Phone /></el-icon>
                    {{ homePageData.phone ? homePageData.phone : "-" }}
                  </div>
                </el-descriptions-item>
                <el-descriptions-item label="邮箱">
                  <div class="cell-item">
                    <el-icon class="el-input__icon mr5"><ele-Message /></el-icon>
                    {{ homePageData.email ? homePageData.email : "-" }}
                  </div>
                </el-descriptions-item>
                <el-descriptions-item label="微信">
                  <div class="cell-item">
                    <SvgIcon name="my-weixin" :size="20" class="mr5"></SvgIcon>
                    {{ homePageData.wechat ? homePageData.wechat : "-" }}
                  </div>
                </el-descriptions-item>
                <el-descriptions-item label="QQ">
                  <div class="cell-item">
                    <SvgIcon name="fa fa-qq" :size="15" class="mr5"></SvgIcon>
                    {{ homePageData.qq ? homePageData.qq : "-" }}
                  </div>
                </el-descriptions-item>
              </el-descriptions>
              <el-descriptions title="个人信息" direction="vertical" :column="column">
                <el-descriptions-item label="性别"> {{ homePageData.sex ? homePageData.sex : "-" }}</el-descriptions-item>
                <el-descriptions-item label="公司">{{ homePageData.company ? homePageData.company : "-" }}</el-descriptions-item>
                <el-descriptions-item label="行业"> {{ homePageData.industry ? homePageData.industry : "-" }}</el-descriptions-item>
                <el-descriptions-item label="职位"> {{ homePageData.position ? homePageData.position : "-" }}</el-descriptions-item>
              </el-descriptions>
              <el-descriptions title="官网链接" direction="vertical" :column="column" v-if="homePageData.link">
                <el-descriptions-item label="">
                  <div class="cell-item" @click="onSkip(homePageData)">
                    <el-icon class="mr5"> <ele-Share /> </el-icon>
                    <span class="link"> {{ homePageData.link }} </span>
                  </div>
                </el-descriptions-item>
              </el-descriptions>
						</el-card>
					</div>
				</div>
			</div>
			<div class="noAccess-container" v-else>
				<noAccess />
			</div>
		</template>
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { homePageTypes, pageHomeParamsTypes } from '/@/api/homePage/types';
import { useRouter } from 'vue-router';
import { useHomePageApi } from '/@/api/homePage';
import { ElMessage } from 'element-plus';
import noAccess from '/@/views/error/404.vue';
import useSkip from '/@/hooks/useSkip';
import useFile from '/@/hooks/useFile';
import NProgress from 'nprogress';
import {NextLoading} from "/@/utils/loading";
const registerMain = useFile().getFile('register');
const router = useRouter();
const column = ref(3);
const isOpen = ref(false);
const homePageData = ref<homePageTypes>({
	company: '',
	introduction: '',
	address: '',
	phone: '',
	email: '',
	wechat: '',
	qq: '',
	sex: '',
	industry: '',
	position: '',
	link: '',
});
const visiable = ref(false);
const { onSkip } = useSkip();
// 获取homepage数据
const getHomePageData = async (userID: string) => {
	if (!userID) {
		ElMessage.warning('抱歉，您没有登录权限');
		return;
	}
	let data: pageHomeParamsTypes = {
		userID: userID,
	};
	await useHomePageApi()
		.getHomePage(data)
		.then((res: any) => {
			visiable.value = true;
			if (res.code === "SUCCESS") {
				isOpen.value = true;
				homePageData.value = res.data.homepage;
				return;
			} else {
        isOpen.value = false;
      }
		});
};

onMounted(() => {
	const clientWidth = document.body.clientWidth;
	if (clientWidth < 1000) {
		column.value = 1;
	}
	getHomePageData(router.currentRoute.value.query.userID as string);
});
</script>

<style scoped lang="scss">
.visiable {
	width: 100%;
	height: 100%;
	background: url('../../../assets/bg.png') no-repeat center center / cover;
	//兼容pc和移动端
	.home-page-container {
		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
		width: 100%;
		height: 100%;
		.left-container {
			display: flex;
			flex-direction: column;
			justify-content: center;
			align-items: center;
			width: 0%;
			@media screen and (min-width: 768px) {
				width: 40%;
				height: 100%;
				.img {
					width: 400px;
					height: 400px;
				}
				.title {
					font-size: 38px;
					font-weight: bold;
					color: #ffffff;
				}
			}
			@media screen and (max-width: 768px) {
				width: 0%;
				height: 0%;
				.img {
					width: 200px;
					height: 200px;
				}
				.title {
					font-size: 20px;
					font-weight: bold;
					color: #ffffff;
				}
			}
		}

		.right-container {
			display: flex;
			flex-direction: column;
			justify-content: center;
			align-items: center;
			@media screen and (min-width: 768px) {
				width: 60%;
				height: 100%;
			}
			@media screen and (max-width: 768px) {
				width: 100%;
				height: 100%;
			}

			.card-container {
				width: 80%;
				height: 80%;

				.empty {
					width: 100%;
					height: 100%;

					.el-card {
						display: flex;
						justify-content: center;
						align-items: center;
						height: 100%;
					}
				}

				.el-card {
					width: 100%;
					height: 100%;
					overflow-y: scroll;

					.el-descriptions {
						width: 100%;
						height: 100%;

						.el-descriptions__title {
							font-size: 20px;
							font-weight: bold;
							color: #ffffff;
						}

						.el-descriptions__item {
							font-size: 16px;
							color: #ffffff;
						}

						.cell-item {
							display: flex;
							align-items: center;
						}

						.link {
							cursor: pointer;
							color: var(--el-color-primary);
						}
					}
				}
			}
		}
	}

	.noAccess-container {
		width: 100%;
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}
}
</style>
