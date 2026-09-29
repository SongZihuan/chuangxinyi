<template>
  <el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="869px" @close="closeDialog"
             destroy-on-close>
    <el-descriptions :column="column" border>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr2">
              <ele-User/>
            </el-icon>
            权限ID
          </div>
        </template>
        {{ state.domainInfo.id }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr2">
              <ele-User/>
            </el-icon>
            权限名称
          </div>
        </template>
        {{ state.domainInfo.name }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr2">
              <ele-User/>
            </el-icon>
            权限标识
          </div>
        </template>
        {{ state.domainInfo.sign }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr2">
              <ele-User/>
            </el-icon>
            权限描述
          </div>
        </template>
        {{ state.domainInfo.describe }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr2">
              <ele-User/>
            </el-icon>
            权限排序
          </div>
        </template>
        {{ state.domainInfo.sort }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            <el-icon class="mr2">
              <ele-User/>
            </el-icon>
            权限状态
          </div>
        </template>
        <el-tag type="success" v-if="state.domainInfo.status === 1">启用</el-tag>
        <el-tag type="info" v-else>禁用</el-tag>
      </el-descriptions-item>
    </el-descriptions>
  </el-dialog>
</template>
<script setup lang="ts">
import {reactive} from "vue";

const state = reactive({
  dialog: {
    isShowDialog: false,
    title: '站点信息查看',
    submitTxt: '',
  },
  domainInfo: {} as {
    id: number;
    name: string;
    sign: string;
    describe: string;
    sort: number;
    status: number;
  },
});
const column = 1;
const openDialog = (row?: { webDomain: string; webID: number; webName: string }) => {
  if (!row) return;
  state.domainInfo = JSON.parse(JSON.stringify(row));
  state.dialog.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
  state.dialog.isShowDialog = false;
};
defineExpose({
  openDialog,
  closeDialog,
});
</script>
<style scoped lang="scss">

</style>