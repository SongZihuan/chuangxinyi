<template>
	<div class="home-page-container container layout-padding">
		<el-card shadow="hover">
			<template v-if="isOpen">
				<el-descriptions title="公司信息" direction="vertical" :column="column">
					<template #extra>
						<el-button type="primary" text @click="onEdit" v-if="authUser.updateHomepage">修改</el-button>
						<el-button type="danger" text @click="onDel">清空</el-button>
					</template>
					<el-descriptions-item label="公司">{{ homePageData.company }}</el-descriptions-item>
					<el-descriptions-item label="公司简介">{{ homePageData.introduction }}</el-descriptions-item>
				</el-descriptions>
				<el-descriptions title="联系信息" direction="vertical" :column="column">
					<el-descriptions-item label="地址">
						<div class="cell-item">
							<SvgIcon name="ele-MapLocation" :size="20" class="mr5"></SvgIcon>
							{{ homePageData.address }}
						</div>
					</el-descriptions-item>
					<el-descriptions-item label="电话">
						<div class="cell-item">
							<el-icon class="mr5"><ele-Phone /></el-icon>
							{{ homePageData.phone }}
						</div>
					</el-descriptions-item>
					<el-descriptions-item label="邮箱">
						<div class="cell-item">
							<el-icon class="el-input__icon mr5"><ele-Message /></el-icon>
							{{ homePageData.email }}
						</div>
					</el-descriptions-item>
					<el-descriptions-item label="微信">
						<div class="cell-item">
							<SvgIcon name="my-weixin" :size="20" class="mr5"></SvgIcon>
							{{ homePageData.wechat }}
						</div>
					</el-descriptions-item>
					<el-descriptions-item label="QQ">
						<div class="cell-item">
							<SvgIcon name="fa fa-qq" :size="15" class="mr5"></SvgIcon>
							{{ homePageData.qq }}
						</div>
					</el-descriptions-item>
				</el-descriptions>
				<el-descriptions title="个人信息" direction="vertical" :column="column">
					<el-descriptions-item label="性别"> {{ homePageData.sex }}</el-descriptions-item>
					<el-descriptions-item label="行业"> {{ homePageData.industry }}</el-descriptions-item>
					<el-descriptions-item label="职位"> {{ homePageData.position }}</el-descriptions-item>
				</el-descriptions>
				<el-descriptions title="外部链接" direction="vertical" :column="column">
					<el-descriptions-item label="链接">
						<div class="cell-item" @click="onSkip">
							<el-icon class="mr5"> <ele-Share /> </el-icon>
							<span class="link"> 点击跳转</span>
						</div>
					</el-descriptions-item>
				</el-descriptions>
			</template>
			<el-empty v-else description="您暂未开启主页">
				<el-button type="primary" @click="onEdit" v-if="authUser.updateHomepage">开启主页</el-button>
			</el-empty>
		</el-card>
		<Edit ref="editRef" :home-page-data="homePageData" @refresh="getHomePageData"></Edit>
	</div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { homePageTypes, pageHomeParamsTypes } from '/@/api/homePage/types';
import { useRouter } from 'vue-router';
import { useHomePageApi } from '/@/api/homePage';
import { Session } from '/@/utils/storage';
import { ElMessage, ElMessageBox } from 'element-plus';
import { message } from '/@/utils/message';
import Edit from '/@/views/userCenter/homePage/component/edit.vue';
import useSubAuth from '/@/hooks/useSubAuth';

const router = useRouter();
const column = ref(3);
const editRef = ref();
const authUser = useSubAuth();
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '新增工单',
	submitTxt: '新增',
});

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
// 获取homepage数据
const getHomePageData = async (userID: string) => {
	if (!userID) {
		ElMessage.warning('抱歉，获取用户失败');
		return;
	}
	let data: pageHomeParamsTypes = {
		userID: userID,
	};
	setTimeout(() => {
		useHomePageApi()
			.getHomePage(data)
			.then((res: any) => {
				if (res.code === "SUCCESS") {
					homePageData.value = res.data.homepage;
          isOpen.value = true;
				} else {
          isOpen.value = false;
        }
			});
	}, 800);
};
// 修改
const onEdit = () => {
	editRef.value.openDialog();
};
//删除
const onDel = () => {
	ElMessageBox.confirm(`此操作将清空所有信息, 是否继续?`, '提示', {
		confirmButtonText: '清空',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useHomePageApi()
				.updateHomePage({ isDelete: false })
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						message('清空成功', { type: 'success' });
						getHomePageData(Session.get('userInfo').user.id);
					}
				});
		})
		.catch(() => {});
};

const onSkip = () => {
	if (!homePageData.value.link) {
		ElMessage.warning('抱歉，您没有设置外部链接');
		return;
	}
	router.push({
		path: '/externalLinkSkip',
		query: {
			link: homePageData.value.link,
		},
	});
};
onMounted(() => {
	const clientWidth = document.body.clientWidth;
	if (clientWidth < 1000) {
		column.value = 1;
	}
	getHomePageData(Session.get('userInfo').user.id);
});
</script>

<style scoped lang="scss">
.home-page-container {
	.cell-item {
		display: flex;
		align-items: center;
		font-size: 14px;
	}
	.link {
		cursor: pointer;
		color: var(--el-color-primary);
	}
}
</style>
