<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="850px" destroy-on-close @close="closeDialog">
    <el-row v-if="type !== 'admin'">
      <el-col :span="24">
        <div class="mb20">
          <el-button text type="primary" size="mini" @click="changAdminUser" v-if="nowToken !== loginToken.token">点击切换回登录账号</el-button>
        </div>
      </el-col>
    </el-row>
		<admin-sub-user-table v-if="type === 'admin'" :state="state" />
		<user-sub-user-table v-else :state="state" @changUser="changUser" />
	</el-dialog>
</template>

<script setup lang="ts">
import {computed, reactive, ref} from 'vue';
import { useSubUserApi } from '/@/api/user/subuser';
import { subAccountStatsTypes } from '/@/views/system/user/types';
import { ElMessage } from 'element-plus';
import { useSubUserCenterApi } from '/@/api/subUser';
import { useLoginSignIn, useLoginSignSub } from '/@/hooks/useLoginSignIn';
import { Local, Session } from '/@/utils/storage';
import AdminSubUserTable from '/@/views/system/user/component/adminSubUserTable.vue';
import UserSubUserTable from '/@/views/system/user/component/userSubUserTable.vue';
import { subUserAccountStatsTypes } from '/@/views/accountManagement/subUser/types';
import { NextLoading } from '/@/utils/loading';

const props = defineProps({
	type: {
		type: String,
		default: 'admin',
	},
});
const info = ref<any>({});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '子账户信息',
	submitTxt: '',
});

const state = reactive<subAccountStatsTypes | subUserAccountStatsTypes>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			id: '',
			uid: '',
		},
	},
});
const getAdminTableData = async () => {
	state.tableData.loading = true;
	await useSubUserApi()
		.subUserInfo(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.user;
				state.tableData.loading = false;
			}
		});
};
const getUserTableData = async () => {
	state.tableData.loading = true;
	await useSubUserCenterApi()
		.getSubUserListByLoginToken()
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.user;
				state.tableData.loading = false;
			}
		});
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const { onSignSub } = useLoginSignSub();
const { resetSignIn } = useLoginSignIn();

const changUser = async (row: any) => {
  closeDialog();
  let res = await onSignSub(row.id);
  if (res) {
    ElMessage.success('切换成功');
    nowToken = Session.get('token')
  } else {
    ElMessage.success('切换失败');

    let resetRes = await resetSignIn(false);
    if (!resetRes) {
      Session.clear()
      Local.clear();
      window.location.href = "/login"
    }
  }
  NextLoading.done()
};
const changAdminUser = async () => {
  closeDialog();
  let resetRes = await resetSignIn(false);
  if (!resetRes) {
    Session.clear()
    Local.clear();
    window.location.reload()
  } else {
    ElMessage.success('切换成功');
    nowToken = Session.get('token')
  }
  NextLoading.done()
};

const openDialog = async (row?: any) => {
	state.tableData.param.uid = row.id;
	if (props.type == 'admin') {
		await getAdminTableData();
		info.value = JSON.parse(JSON.stringify(row));
	} else {
		await getUserTableData();
		info.value = JSON.parse(JSON.stringify(row));
		dialog.title = '账号切换';
	}
	dialog.isShowDialog = true;
};

let nowToken = Session.get('token')
let loginToken = Session.get('login-token')

// 暴露变量
defineExpose({
	openDialog,
});
</script>
