<template>
	<div class="mt20" v-if="hasInvite">
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
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import dayjs from 'dayjs';
import { useInviteApi } from '/@/api/invite';
import { userColorTypes, UserStatusEnum, userTypes } from '/@/data/enum';
import commonFunction from '/@/utils/commonFunction';
const useInviteApiCollect = useInviteApi();
const { copyText } = commonFunction();
const userColorEnumTypes: userTypes = userColorTypes;
const UserStatusEnumTypes: userTypes = UserStatusEnum;

const column = ref(2);
const hasInvite = ref();
const userInfo = ref({});
const getUserInfo = () => {
	useInviteApiCollect.inviterInfo().then((res: any) => {
		if (res.code === "SUCCESS") {
			hasInvite.value = res.data.hasInvite;
			userInfo.value = res.data.invite;
		}
	});
};
const copyUserId = (val: string) => {
	copyText(val);
};
onMounted(() => {
	const clientWidth = document.body.clientWidth;
	if (clientWidth < 1000) {
		column.value = 1;
	}
	getUserInfo();
});
</script>

<style lang="scss" scoped></style>
