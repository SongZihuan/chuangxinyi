<template>
	<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="869px" @close="closeDialog" destroy-on-close>
		<div class="mt20">
			<el-card shadow="hover">
				<el-descriptions title="邀请人名片" :column="column" border>
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
					<el-descriptions-item label="用户名"
						><template #label>
							<div class="cell-item">
								<el-icon class="mr4">
									<ele-User />
								</el-icon>
								用户名称
							</div>
						</template>
						{{ userInfo.userName }}</el-descriptions-item
					>
					<el-descriptions-item label="昵称">
						<template #label>
							<div class="cell-item">
								<el-icon class="mr4">
									<ele-User />
								</el-icon>
								用户昵称
							</div>
						</template>
						<span>{{ userInfo.nickname }}</span></el-descriptions-item
					>
					<el-descriptions-item label="手机"
						><template #label>
							<div class="cell-item">
								<el-icon class="mr4"> <ele-Iphone /> </el-icon>
								手机号
							</div> </template
						>{{ userInfo.phone || '-' }}</el-descriptions-item
					>
					<el-descriptions-item label="邮箱"
						><template #label>
							<div class="cell-item">
								<el-icon class="el-input__icon"><ele-Message /></el-icon>
								邮箱
							</div>
						</template>
						{{ userInfo.email }}</el-descriptions-item
					>
					<el-descriptions-item label="注册时间">
						<template #label>
							<div class="cell-item">
								<SvgIcon name="my-peoples" :size="14" color="#606266"></SvgIcon>
								注册时间
							</div> </template
						>{{ dayjs.unix(userInfo.createAt).format('YYYY-MM-DD HH:mm:ss') || '-' }}</el-descriptions-item
					>
					<el-descriptions-item>
						<template #label>
							<div class="cell-item">
								<SvgIcon name="ele-CollectionTag" color="#606266"></SvgIcon>
								<span class="mt4"> 状态</span>
							</div>
						</template>
						<el-text class="mx-1" :type="userColorEnumTypes[userInfo.status]" size="mini">{{ UserStatusEnumTypes[userInfo.status] }}</el-text>
					</el-descriptions-item>
				</el-descriptions>
			</el-card>
		</div>
	</el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import dayjs from 'dayjs';
import { userColorTypes, UserStatusEnum, userTypes } from '/@/data/enum';
import commonFunction from '/@/utils/commonFunction';
import { userApi } from '/@/api/system/user';
const { copyText } = commonFunction();
const userColorEnumTypes: userTypes = userColorTypes;
const UserStatusEnumTypes: userTypes = UserStatusEnum;

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
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '',
		submitTxt: '',
	},
	ruleForm: {},
});
const getUserInfo = (userID: number) => {
	userApi()
		.getUserDataData({ uid: userID })
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

<style lang="scss" scoped></style>
