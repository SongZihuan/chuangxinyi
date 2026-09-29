<template>
	<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="869px" @close="closeDialog" destroy-on-close>
		<div class="mt20">
			<el-card>
				<div class="main">
					<div class="user-info">
						<div class="user-info-left">
							<div class="user-head">
								<div style="width: 60px; height: 60px">
									<showImage :user="userInfo.id" />
								</div>
								<div class="user-head-right">
									<div>账号</div>
									<div v-if="userInfo && userInfo.phone">
										{{ userInfo.nickname || userInfo.userName || userInfo.phone || '异常用户' }}
									</div>
								</div>
							</div>
						</div>
					</div>
					<el-descriptions :column="column" border>
						<el-descriptions-item>
							<template #label>
								<div class="cell-item">
									<el-icon class="mr4">
										<ele-Avatar />
									</el-icon>
									微信头像
								</div>
							</template>
							<el-image
								:style="{ width: `40px`, height: `40px`, borderRadius: '6px' }"
								:src="userInfo.wxHeader"
								:zoom-rate="2"
								:preview-src-list="[userInfo.wxHeader]"
								preview-teleported
								v-if="userInfo.wxHeader"
								fit="cover"
								close-on-press-escape
							/>
						</el-descriptions-item>
						<el-descriptions-item>
							<template #label>
								<div class="cell-item">
									<SvgIcon name="my-weixin" :size="14" class="mr5" color="#606266"></SvgIcon>
									微信名称
								</div>
							</template>
							{{ userInfo.wxNickName }}
						</el-descriptions-item>
						<el-descriptions-item>
							<template #label>
								<div class="cell-item">
									<el-icon class="mr4">
										<ele-User />
									</el-icon>
									用户ID
								</div>
							</template>
							{{ userInfo.id }}
							<el-tooltip class="box-item" effect="dark" content="点击复制" placement="top-start">
								<el-icon class="ml4 cursor-pointer" @click="copyUserId(userInfo.id)">
									<ele-CopyDocument />
								</el-icon>
							</el-tooltip>
						</el-descriptions-item>
						<el-descriptions-item>
							<template #label>
								<div class="cell-item">
									<el-icon class="mr4">
										<ele-User />
									</el-icon>
									用户名称
								</div>
							</template>
							{{ userInfo.userName }}
						</el-descriptions-item>
						<el-descriptions-item>
							<template #label>
								<div class="cell-item">
									<el-icon class="mr4">
										<ele-User />
									</el-icon>
									用户昵称
								</div>
							</template>
							{{ userInfo.nickname }}
						</el-descriptions-item>
						<el-descriptions-item>
							<template #label>
								<div class="cell-item">
									<el-icon class="mr4">
										<ele-Iphone />
									</el-icon>
									手机号码
								</div>
							</template>
							{{ userInfo.phone }}
						</el-descriptions-item>
						<el-descriptions-item>
							<template #label>
								<div class="cell-item">
									<el-icon class="el-input__icon"><ele-Message /></el-icon>
									邮箱
								</div>
							</template>
							{{ userInfo.email }}
						</el-descriptions-item>
						<el-descriptions-item label="注册时间">
							<template #label>
								<div class="cell-item">
									<SvgIcon name="my-peoples" :size="14" color="#606266"></SvgIcon>
									注册时间
								</div>
							</template>
							{{ dayjs.unix(userInfo.createAt).format('YYYY-MM-DD HH:mm:ss') || '-' }}
						</el-descriptions-item>
						<el-descriptions-item>
							<template #label>
								<div class="cell-item">
									<SvgIcon name="my-status" :size="14" color="#606266"></SvgIcon>
									状态
								</div>
							</template>
							<el-text class="mx-1" :type="userColorEnumTypes[userInfo.status]" size="mini">
								{{ UserStatusEnumTypes[userInfo.status] }}
							</el-text>
						</el-descriptions-item>
					</el-descriptions>
				</div>
			</el-card>
		</div>
	</el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import dayjs from 'dayjs';
import { userColorTypes, UserStatusEnum, userTypes } from '/@/data/enum';
import showImage from '/@/components/showImage/index.vue';
import commonFunction from '/@/utils/commonFunction';
import { userApi } from '/@/api/system/user';

const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '',
		submitTxt: '',
	},
	ruleForm: {},
});
const { copyText } = commonFunction();
const column = ref(2);

const userInfo = ref({
	id: 0,
	uid: '',
	status: '',
	signin: false,
	father: 0,
	invite: 0,
	tokenExpiration: 0,
	roleID: 0,
	roleName: '',
	roleSign: '',
	isAdmin: false,
	createAt: 0,
	phone: '',
	email: '',
	nickname: '',
	header: '',
	wxOpenID: '',
	wxUnionID: '',
	fuwuhaoOpenID: '',
	wxNickName: '',
	wxHeader: '',
	wxWebHook: '',
	hasPassword: false,
	userName: '',
	has2FA: false,
	addressName: '',
	addressPhone: '',
	addressEmail: '',
	addressProvince: '',
	addressCity: '',
	addressDistrict: '',
	addressAddress: '',
	addressArea: [],
});

const userColorEnumTypes: userTypes = userColorTypes;
const UserStatusEnumTypes: userTypes = UserStatusEnum;
const getUserInfo = (userID: number) => {
	userApi()
		.getUserDataData({ id: userID })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				userInfo.value = res.data;
			}
		});
};
const copyUserId = (val: string) => {
	copyText(val);
};
const openDialog = async (row?: any) => {
	if (!row) return;
	state.ruleForm = JSON.parse(JSON.stringify(row));
	state.dialog.title = '父级信息查看';
	state.dialog.submitTxt = '修 改';
	getUserInfo(row.id);
	state.dialog.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
	state.dialog.isShowDialog = false;
};
onMounted(() => {
	const clientWidth = document.body.clientWidth;
	if (clientWidth < 1000) {
		column.value = 1;
	}
});
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style lang="scss" scoped>
.cell-item {
	display: flex;
	align-items: center;
	line-height: 24px;
}

.main-title {
	font-size: 16px;
	font-weight: 600;
	color: #000000;
	line-height: 24px;
}

.main {
	.user-info {
		display: flex;
		justify-content: space-between;
		align-items: center;

		.user-info-right {
			display: flex;
			flex-direction: row;
		}

		.user-head {
			padding: 16px 0px 20px 0px;
			display: flex;
			border-bottom: 1px solid #f1f2f5;

			.user-head-right {
				margin-left: 16px;

				:nth-child(1) {
					height: 22px;
					font-size: 14px;
					font-weight: 400;
					color: #61687c;
					line-height: 22px;
				}

				:nth-child(2) {
					height: 22px;
					font-size: 14px;
					font-weight: 400;
					color: #414960;
					line-height: 22px;
					margin-top: 8px;
				}

				.roleName {
					color: var(--el-text-color-secondary);
				}
			}
		}

		.item-action {
			color: var(--el-color-primary);
			line-height: 40px;
			padding-top: 10px;
			padding-left: 8px;
			cursor: pointer;
		}
	}
}
</style>
