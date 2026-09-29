<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="869px" destroy-on-close @close="closeDialog">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="search-header mb15">

        <el-form :inline="true" :model="state.tableData.param">
          <!--        // 来源网站-->
          <el-form-item label="来源网站">
            <el-select v-model="webID" placeholder="请选择网站来源">
              <el-option label="全部" :value="0"></el-option>
              <el-option v-for="item in websiteList" :key="item.id" :label="item.name" :value="item.id"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="danger" @click="clearAllAuth">
              <el-icon>
                <DeleteFilled/>
              </el-icon>
              清除所有用户授权
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
        <el-table-column prop="webName" label="网站名称" show-overflow-tooltip align="left"></el-table-column>
        <el-table-column prop="isLogin" label="是否在线" show-overflow-tooltip align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.isLogin" type="success">在线</el-tag>
            <el-tag v-else type="warning">下线</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="ip地址" show-overflow-tooltip align="left"></el-table-column>
        <el-table-column prop="geo" label="地理位置" show-overflow-tooltip align="center"></el-table-column>
        <!--        取消授权-->
        <el-table-column label="操作" width="220" fixed="right" align="center">
          <template #default="scope">
            <el-button text type="primary" @click="onTabelRowDel(scope.row)">取消授权</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </el-dialog>
</template>

<script setup lang="ts">
import {reactive} from 'vue';
import {useoauth2Api} from "/@/api/oauth2";
import {ElMessage, ElMessageBox} from "element-plus";
import {ref} from "vue-demi";
import {userApi} from "/@/api/system/user";
import {websiteListType} from "/@/views/oauth2/authorizationList/types";

const state = reactive({
  tableData: {
    data: [],
    total: 0,
    loading: false,
    param: {
      page: 1,
      pagesize: 10,
      uid: 0,
    },
  },
});
const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: '',
  title: '用户授权列表',
  submitTxt: '',
});
const websiteList = ref<websiteListType>([]);
const webID = ref<number>(0);
// 获取网站列表
const getWebsiteList = () => {
  useoauth2Api()
      .oauth2DomainList({page: 1, pagesize: 10000})
      .then((res: any) => {
        if (res.code === "SUCCESS") {
          websiteList.value = res.data.website;
        }
      });
};
//搜索
// 初始化表格数据
const getTableData = () => {
  state.tableData.loading = true;
  useoauth2Api()
      .oauth2List({...state.tableData.param, webID: webID.value})
      .then((res: any) => {
        if (res.code === "SUCCESS") {
          state.tableData.data = res.data.record;
          state.tableData.total = res.data.count;
          state.tableData.loading = false;
        }
      });
};
const closeDialog = () => {
  dialog.isShowDialog = false;
};

const openDialog = async (row?: any) => {
  state.tableData.param.uid = row.id
  dialog.isShowDialog = true;
  getWebsiteList();
  getTableData();
};
// 取消所有授权
const clearAllAuth = () => {
  const webName = webID.value == 0 ? '全部' : websiteList.value.find((item) => item.id == webID.value)?.name;
  ;
  ElMessageBox.confirm(`确定要取消${webName}的所有授权吗，此操作将导致所有授权清空？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        userApi()
            .deleteUserOauth2All({webID: webID.value, uid: state.tableData.param.uid})
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
// 取消授权
const onTabelRowDel = (row: any) => {
  ElMessageBox.confirm('确定要取消授权吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        userApi()
            .deleteUserOauth2({token: row.deleteToken})
            .then((res: any) => {
              if (res.code === "SUCCESS") {
                ElMessage.success('取消成功');
                getTableData();
              }
            });
      })
      .catch(() => {
      });
};

// 暴露变量
defineExpose({
  openDialog,
});
</script>

<style scoped lang="scss"></style>
