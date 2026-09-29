<template>
	<div class="mt20">
		<el-card shadow="hover">
			<el-descriptions title="地址信息" direction="vertical" :column="column">
				<template #extra>
					<el-button type="primary" text @click="onEdit">修改</el-button>
				</template>
				<el-descriptions-item label="姓名"
					><span>{{ userInfo.name }}</span></el-descriptions-item
				>
				<el-descriptions-item label="手机">{{ userInfo.phone || '-' }}</el-descriptions-item>
				<el-descriptions-item label="邮箱">{{ userInfo.email || '-' }}</el-descriptions-item>
				<el-descriptions-item label="省份">{{ userInfo.province || '-' }}</el-descriptions-item>
				<el-descriptions-item label="市区">{{ userInfo.city || '-' }}</el-descriptions-item>
				<el-descriptions-item label="区">{{ userInfo.district || '-' }}</el-descriptions-item>
				<el-descriptions-item label="地址">
					<template #label> 地址 <el-tag size="small"> 收货地址 </el-tag> </template>
					{{ userInfo.address || '-' }}
				</el-descriptions-item>
			</el-descriptions>
		</el-card>
		<!-- 编辑个人信息 -->
		<Edit ref="editRef" :userInfo="userInfo" @refresh="getUserInfo" />
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Edit from '/@/views/userCenter/userInfo/component/edit.vue';
import { useUserApi } from '/@/api/user/user';
import { Session } from '/@/utils/storage';
const useUserApiCollect = useUserApi();
const column = ref(3);
const editRef = ref();
const userInfo = ref({
	name: '',
	phone: '',
	email: '',
	province: '',
	city: '',
	district: '',
	address: '',
	areas: [],
});
const getUserInfo = () => {
	useUserApiCollect.userInfo().then((res: any) => {
		if (res.code === "SUCCESS") {
			userInfo.value = res.data.address;
			Session.set('userData', res.data);
		}
	});
};
//编辑个人信息
const onEdit = () => {
	editRef.value.openDialog();
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
