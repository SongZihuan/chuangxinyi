<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="450px" destroy-on-close @close="closeDialog">
    <el-descriptions :column="column" border>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            用户ID
          </div>
        </template>
        {{ ruleForm.id }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            手机号
          </div>
        </template>
        {{ ruleForm.phone }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">
            状态
          </div>
        </template>
        <el-tag v-if="ruleForm.status === 'NORMAL'" type="success">正常</el-tag>
        <el-tag v-else-if="ruleForm.status === 'DISABLE'" type="danger">禁用</el-tag>
        <el-tag v-else-if="ruleForm.status === 'DELETE'" type="danger">删除</el-tag>
      </el-descriptions-item>
    </el-descriptions>
  </el-dialog>
</template>

<script setup lang="ts" name="modePaymentDialog">
import { reactive, ref, nextTick } from 'vue';
import dayjs from "dayjs";
const emit = defineEmits(['queryOrderBack']);
const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: '',
  title: '购买人信息',
  submitTxt: '确认',
});
const ruleForm = ref<any>({
  phone: '',
  status: '',
  createAt: '',
  inviteCount: '',
});

const column = 1;

const openDialog = (row:any) => {
  if (!row) return;
  ruleForm.value = JSON.parse(JSON.stringify(row));
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

<style lang="scss" scoped>
</style>
