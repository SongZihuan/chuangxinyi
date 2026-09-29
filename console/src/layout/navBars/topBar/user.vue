<template>
	<div class="layout-navbars-breadcrumb-user pr15" :style="{ flex: layoutUserFlexNum }">
		<div class="layout-navbars-breadcrumb-user-icon" @click="onSearchClick" v-if="isH5">
			<el-icon title="菜单搜索">
				<ele-Search />
			</el-icon>
		</div>
		<el-dropdown :show-timeout="70" :hide-timeout="50" trigger="click" v-if="isH5" @command="onCommand">
			<div class="layout-navbars-breadcrumb-user-icon">账号</div>
			<template #dropdown>
				<el-dropdown-menu>
					<el-dropdown-item :command="item.path" v-for="(item, index) in accountList" :key="index" @click="onGoto">{{ item.name }}</el-dropdown-item>
				</el-dropdown-menu>
			</template>
		</el-dropdown>
		<el-dropdown :show-timeout="70" :hide-timeout="50" trigger="click" v-if="!isH5" @command="onCommand">
			<div class="layout-navbars-breadcrumb-user-icon">快捷导航</div>
			<template #dropdown>
				<el-dropdown-menu>
					<el-dropdown-item :command="item.path" v-for="(item, index) in routeList" :key="index">{{ item.name }}</el-dropdown-item>
				</el-dropdown-menu>
			</template>
		</el-dropdown>
		<!--    身份认证-->
		<div class="layout-navbars-breadcrumb-user-icon" @click="onGoto('/IdentityAuthentication')" v-if="isH5">实名认证</div>
		<div class="layout-navbars-breadcrumb-user-icon" @click="onGoto('/user')" v-if="isH5">我的钱包</div>
		<div class="layout-navbars-breadcrumb-user-icon" @click="onGoto('/workorder/user')" v-if="isH5">工单</div>
		<div class="layout-navbars-breadcrumb-user-icon" @click="onGoto('/audit/user')" v-if="isH5">审计</div>
		<div class="layout-navbars-breadcrumb-user-icon" ref="userNewsBadgeRef" v-click-outside="onUserNewsClick">
			<el-badge :is-dot="announcementList.length > 0 || messageList.length > 0">
				<el-icon title="消息">
					<ele-Bell />
				</el-icon>
			</el-badge>
		</div>
		<el-popover
			ref="userNewsRef"
			:virtual-ref="userNewsBadgeRef"
			placement="bottom"
			trigger="click"
			transition="el-zoom-in-top"
			virtual-triggering
			:width="300"
			:persistent="false"
		>
			<UserNews ref="userNewsListRef" />
		</el-popover>
		<div class="layout-navbars-breadcrumb-user-icon mr10" @click="onScreenfullClick" v-if="isH5">
			<SvgIcon :name="!state.isScreenfull ? 'my-screen' : 'my-off-screen'" :size="16" color="#61687c"></SvgIcon>
		</div>

		<div class="dropdown">
			<div class="dropdown-child">
				<div class="user-name">
					{{ userInfos.user.nickname || userInfos.user.userName || userInfos.user.phone || '异常用户' }}
				</div>

				<el-icon class="el-icon--right ml4">
					<ele-ArrowDown />
				</el-icon>
			</div>
			<div class="user">
				<div class="user-box">
					<div class="user-photo">
						<showImage :user="userInfos.user.id" />
					</div>
					<div class="user-head-right">
						<div>当前账号</div>
						<div v-if="userInfos && userInfos.user.phone && userInfos.role">
							{{ userInfos.user.nickname || userInfos.user.userName || userInfos.user.phone || '异常用户'
							}}<span v-if="userInfos.user.nickname">({{ userInfos.user.userName || userInfos.user.phone }})</span>-<span class="roleName">{{
								userInfos.role.name
							}}</span>
						</div>
					</div>
				</div>
				<!--  用户主页分享      -->
				<el-divider />
				<div class="user-list">
					<div class="user-item">
						<div class="user-label">用户ID</div>
						<div class="user-value">{{ userInfos.user.id }}</div>
						<el-tooltip class="box-item" effect="dark" content="点击复制" placement="top-start">
							<el-icon class="cursor-pointer ml4" @click="copyUserId(userInfos.user.id)">
								<ele-CopyDocument />
							</el-icon>
						</el-tooltip>
					</div>
					<div class="user-item">
						<div class="user-label">用户昵称</div>
						<div class="user-value">{{ userInfos.user.nickname ? userInfos.user.nickname : '暂未设置昵称' }}</div>
					</div>
					<div class="user-item">
						<div class="user-label">用户名</div>
						<div class="user-value">{{ userInfos.user.username ? userInfos.user.username : '暂未设置用户名' }}</div>
					</div>
					<div class="user-item">
						<div class="user-label">用户手机</div>
						<div class="user-value">{{ userInfos.user.phone }}</div>
					</div>
					<div class="user-item">
						<div class="user-label">用户邮箱</div>
						<div class="user-value">{{ userInfos.user.email ? userInfos.user.email : '暂未绑定邮箱' }}</div>
					</div>
				</div>
				<el-divider> </el-divider>
				<!--  用户主页分享      -->
				<div class="user-list">
					<div class="user-item">
						<div class="user-label">用户主页分享</div>
						<div class="user-value">
							<!--              复制分享链接-->
							<div class="btns">
								<div class="btn" @click="copyShare()">
									<el-icon :size="16" color="#61687c">
										<Share />
									</el-icon>
									复制分享链接
								</div>
							</div>
						</div>
					</div>
				</div>
				<el-divider> </el-divider>
				<div class="btn">
					<el-button plain @click="onHandleCommandClick('checkoutUser')">切换账号</el-button>
					<el-button plain @click="onHandleCommandClick('logOut')">退出登录</el-button>
				</div>
			</div>
		</div>

		<Search ref="searchRef" />
		<!-- 获取子账户   -->
		<subUserDialog ref="subUserDialogRef" type="user" />
	</div>
</template>

<script setup lang="ts" name="layoutBreadcrumbUser">
import { defineAsyncComponent, ref, unref, computed, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessageBox, ElMessage, ClickOutside as vClickOutside } from 'element-plus';
import screenfull from 'screenfull';
import { storeToRefs } from 'pinia';
import { useUserInfo } from '/@/stores/userInfo';
import showImage from '/@/components/showImage/index.vue';
import { Session, Local } from '/@/utils/storage';

import { Share } from '@element-plus/icons-vue';
import commonFunction from '/@/utils/commonFunction';
import { useSocketListInfo } from '/@/stores/socketListInfo';
import { useLoginApi } from '/@/api/login';
// 引入组件
const UserNews = defineAsyncComponent(() => import('/@/layout/navBars/topBar/userNews.vue'));
const Search = defineAsyncComponent(() => import('/@/layout/navBars/topBar/search.vue'));
const subUserDialog = defineAsyncComponent(() => import('/@/views/system/user/component/subUserDialog.vue'));

// 定义变量内容
const { copyText, isMobile } = commonFunction();
const userNewsRef = ref();
const userNewsListRef = ref();
const userNewsBadgeRef = ref();
const router = useRouter();
const stores = useUserInfo();
const isH5 = !isMobile();
const socketInfo = useSocketListInfo();
const { messageList, announcementList } = storeToRefs(socketInfo);
const { userInfos } = storeToRefs(stores) as any;
const searchRef = ref();
const subUserDialogRef = ref();
const state = reactive({
	isScreenfull: false,
	disabledSize: 'large',
});

const accountList = ref([
	{ name: '本账号', path: '/home' },
	{ name: '子账号', path: '/accountManagement/subUser/list' },
]);
const routeList = ref([
	{ name: '工单', path: '/workorder/user' },
	{ name: '审计', path: '/audit/user' },
	{ name: '实名认证', path: '/IdentityAuthentication' },
	{ name: '我的钱包', path: '/user' },
	{ name: '我的钱包', path: '/user' },
]);
const onCommand = (path: string) => {
	onGoto(path);
};
// 设置分割样式
const layoutUserFlexNum = computed(() => {
	return '1';
});

// 全屏点击时
const onScreenfullClick = () => {
	if (!screenfull.isEnabled) {
		ElMessage.warning('暂不不支持全屏');
		return false;
	}
	screenfull.toggle();
	screenfull.on('change', () => {
		if (screenfull.isFullscreen) state.isScreenfull = true;
		else state.isScreenfull = false;
	});
};
// 消息通知点击时
const onUserNewsClick = () => {
	unref(userNewsRef).popperRef?.delayHide?.();
};
const copyUserId = (val: string) => {
	copyText(val);
};
// 复制主页分享链接
const copyShare = () => {
	const userInfos = Session.get('userInfo');
	const input = document.createElement('input');
	input.setAttribute('readonly', 'readonly');
	input.setAttribute('value', new URL(window.location.href).origin + `/visitor?userID=${userInfos.user.id}`);
	document.body.appendChild(input);
	input.select();
	input.setSelectionRange(0, 9999);
	if (document.execCommand('copy')) {
		document.execCommand('copy');
		ElMessage.success('复制成功');
	}
	document.body.removeChild(input);
};

// 下拉菜单点击时
const onHandleCommandClick = (path: string) => {
	if (path === 'logOut') {
		ElMessageBox({
			closeOnClickModal: false,
			closeOnPressEscape: false,
			title: '提示',
			message: '此操作将退出登录, 是否继续?',
			showCancelButton: true,
			confirmButtonText: '确定',
			cancelButtonText: '取消',
			buttonSize: 'default',
			beforeClose: (action, instance, done) => {
				if (action === 'confirm') {
					instance.confirmButtonLoading = true;
					instance.confirmButtonText = '退出中';
					setTimeout(() => {
						done();
						setTimeout(() => {
							instance.confirmButtonLoading = false;
						}, 300);
					}, 700);
				} else {
					done();
				}
			},
		})
			.then(async () => {
				useLoginApi()
					.logout()
					.then(() => {
						// 清除LocalStorage中的hierarchy
						// 清除缓存/token等
						Session.clear();
						Local.clear();
						window.location.href = '/login'; // 回到登录页
					});
			})
			.catch(() => {});
	} else if (path === 'checkoutUser') {
		let userInfos = Session.get('userInfo');
		let token = Session.get('token');
		// 深拷贝
		const subAuth = JSON.parse(JSON.stringify(Session.get('userType')));
		if (userInfos && userInfos.user && userInfos.user.id && token) {
			subUserDialogRef.value.openDialog({ id: userInfos.user.id, data: { ...userInfos.user, token: token, subType: subAuth } });
		}
	} else {
		router.push(path);
	}
};
// 菜单搜索点击
const onSearchClick = () => {
	searchRef.value.openSearch();
};

const onGoto = (path: string) => {
	router.push(path);
};
// 初始化组件大小/i18n
const initI18nOrSize = (value: string, attr: string) => {
	(<any>state)[attr] = Local.get('themeConfig')[value];
};
// 页面加载时
onMounted(() => {
	if (Local.get('themeConfig')) {
		initI18nOrSize('globalComponentSize', 'disabledSize');
	}
});
</script>

<style scoped lang="scss">
@import '../../../theme/mixins/index.scss';
.layout-navbars-breadcrumb-user {
	display: flex;
	align-items: center;
	justify-content: flex-end;
	background: #fff;
	&-link {
		height: 100%;
		display: flex;
		align-items: center;
		white-space: nowrap;
		&-photo {
			width: 25px;
			height: 25px;
			border-radius: 100%;
		}
	}
	&-icon {
		padding: 0 10px;
		cursor: pointer;
		color: var(--next-bg-topBarColor);
		height: 50px;
		line-height: 50px;
		display: flex;
		align-items: center;
		&:hover {
			background: var(--next-color-user-hover);
			i {
				display: inline-block;
				animation: logoAnimation 0.3s ease-in-out;
			}
		}
	}
	:deep(.el-dropdown) {
		color: var(--next-bg-topBarColor);
	}
	:deep(.el-badge) {
		height: 40px;
		line-height: 40px;
		display: flex;
		align-items: center;
	}
	:deep(.el-badge__content.is-fixed) {
		top: 12px;
	}
}
.dropdown {
	position: relative;
	&-child {
		min-width: 50px;
		display: flex;
		flex-direction: row;
		cursor: pointer;
		height: 50px;
		align-items: center;
		.user-name {
			max-width: 200px;
			height: 20px;
			@include text-ellipsis(1);
		}
	}
	&:hover .user {
		display: block;
	}

	.user {
		width: 300px;
		min-height: 150px;
		padding: 15px;
		background-color: #fff;
		display: none;
		position: absolute;
		top: 50px;
		right: -15px;
		z-index: 2004;
		box-shadow: 0 8px 16px #0000001f;
		&:hover {
			display: block;
		}
		.user-box {
			display: flex;
			flex-direction: row;
		}
		.user-photo {
			width: 36px;
			height: 36px;
			display: flex;
		}
		.user-head-right {
			margin-left: 16px;
			min-height: 36px;
			cursor: pointer;
			:nth-child(1) {
				font-size: 14px;
			}

			:nth-child(2) {
				font-size: 14px;
				font-weight: 400;
				color: #414960;
			}

			.roleName {
				color: var(--el-text-color-secondary);
			}
		}
		.user-list {
			display: flex;
			flex-direction: column;
			color: #555;
			font-size: 12px;
			.user-item {
				display: flex;
				flex-direction: row;
				align-items: center;
				.user-label {
					min-width: 50px;
					margin-right: 8px;
					color: #999;
				}
				.user-value {
					flex: 1;
					@include text-ellipsis(1);
					.btns {
						display: flex;
						flex-direction: row;
						justify-content: space-between;
						.btn {
							display: flex;
							flex-direction: row;
							justify-content: center;
							align-items: center;
							color: #61687c;
							font-size: 12px;
							border: 1px solid #ebeef5;
							padding: 5px 10px;
							cursor: pointer;
						}
						.btn:hover {
							border-color: #409eff;
						}
					}
				}
			}
		}
	}
	.user-title {
		font-weight: 600;
		color: #61687c;
		margin-bottom: 8px;
	}
	.btn {
		display: flex;
		flex-direction: row;
		justify-content: center;
	}
}
:deep(.el-divider--horizontal) {
	margin: 10px 0px;
}
</style>
