<template>
  <el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
    <!-- 剩余余额  -->
    <el-form ref="formRef" :model="ruleForm" size="default" label-width="65px" :rules="rules">
      <el-row :gutter="35">
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
          <el-form-item label="用户ID" prop="uid">
            <el-input v-model="ruleForm.uid" placeholder="请输入用户ID" clearable></el-input>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="invoiceDialog">
import {reactive, ref, onMounted} from 'vue';
import {type FormRules, type FormInstance} from 'element-plus';
import {ElMessage} from 'element-plus';
import {useCouponApi} from '/@/api/cardRoll/discountPackage/index';

const useCouponApiCollect = useCouponApi();
const formRef = ref();
let ruleForm = ref({
  uid: null,
  id: null,
  discountID: null,
});
const rules = reactive<FormRules>({
  uid: [
    {
      required: true,
      message: '请输入用户ID',
      trigger: 'blur',
    },
  ],
});

const dialog = reactive<dialogParams>({
  isShowDialog: false,
  type: 'add',
  title: '分配用户优惠包',
  submitTxt: '分配',
});

const openDialog = (type: string, row: any) => {
  ruleForm.value.discountID = row.id;
  dialog.isShowDialog = true;
};
//重置
const closeDialog = () => {
  dialog.isShowDialog = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate((valid) => {
    if (valid) {
      useCouponApiCollect.allocationUser({...ruleForm.value, uid: ruleForm.value.uid}).then((res: any) => {
        if (res.code === "SUCCESS") {
          closeDialog();
          ElMessage({
            type: 'success',
            message: '获取成功',
          });
        }
      });
    } else {
      return false;
    }
  });
};

onMounted(() => {
});
defineExpose({
  openDialog,
  closeDialog,
});
</script>

<style lang="scss" scoped>
.tip {
  display: flex;
  padding: 10px 10px;
}
</style>
