<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" height="500px" center>
    <template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">发送</el-button>
			</span>
    </template>
  </el-dialog>
</template>
<script setup lang="ts">
import {reactive, ref, toRefs} from "vue";

const state = reactive({
  dialog: {
    isShowDialog: false,
    type: '',
    title: '发送邮件',
    submitTxt: '',
  },
});
let ruleForm = ref({
  uid: null,
});
const closeDialog = () => {
  dialog.value.isShowDialog = false;
};
const { dialog } = toRefs(state);
const openDialog = (row: any) => {
  ruleForm.value.uid = row.id;
  dialog.value.isShowDialog = true;
};

defineExpose({
  openDialog,
});
</script>
<style scoped lang="scss">

</style>