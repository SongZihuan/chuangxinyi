<template>
  <div>
    <el-table
        :data="tableData.data"
        v-loading="tableData.loading"
        style="width: 100%"
        row-key="userID"
        default-expand-all
        :tree-props="{ children: 'son', hasChildren: 'hasChildren' }"
    >
      <el-table-column prop="userID" label="用户数字ID" show-overflow-tooltip align="center"></el-table-column>
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
          <el-tag type="success" v-if="scope.row.status == 'NORMAL'">正常</el-tag>
          <el-tag type="success" v-else-if="scope.row.status == 'REGISTER'">注册</el-tag>
          <el-tag type="danger" v-else-if="scope.row.status == 'BANNED'">封禁</el-tag>
          <el-tag type="warning" v-else-if="scope.row.status == 'DELETE'">注销</el-tag>
          <el-tag type="danger" v-else-if="scope.row.status == 'FREEZE'">冻结</el-tag>
        </template>
      </el-table-column>
      <!--    查看子账户token-->
      <el-table-column label="操作" fixed="right" align="center">
        <template #default="scope">
          <el-button text type="primary" @click="oncheckToken(scope.row)" v-if="scope.row.status == 'NORMAL'">在线列表</el-button>
        </template>
      </el-table-column>

    </el-table>
    <sub-token-list ref="subTokenListRef"  />
  </div>
</template>
<script setup lang="ts">
import {reactive, ref} from 'vue';
import {subAccountStatsTypes} from "/@/views/system/user/types";
import SubTokenList from "/@/views/system/user/component/oauth2/subTokenList.vue";
const props = defineProps({
  state: {
    type: Object,
    required: true,
  }
});
const subTokenListRef = ref();
const oncheckToken = (row: any) => {
  subTokenListRef.value.openDialog(row);
};
const state = props.state as subAccountStatsTypes;
const tableData = reactive(state.tableData) ;
</script>
<style scoped lang="scss"></style>
