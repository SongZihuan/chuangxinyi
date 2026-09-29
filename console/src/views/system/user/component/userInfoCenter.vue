<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="869px" @close="closeDialog" destroy-on-close>
		<div class="home-page-container container">
			<el-card shadow="hover">
				<template v-if="isOpen">
          <el-descriptions title="用户主页" direction="vertical" :column="column">
            <el-descriptions-item>
              <div class="cell-item" @click="goHome">
                <el-icon class="mr5"><ele-House /></el-icon>
                <span class="link">点击查看</span>
              </div>
            </el-descriptions-item>
          </el-descriptions>
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
				</template>
				<el-empty v-else description="您暂未开启主页"> </el-empty>
			</el-card>
		</div>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { homePageTypes, pageHomeParamsTypes } from '/@/api/homePage/types';
import { useRouter } from 'vue-router';
import { useHomePageApi } from '/@/api/homePage';
import { ElMessage } from 'element-plus';

const router = useRouter();
const column = ref(3);
const isOpen = ref(true);
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
const id = ref('');
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '用户主页',
	submitTxt: '新增',
});
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
			if (res.code === "SUCCESS") {
				isOpen.value = true;
				homePageData.value = res.data.homepage;
				return;
			} else {
        isOpen.value = false;
      }
		});
};
const openDialog = (row: any) => {
	if (!row.id) {
		ElMessage.warning('用户获取失败');
		return;
	}
  id.value = row.id;
	getHomePageData(row.id as string);
	dialog.isShowDialog = true;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
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
const goHome = () => {
  if (!id.value) {
    ElMessage.warning('抱歉，您暂未开始主页');
    return;
  }
  window.open(new URL(window.location.href).origin + `/visitor?userID=${id.value}`, '_blank');
};
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style scoped lang="scss">
//兼容pc和移动端
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
