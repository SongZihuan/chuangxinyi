<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="869px" @close="closeDialog" destroy-on-close>
		<el-card shadow="hover" class="mb20">
			<el-descriptions title="用户信息" direction="vertical" :column="column">
				<el-descriptions-item label="用户ID">{{ homePageData.id || '-' }}</el-descriptions-item>
				<el-descriptions-item label="用户UID">{{ homePageData.uid || '-' }}</el-descriptions-item>
				<el-descriptions-item label="用户状态">
					<el-tag v-if="homePageData.status === 'REGISTER'" type="success">注册</el-tag>
					<el-tag v-else-if="homePageData.status === 'NORMAL'" type="success">正常</el-tag>
					<el-tag v-else-if="homePageData.status === 'BANNED'" type="danger">禁用</el-tag>
					<el-tag v-else-if="homePageData.status === 'DELETE'" type="danger">注销</el-tag>
					<el-tag v-else-if="homePageData.status === 'FREEZE'" type="danger">冻结</el-tag>
				</el-descriptions-item>
				<el-descriptions-item label="是否管理员">
					<el-tag v-if="homePageData.isAdmin" type="success">是</el-tag>
					<el-tag v-else type="danger">否</el-tag>
				</el-descriptions-item>
				<el-descriptions-item label="手机号">{{ homePageData.phone || '-' }}</el-descriptions-item>
				<el-descriptions-item label="邮箱">{{ homePageData.email || '-' }}</el-descriptions-item>
				<el-descriptions-item label="昵称">{{ homePageData.nickname || '-' }}</el-descriptions-item>
				<el-descriptions-item label="微信OpenID">{{ homePageData.wxOpenID || '-' }}</el-descriptions-item>
				<el-descriptions-item label="微信UnionID">{{ homePageData.wxUnionID || '-' }}</el-descriptions-item>
				<el-descriptions-item label="服务号OpenID">{{ homePageData.fuwuhaoOpenID || '-' }}</el-descriptions-item>
				<el-descriptions-item label="微信昵称">{{ homePageData.wxNickName || '-' }}</el-descriptions-item>
				<el-descriptions-item label="微信头像">
					<el-image
						v-if="homePageData.wxHeader"
						:src="getFile(homePageData.wxHeader)"
						fit="cover"
						:preview-src-list="[getFile(homePageData.wxHeader)]"
					/>
					<div v-else>-</div>
				</el-descriptions-item>
				<el-descriptions-item label="微信WebHook">{{ homePageData.wxWebHook || '-' }}</el-descriptions-item>
				<el-descriptions-item label="是否有密码">
					<el-tag v-if="homePageData.hasPassword" type="success">是</el-tag>
					<el-tag v-else type="danger">否</el-tag>
				</el-descriptions-item>
				<el-descriptions-item label="用户名">{{ homePageData.userName || '-' }}</el-descriptions-item>
				<el-descriptions-item label="是否有2FA">
					<el-tag v-if="homePageData.has2FA" type="success">是</el-tag>
					<el-tag v-else type="danger">否</el-tag>
				</el-descriptions-item>
				<el-descriptions-item label="创建时间">
					{{ dayjs.unix(homePageData.createAt).format('YYYY-MM-DD HH:mm:ss') }}
				</el-descriptions-item>
			</el-descriptions>
		</el-card>
		<el-card shadow="hover" class="mb20">
			<el-descriptions title="用户角色" direction="vertical" :column="column">
				<el-descriptions-item label="角色ID">{{ homePageData.roleID || '-' }}</el-descriptions-item>
				<el-descriptions-item label="角色名称">{{ homePageData.roleName || '-' }}</el-descriptions-item>
				<el-descriptions-item label="角色标识">{{ homePageData.roleSign || '-' }}</el-descriptions-item>
			</el-descriptions>
		</el-card>
		<el-card shadow="hover">
			<el-descriptions title="地址信息" direction="vertical" :column="column">
				<el-descriptions-item label="姓名">{{ homePageData.addressName || '-' }}</el-descriptions-item>
				<el-descriptions-item label="手机号">{{ homePageData.addressPhone || '-' }}</el-descriptions-item>
				<el-descriptions-item label="邮箱">{{ homePageData.addressEmail || '-' }}</el-descriptions-item>
				<el-descriptions-item label="省份">{{ homePageData.addressProvince || '-' }}</el-descriptions-item>
				<el-descriptions-item label="城市">{{ homePageData.addressCity || '-' }}</el-descriptions-item>
				<el-descriptions-item label="区县">{{ homePageData.addressDistrict || '-' }}</el-descriptions-item>
				<el-descriptions-item label="详细地址">{{ homePageData.addressAddress || '-' }}</el-descriptions-item>
				<el-descriptions-item label="地址区域">
					<div class="cell-item" v-for="(item, index) in homePageData.addressArea" :key="index">
						<el-icon class="mr5">
							<ele-MapLocation />
						</el-icon>
						{{ item || '-' }}
					</div>
				</el-descriptions-item>
			</el-descriptions>
		</el-card>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { pageHomeParamsTypes, userDataTypes } from '/@/api/homePage/types';
import { ElMessage } from 'element-plus';
import { userApi } from '/@/api/system/user';
import dayjs from 'dayjs';

const column = ref(1);
const homePageData = ref<userDataTypes>({
	id: 0,
	uid: '',
	status: '',
	signin: false,
	father: 0,
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
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '用户信息',
	submitTxt: '新增',
});
// 获取homepage数据
const getHomePageData = async (userID: string) => {
	if (!userID) {
		ElMessage.warning('抱歉，您没有登录权限');
		return;
	}
	let data: pageHomeParamsTypes = {
		uid: userID,
	};
	await userApi()
		.getUserDataData(data)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				homePageData.value = res.data;
				return;
			}
		});
};
const getFile = (file: string) => {
	return `${file}`;
};
const openDialog = (row: any) => {
	if (!row.id) {
		ElMessage.warning('用户获取失败');
		return;
	}
	getHomePageData(row.id as string);
	dialog.isShowDialog = true;
};
const closeDialog = () => {
	dialog.isShowDialog = false;
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
