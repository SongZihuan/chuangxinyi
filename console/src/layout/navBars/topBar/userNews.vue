<template>
	<div class="layout-navbars-breadcrumb-user-news">
		<div class="head-box">
			<el-tabs v-model="tabsActiveName">
				<el-tab-pane label="站内信" name="first"></el-tab-pane>
				<el-tab-pane label="公告" name="second"></el-tab-pane>
			</el-tabs>
		</div>
		<div class="content-box">
			<template v-if="announcementList.length > 0 && tabsActiveName == 'second'">
				<div class="content-box-item" v-for="(v, k) in announcementList" :key="k">
					<div class="content-box-title">
						<div class="svgicon">
							<SvgIcon name="my-announcement" :size="20" color="#e6a23c"></SvgIcon>
						</div>
						<div>{{ v.title }}</div>
					</div>
					<div class="content-box-time">{{ dayjs.unix(v.startAt).format('YYYY-MM-DD') }}至{{ dayjs.unix(v.stopAt).format('YYYY-MM-DD') }}</div>
					<div class="content-box-msg" @click.stop="showContent(v)">查看</div>
				</div>
				<div class="all" @click="goAnnouncement">前往公告中心</div>
			</template>
			<template v-else-if="messageList.length > 0 && tabsActiveName == 'first'">
				<div class="content-box-item" v-for="(v, k) in messageList" :key="k" @click="goMessage">
					<div class="content-box-title">
						<div class="svgicon2">
							<SvgIcon name="my-message" :size="20" color="#f56c6c"></SvgIcon>
						</div>
						<div>{{ v.title }}</div>
					</div>
					<div class="content-box-time">{{ v.content }}</div>
					<div class="content-box-msg" @click.stop="handleView(v)">查看</div>
				</div>
				<div class="all" @click="goMessage">前往站内信中心</div>
			</template>
			<el-empty description="暂无数据" v-else></el-empty>
		</div>
		<Detail-dialog ref="detailDialogRef"></Detail-dialog>
		<ShowHtml ref="showHtmlRef" :title="currentInfo.title" :contentHtml="currentInfo.content" :subtitle="currentInfo.subtitle" />
	</div>
</template>

<script setup lang="ts" name="layoutBreadcrumbUserNews">
import { ref } from 'vue';
import dayjs from 'dayjs';
import { useRouter } from 'vue-router';
import ShowHtml from '/@/components/ShowHtml/index.vue';
import DetailDialog from '/@/views/message/user/component/detailDialog.vue';
import { useSocketListInfo } from '/@/stores/socketListInfo';
import { storeToRefs } from 'pinia';
const showHtmlRef = ref();
const detailDialogRef = ref();
const tabsActiveName = ref('first');
const router = useRouter();

const socketInfo = useSocketListInfo();
const { messageList, announcementList } = storeToRefs(socketInfo);
// eslint-disable-next-line vue/no-dupe-keys
const currentInfo = ref({ title: '', content: '', subtitle: '' });
const showContent = (row: any) => {
	currentInfo.value = row;
	currentInfo.value.subtitle = dayjs.unix(row.startAt).format('YYYY-MM-DD') + '至' + dayjs.unix(row.stopAt).format('YYYY-MM-DD');
	showHtmlRef.value.openDialog();
};
const goMessage = () => {
	router.push({ path: '/message/user' });
};
const goAnnouncement = () => {
	router.push({ path: '/announcement/user' });
};
const handleView = (row: any) => {
	detailDialogRef.value.openDialog(row);
};
</script>

<style scoped lang="scss">
.all {
	width: 100%;
	text-align: center;
	color: var(--el-color-primary);
	cursor: pointer;
}
.layout-navbars-breadcrumb-user-news {
	.head-box {
		width: 100%;
		display: flex;
		box-sizing: border-box;
		color: var(--el-text-color-primary);
		justify-content: space-between;
		height: 35px;
		align-items: center;
		.head-box-btn {
			color: var(--el-color-primary);
			font-size: 13px;
			cursor: pointer;
			opacity: 0.8;
			&:hover {
				opacity: 1;
			}
		}
	}
	.content-box {
		font-size: 13px;
		.content-box-item {
			padding-top: 12px;
			cursor: pointer;
			&:last-of-type {
				padding-bottom: 12px;
			}
			.content-box-title {
				display: flex;
				flex-direction: row;
				align-items: center;
				font-weight: bold;
				.svgicon {
					width: 32px;
					height: 32px;
					background: var(--next-color-warning-lighter);
					border-radius: 50%;
					display: flex;
					flex-direction: row;
					align-items: center;
					justify-content: center;
					margin-right: 6px;
				}
				.svgicon2 {
					width: 32px;
					height: 32px;
					background: var(--next-color-danger-lighter);
					border-radius: 50%;
					display: flex;
					flex-direction: row;
					align-items: center;
					justify-content: center;
					margin-right: 6px;
				}
			}
			.content-box-msg {
				color: var(--el-color-primary);
				margin-top: 5px;
				margin-bottom: 5px;
				text-align: right;
				cursor: pointer;
			}
			.content-box-time {
				color: var(--el-text-color-secondary);
				margin-top: 5px;
			}
		}
	}
	:deep(.el-empty__description p) {
		font-size: 13px;
	}
}
</style>
