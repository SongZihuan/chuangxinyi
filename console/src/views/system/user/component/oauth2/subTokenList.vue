<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="869px" destroy-on-close @close="closeDialog">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="search-header mb15">

        <el-form :inline="true" :model="state.tableData.param">
          <!--        // 来源网站-->
          <el-form-item>
            <el-button type="danger" @click="clearAllToken">
              <el-icon>
                <DeleteFilled/>
              </el-icon>
              清除所有在线子用户
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
        <el-table-column prop="user.id" label="子用户ID" show-overflow-tooltip width="150" align="center"></el-table-column>
        <el-table-column prop="user.phone" label="子用户" show-overflow-tooltip width="150" align="center"></el-table-column>
        <el-table-column prop="ip" label="ip地址" show-overflow-tooltip width="150" align="center"></el-table-column>
        <el-table-column prop="geo" label="地理位置" show-overflow-tooltip width="200" align="center"></el-table-column>
        <el-table-column prop="nowIP" label="当前ip地址" show-overflow-tooltip width="150"
                         align="center"></el-table-column>
        <el-table-column prop="nowGeo" label="当前地理位置" show-overflow-tooltip width="200"
                         align="center"></el-table-column>
        <el-table-column prop="subType" label="Token类型" show-overflow-tooltip width="150" align="center">
          <template #default="scope">
            <span v-if="scope.row.subType === 'UserNotToken'"> 无 </span>
            <span v-else-if="scope.row.subType === 'UserRootToken'"> 根用户 </span>
            <span v-else-if="scope.row.subType === 'UserSonToken'"> 子用户 </span>
            <span v-else-if="scope.row.subType === 'UserFatherToken'"> 父用户 </span>
            <span v-else-if="scope.row.subType === 'UserRootFatherToken'"> 根父用户 </span>
            <span v-else-if="scope.row.subType === 'UserUncleToken'"> 协作人 </span>
            <span v-else-if="scope.row.subType === 'UserHighAuthorityRootToken'"> 子用户（高权限） </span>
            <span v-else-if="scope.row.subType === 'UserWebsiteToken'"> 外站授权 </span>
          </template>
        </el-table-column>
        <el-table-column prop="isLogin" label="是否登录" show-overflow-tooltip width="150"
                         align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.isLogin" type="success">是</el-tag>
            <el-tag v-else type="danger">否</el-tag>
          </template>
        </el-table-column>
        <!--        取消授权-->
        <!--        <el-table-column label="操作" width="220" fixed="right" align="center">-->
        <!--          <template #default="scope">-->
        <!--            <el-button text type="primary" @click="onTabelRowDel(scope.row)">取消授权</el-button>-->
        <!--          </template>-->
        <!--        </el-table-column>-->
      </el-table>
    </el-card>
  </el-dialog>
  <father-info ref="fatherInfoRef" />
</template>

<script setup lang="ts" name="subTokenList">
import {reactive, ref} from 'vue';
import {ElMessage, ElMessageBox} from "element-plus";
import {userApi} from "/@/api/system/user";
import FatherInfo from "/@/views/onlineUser/components/fatherInfo.vue";

const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: '',
  title: '子用户在线列表',
  submitTxt: '',
});
const fatherInfoRef = ref();
const state = reactive({
  tableData: {
    data: [],
    total: 0,
    loading: false,
    uid: 0,
    param: {
      page: 1,
      pagesize: 10,
    },
  },
});
const checkSubUserInfo = (row: any) => {
  fatherInfoRef.value.openDialog(row.user,"子级信息查看");
};
const clearAllToken = () => {
  ElMessageBox.confirm('确定要将该用户所有子用户都下线吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        userApi()
            .deleteUserOnlineSubUserBySite({ id: state.tableData.uid})
            .then((res: any) => {
              if (res.code === "SUCCESS") {
                ElMessage.success('清除成功');
                getTableData();
              }
            });
      })
      .catch(() => {
      });
};
//搜索
// 初始化表格数据
const getTableData = () => {
  state.tableData.loading = true;
  userApi()
      .getUserOnlineSubUser({uid: state.tableData.uid})
      .then((res: any) => {
        if (res.code === "SUCCESS") {
          state.tableData.data = res.data.token;
          state.tableData.total = res.data.count;
          state.tableData.loading = false;
        }
      });
};
const closeDialog = () => {
  dialog.isShowDialog = false;
};

const openDialog = async (row?: any) => {
  state.tableData.uid = row.id
  dialog.isShowDialog = true;
  getTableData();
};

// 暴露变量
defineExpose({
  openDialog,
});
</script>

<style scoped lang="scss"></style>
