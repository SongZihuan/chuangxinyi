<template>
	<el-card shadow="never">
    <el-row>
      <el-col :span="6">
        <el-statistic title="用户账号余额（元）" :value="formatAmount(props.financeInfo.balance)" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="待开发票金额（元）" :value="formatAmount(props.financeInfo.notBilled)" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="待提现收益（元）" :value="formatAmount(props.financeInfo.notWithdraw)" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="待入账收益（元）" :value="formatAmount(props.financeInfo.waitWithdraw + props.financeInfo.waitBalance)" />
      </el-col>
    </el-row>
	</el-card>
</template>
<script setup lang="ts">
import { formatAmount } from '/@/utils/formatAmount';
import { ref } from 'vue';
import { userWalletTypes } from '/@/views/system/user/types';

const column = ref(2);

interface Props {
	financeInfo: userWalletTypes;
}

const props = withDefaults(defineProps<Props>(), {});
</script>
<style scoped lang="scss">
.el-col {
  text-align: center;
}
</style>
