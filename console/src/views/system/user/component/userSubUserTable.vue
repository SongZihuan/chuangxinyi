<template>
	<el-table
		:data="tableData.data"
		v-loading="tableData.loading"
		style="width: 100%"
		row-key="id"
		default-expand-all
		:tree-props="{ children: 'son', hasChildren: 'hasChildren' }"
	>
		<el-table-column prop="id" label="用户ID" show-overflow-tooltip align="center"></el-table-column>
		<el-table-column prop="phone" label="手机号" show-overflow-tooltip align="center"></el-table-column>
		<el-table-column prop="roleName" label="角色" show-overflow-tooltip align="center"></el-table-column>
		<el-table-column prop="lineal" label="关系" show-overflow-tooltip align="center">
			<template #default="scope">
				<el-tag type="success" v-if="scope.row.lineal">子账号</el-tag>
				<el-tag type="success" v-else-if="!scope.row.lineal && scope.row.nephewStatus === 2">协作账号</el-tag>
        <el-tag type="info" v-else>未确认协作账号</el-tag>
			</template>
		</el-table-column>
		<el-table-column prop="status" label="状态" show-overflow-tooltip align="center">
			<template #default="scope">
        <el-tag type="success" v-if="scope.row.status == 'NORMAL' || scope.row.status == 'REGISTER'">正常</el-tag>
        <el-tag type="warning" v-else-if="scope.row.status == 'DELETE'">注销</el-tag>
        <el-tag type="danger" v-else>封禁</el-tag>
			</template>
		</el-table-column>
		<el-table-column label="操作" fixed="right" align="center" v-if="authUser.getSonToken">
			<template #default="scope">
				<el-button text type="primary" size="mini" @click="changUser(scope.row)" v-if="scope.row.status == 'NORMAL'">切换</el-button>
				<el-button type="info" disabled v-else>账号已禁用</el-button>
			</template>
		</el-table-column>
	</el-table>
</template>
<script setup lang="ts">
import { reactive } from 'vue';
import { subUserAccountStatsTypes } from '/@/views/accountManagement/subUser/types';
import dayjs from 'dayjs';
import useSubAuth from '/@/hooks/useSubAuth';

const props = defineProps({
	state: {
		type: Object,
		required: true,
	},
});
const emit = defineEmits(['changUser']);
const state = props.state as subUserAccountStatsTypes;
const authUser = useSubAuth();
const tableData = reactive(state.tableData);
const changUser = (row: any) => {
	emit('changUser', row);
};
</script>
<style scoped lang="scss"></style>
