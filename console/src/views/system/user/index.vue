<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="system-user-search mb15">
				<el-form :inline="true" :model="search">
					<el-form-item label="">
						<el-date-picker
							v-model="search.range"
							type="datetimerange"
							range-separator="至"
							start-placeholder="开始日期"
							end-placeholder="结束日期"
						/>
					</el-form-item>
					<el-form-item>
						<el-button type="primary" @click="handleSearch" class="mb10">
							<el-icon>
								<ele-Search />
							</el-icon>
							查询
						</el-button>
						<el-button type="primary" @click="openDialog('avatarUpload')" class="mb10">
							<el-icon>
								<ele-Edit />
							</el-icon>
							默认头像
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<el-table
				:data="state.tableData.data"
				v-loading="state.tableData.loading"
				@selection-change="handleSelectionChange"
				@cell-click="cellClick"
				border
			>
				<el-table-column prop="numberID" label="用户数字ID" show-overflow-tooltip width="120" align="left"></el-table-column>
				<el-table-column prop="id" label="用户ID" show-overflow-tooltip align="left" width="100"></el-table-column>
				<el-table-column prop="phone" label="手机号" show-overflow-tooltip min-width="150" align="center"></el-table-column>
				<el-table-column prop="userRealName" label="真实姓名" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column label="头像" width="100" align="center">
					<template v-slot="scope">
						<div class="center">
							<el-image
								:style="{ width: `50px`, height: `50px`, borderRadius: '4px' }"
								:src="showImageCalc(scope.row.id)"
								:zoom-rate="1.2"
								:preview-src-list="[showImageCalc(scope.row.id)]"
								preview-teleported
								fit="cover"
								close-on-press-escape
							/>
						</div>
					</template>
				</el-table-column>
				<el-table-column prop="wechatHeader" label="微信头像" width="100" align="center">
					<template v-slot="scope">
						<div class="center">
							<el-image
								:style="{ width: `50px`, height: `50px`, borderRadius: '6px' }"
								:src="scope.row.wechatHeader"
								:zoom-rate="2"
								:preview-src-list="[scope.row.wechatHeader]"
								preview-teleported
								v-if="scope.row.wechatHeader"
								fit="cover"
								close-on-press-escape
							/>
						</div>
					</template>
				</el-table-column>
				<el-table-column prop="wechatNickName" label="微信昵称" show-overflow-tooltip width="100" align="center"></el-table-column>
				<el-table-column prop="companyName" label="公司名称" show-overflow-tooltip min-width="150" align="center"></el-table-column>
				<el-table-column prop="roleName" label="角色" show-overflow-tooltip align="center" width="100"></el-table-column>
				<el-table-column prop="status" label="单点登录" show-overflow-tooltip width="90" align="center">
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.signin">是</el-tag>
						<el-tag type="info" v-else>否</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="status" label="用户状态" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-tag :type="userColorTypesEnumData[scope.row.status]">{{ UserStatusEnumData[scope.row.status] }}</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="createTime" label="创建时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="160" fixed="right" align="center">
					<template #default="scope">
						<div class="menus-dropdown">
							<el-dropdown>
								<span class="el-dropdown-link">
									查看
									<el-icon class="el-icon--right">
										<arrow-down />
									</el-icon>
								</span>
								<template #dropdown>
									<el-dropdown-menu>
										<el-dropdown-item @click="onOpenUserInfo(scope.row)">用户详情</el-dropdown-item>
										<el-dropdown-item @click="onOpenUserCenter(scope.row)">用户主页</el-dropdown-item>
										<el-dropdown-item @click="onOpenFundUserInfo(scope.row)">用户信息</el-dropdown-item>
										<el-dropdown-item @click="onOpenFinance(scope.row)"> 资金信息</el-dropdown-item>
										<el-dropdown-item @click="onOpenUserrealName(scope.row)">实名信息</el-dropdown-item>
										<el-dropdown-item @click="onOpenSubUser(scope.row)"> 子账户</el-dropdown-item>
										<el-dropdown-item @click="onOpenDiscount(scope.row)">优惠券</el-dropdown-item>
										<el-dropdown-item @click="onPayProcess(scope.row)">充值记录</el-dropdown-item>
										<el-dropdown-item @click="onPayHistory(scope.row)">消费记录</el-dropdown-item>
										<el-dropdown-item @click="onOpenToken(scope.row)">用户在线列表</el-dropdown-item>
										<el-dropdown-item @click="onOpenSubToken(scope.row)">子用户在线列表</el-dropdown-item>
										<el-dropdown-item @click="onOpenOauth2(scope.row)">授权在线列表</el-dropdown-item>
										<el-dropdown-item @click="onOpenOauth2Record(scope.row)">授权记录</el-dropdown-item>
										<el-dropdown-item @click="onOpenBanned(scope.row)">开通列表</el-dropdown-item>
										<el-dropdown-item @click="onInvite(scope.row)">邀请用户列表</el-dropdown-item>
										<el-dropdown-item @click="onOpenAudit(scope.row)">查看审计内容</el-dropdown-item>
									</el-dropdown-menu>
								</template>
							</el-dropdown>
							<el-dropdown>
								<span class="el-dropdown-link">
									操作
									<el-icon class="el-icon--right">
										<arrow-down />
									</el-icon>
								</span>
								<template #dropdown>
									<el-dropdown-menu>
										<el-dropdown-item @click="onCopyID(scope.row)">复制ID</el-dropdown-item>
										<el-dropdown-item @click="onOpenBindRole(scope.row)">角色绑定</el-dropdown-item>
										<el-dropdown-item @click="onCreateWorkOrder(scope.row)">创建工单</el-dropdown-item>
										<el-dropdown-item @click="onSend(scope.row, 2)"> 发送短信</el-dropdown-item>
										<el-dropdown-item @click="onSend(scope.row, 3)"> 发送邮件</el-dropdown-item>
										<el-dropdown-item @click="onSend(scope.row, 4)"> 发送服务号</el-dropdown-item>
										<el-dropdown-item @click="onSend(scope.row, 1)"> 发送站内信</el-dropdown-item>
										<el-dropdown-item @click="onSend(scope.row, 5)"> 发送企微机器人信息</el-dropdown-item>
										<el-dropdown-item @click="onAdd(scope.row)">添加用户订单</el-dropdown-item>
										<el-dropdown-item @click="onAddDefary(scope.row)">添加用户消费</el-dropdown-item>
										<el-dropdown-item @click="onBilledAdd(scope.row)">增加开票额度</el-dropdown-item>
										<el-dropdown-item @click="onsubBilled(scope.row)">减少开票额度</el-dropdown-item>
									</el-dropdown-menu>
								</template>
							</el-dropdown>
						</div>
					</template>
				</el-table-column>
			</el-table>
			<el-pagination
				hide-on-single-page
				@size-change="onHandleSizeChange"
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				:page-sizes="[10, 20, 30]"
				v-model:current-page="search.page"
				background
				v-model:page-size="search.pagesize"
				layout="total, sizes, prev, pager, next, jumper"
				:total="state.tableData.total"
			>
			</el-pagination>
		</el-card>
		<!-- 用户主页 -->
		<UserInfo ref="userInfoDialogRef" @refresh="getTableData()" />
		<!-- 用户实名信息 -->
		<userRealName ref="userRealNameRef" @refresh="getTableData()" />
		<!-- 用户信息 -->
		<UserDialog ref="userDialogRef" @refresh="getTableData()" />
		<!-- 绑定角色 -->
		<roleDialog ref="roleRef" @refresh="getTableData()" />
		<!-- 获取子账户   -->
		<subUserDialog ref="subUserDialogRef" />
		<!-- 查看优惠 -->
		<discountDislog ref="discountDislogRef" />
		<!-- 发送站内信 -->
		<message ref="messageRef" />
		<!-- 发送短信 -->
		<sms ref="smsRef" />
		<!-- 发送邮件 -->
		<email ref="emailRef" />
		<!-- 发送服务号 -->
		<fuwuhao ref="fuwuhaoRef" />
		<!-- 发送企业微信机器人 -->
		<wxrobotMessage ref="wxrobotMessageRef" />
		<!-- 用户充值 -->
		<payDialog ref="payDialogRef" />
		<!-- 用户充值记录 -->
		<payHistpry ref="payHistpryRef" />
		<!-- 查询订单 -->
		<searchOrder ref="searchOrderRef" />
		<!--  创建工单  -->
		<action-dialog ref="actionDialogRef" />
		<!-- 上传默认头像 -->
		<defaultAvatar ref="defaultAvatarRef" @refresh="getTableData()" />
		<!--  用户主页  -->
		<userCenter ref="userCenterDialogRef" />
		<!-- 用户信息 -->
		<fundUserInfoDialog ref="fundUserInfoDialogRef" @refresh="getTableData()" />
		<!-- 查看用户资产 -->
		<financeDialog ref="financeDialogRef" />
		<!-- 审计列表 -->
		<auditDialog ref="auditDialogRef" />
		<!-- 增加余额 -->
		<addBlance ref="addBlanceRef" />
		<!-- 减少金额 -->
		<addDefary ref="addDefaryRef" />
		<!-- 增加可开票 -->
		<billedAdd ref="billedAddRef" />
		<!-- 核销开票 -->
		<subBilled ref="subBilledRef" />
		<!-- 查看用户在线token -->
		<tokenList ref="tokenListRef" />
		<!--    查看子用户在线token-->
		<subTokenList ref="subTokenListRef" />
		<!--   查看授权列表-->
		<oauthAuthorization ref="oauthAuthorizationRef" />
		<!--   查看授权记录-->
		<oauthAuthorizationRecord ref="oauthAuthorizationRecordRef" />
		<!-- 封禁列表 -->
		<oauthAuthorizationBanned ref="oauthAuthorizationBannedRef" />
		<!--  邀请列表  -->
		<inviter ref="inviteListRef" />
	</div>
</template>

<script setup lang="ts" name="systemUser">
import { defineAsyncComponent, reactive, onMounted, ref, watch } from 'vue';
import { userApi } from '/@/api/system/user/index';
import type { searchTypes } from './types';
import { ArrowDown } from '@element-plus/icons-vue';
import dayjs from 'dayjs';
import UserInfo from '/@/views/system/user/component/userCenter.vue';
import { userColorTypes, UserStatusEnum, userTypes } from '/@/data/enum';
import ActionDialog from '/@/views/workOrder/admin/component/actionDialog.vue';

const userApiCollect = userApi();
// 引入组件
const UserDialog = defineAsyncComponent(() => import('/@/views/system/user/component/dialog.vue'));
const roleDialog = defineAsyncComponent(() => import('./component/bindRole.vue'));
const subUserDialog = defineAsyncComponent(() => import('./component/subUserDialog.vue'));
const discountDislog = defineAsyncComponent(() => import('./component/discountDislog.vue'));
const message = defineAsyncComponent(() => import('./component/message.vue'));
const sms = defineAsyncComponent(() => import('./component/sms.vue'));
const email = defineAsyncComponent(() => import('./component/email.vue'));
const fuwuhao = defineAsyncComponent(() => import('./component/fuwuhao.vue'));
const wxrobotMessage = defineAsyncComponent(() => import('./component/wxrobotMessage.vue'));
const payDialog = defineAsyncComponent(() => import('./component/payDialog.vue'));
const payHistpry = defineAsyncComponent(() => import('./component/payHistpry.vue'));
const searchOrder = defineAsyncComponent(() => import('./component/searchOrder.vue'));
const userRealName = defineAsyncComponent(() => import('./component/userRealName.vue'));
const defaultAvatar = defineAsyncComponent(() => import('./component/defaultAvatar.vue'));
const userCenter = defineAsyncComponent(() => import('./component/userInfoCenter.vue'));
const fundUserInfoDialog = defineAsyncComponent(() => import('./component/fundUserInfo.vue'));
const financeDialog = defineAsyncComponent(() => import('./component/financeDialog.vue'));
const auditDialog = defineAsyncComponent(() => import('./component/auditDialog.vue'));
const addBlance = defineAsyncComponent(() => import('./component/addBlance.vue'));
const addDefary = defineAsyncComponent(() => import('./component/addDefary.vue'));
const billedAdd = defineAsyncComponent(() => import('./component/billedAdd.vue'));
const subBilled = defineAsyncComponent(() => import('./component/subBilled.vue'));
const tokenList = defineAsyncComponent(() => import('./component/oauth2/tokenList.vue'));
const subTokenList = defineAsyncComponent(() => import('./component/oauth2/subTokenList.vue'));
const oauthAuthorization = defineAsyncComponent(() => import('./component/oauth2/oauthAuthorization.vue'));
const oauthAuthorizationRecord = defineAsyncComponent(() => import('./component/oauth2/oauthAuthorizationRecord.vue'));
const oauthAuthorizationBanned = defineAsyncComponent(() => import('./component/oauth2/oauthAuthrizationBanned.vue'));
const inviter = defineAsyncComponent(() => import('./component/inviter.vue'));
import commonFunction from '/@/utils/commonFunction';

// 定义变量内容
const { copyText } = commonFunction();
// 定义变量内容
const userInfoDialogRef = ref();
const userDialogRef = ref();
const discountDislogRef = ref();
const messageRef = ref();
const roleRef = ref();
const subUserDialogRef = ref();
const smsRef = ref();
const emailRef = ref();
const fuwuhaoRef = ref();
const wxrobotMessageRef = ref();
const payDialogRef = ref();
const payHistpryRef = ref();
const searchOrderRef = ref();
const userRealNameRef = ref();
const actionDialogRef = ref();
const defaultAvatarRef = ref();
const userCenterDialogRef = ref();
const fundUserInfoDialogRef = ref();
const inviteListRef = ref();
const financeDialogRef = ref();
const auditDialogRef = ref();
const addBlanceRef = ref();
const addDefaryRef = ref();
const billedAddRef = ref();
const subBilledRef = ref();
const tokenListRef = ref();
const subTokenListRef = ref();
const oauthAuthorizationRef = ref();
const oauthAuthorizationRecordRef = ref();
const oauthAuthorizationBannedRef = ref();
const search = reactive<searchTypes>({
	range: [],
	page: 1,
	pagesize: 20,
	starttime: null,
	endtime: null,
});
const state = reactive({
	tableData: {
		data: [] as any[],
		total: 0,
		loading: false,
	},
});
const UserStatusEnumData: userTypes = UserStatusEnum;
const userColorTypesEnumData: userTypes = userColorTypes;
const onSend = (row: any, type: number) => {
	if (type === 1) {
		messageRef.value.openDialog(row);
	}
	if (type === 2) {
		smsRef.value.openDialog(row);
	}
	if (type === 3) {
		emailRef.value.openDialog(row);
	}
	if (type === 4) {
		fuwuhaoRef.value.openDialog(row);
	}
	if (type === 5) {
		wxrobotMessageRef.value.openDialog(row);
	}
};
const onCopyID = (row: { id: string }) => {
	copyText(row.id);
};
const onOpenToken = (row: any) => {
	tokenListRef.value.openDialog(row);
};
const onOpenOauth2 = (row: any) => {
	oauthAuthorizationRef.value.openDialog(row);
};
const onOpenBanned = (row: any) => {
	oauthAuthorizationBannedRef.value.openDialog(row);
};
const onOpenOauth2Record = (row: any) => {
	oauthAuthorizationRecordRef.value.openDialog(row);
};
const onOpenSubToken = (row: any) => {
	subTokenListRef.value.openDialog(row);
};
const onAdd = (row: any) => {
	addBlanceRef.value.openDialog(row);
};
const onAddDefary = (row: any) => {
	addDefaryRef.value.openDialog(row);
};
const onBilledAdd = (row: any) => {
	billedAddRef.value.openDialog(row);
};
const onsubBilled = (row) => {
	subBilledRef.value.openDialog(row);
};
const onInvite = (row: any) => {
	inviteListRef.value.openDialog(row);
};
const onPayProcess = (row: any) => {
	payDialogRef.value.openDialog(row);
};
const onPayHistory = (row: any) => {
	payHistpryRef.value.openDialog(row);
};
const onCreateWorkOrder = (row: any) => {
	actionDialogRef.value.openDialog(row, 'edit');
};
const showImageCalc = (id: string) => {
	return import.meta.env.VITE_API_URL + '/public/header/user?id=' + id + `&?date=${new Date().getTime()}`;
};
const openDialog = (type: string) => {
	if (type === 'avatarUpload') {
		defaultAvatarRef.value.openDialog();
	}
};
const onOpenAudit = (row: any) => {
	auditDialogRef.value.openDialog(row);
};
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	userApiCollect.getuserList(search).then((res: any) => {
		if (res.code === 'SUCCESS') {
			state.tableData.data = res.data.user;
			state.tableData.total = res.data.count;
			state.tableData.loading = false;
		}
	});
};
const onOpenUserrealName = (row: any) => {
	userRealNameRef.value.openDialog(row);
};

// 用户信息
const onOpenUserInfo = (row: any) => {
	userInfoDialogRef.value.openDialog(row);
};
// 查看用户资金
const onOpenFinance = (row: any) => {
	financeDialogRef.value.openDialog(row);
};
// 用户首页
const onOpenUserCenter = (row: any) => {
	userCenterDialogRef.value.openDialog(row);
};
// 用户个人信息
const onOpenFundUserInfo = (row: any) => {
	fundUserInfoDialogRef.value.openDialog(row);
};
const onOpenBindRole = (row: any) => {
	roleRef.value.openDialog(row);
};
const onOpenDiscount = (row: any) => {
	discountDislogRef.value.openDialog(row);
};
// 打开子账户信息
const onOpenSubUser = (row: any) => {
	subUserDialogRef.value.openDialog(row);
};
const multipleSelection = ref<number[]>([]);
const handleSelectionChange = (val: number[]) => {
	multipleSelection.value = val;
};

const currentInfo = ref<any>({});
const cellClick = (row: any) => {
	currentInfo.value = row;
};

//搜素
const handleSearch = () => {
	search.page = 1;
	getTableData();
};
const onSearchOrder = (type: string) => {
	searchOrderRef.value.openDialog(type);
};
// 分页改变
const onHandleSizeChange = (val: number) => {
	search.pagesize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	search.page = val;
	getTableData();
};
watch(
	() => search.range,
	() => {
		if (search.range.length > 0) {
			search.starttime = dayjs(search.range[0]).unix();
			search.endtime = dayjs(search.range[1]).unix();
		}
	},
	{
		deep: true,
	}
);
// 页面加载时
onMounted(async () => {
	await getTableData();
});
</script>

<style scoped lang="scss">
// el-dropdown添加滚动条
.el-dropdown-menu {
	max-height: 300px;
	overflow-y: auto;
}
.center {
	display: flex;
	align-items: center;
	justify-content: center;
}

.container {
	:deep(.el-card__body) {
		display: flex;
		flex-direction: column;
		flex: 1;
		overflow: auto;

		.el-table {
			flex: 1;
		}
	}

	.el-dropdown-link {
		text-align: right;
		cursor: pointer;
		line-height: 16px;
		color: var(--el-color-primary);
		display: flex;
		align-items: center;
	}

	.menus-dropdown {
		padding: 0 10px;
		display: flex;
		justify-content: space-between;
	}
}
</style>
