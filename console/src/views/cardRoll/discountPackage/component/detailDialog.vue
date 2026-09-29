<template>
	<div>
		<!--查看详情-->
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="600px" destroy-on-close @close="closeDialog">
			<el-form>
				<el-row :gutter="35">
          <el-col :span="24">
            <el-form-item label="名称">
              <span>{{ info.name }}</span>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="描述">
              <span>{{ info.shortDescribe }}</span>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" v-if="info.dayLimit">
            <el-form-item label="日限额">
              <span>{{ info.dayLimit }}个</span>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" v-if="info.monthLimit">
            <el-form-item label="月限额">
              <span>{{ info.monthLimit }}个</span>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" v-if="info.yearLimit">
            <el-form-item label="年限额">
              <span>{{ info.yearLimit }}个</span>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" v-if="info.limit">
            <el-form-item label="总限额">
              <span>{{ info.limit }}个</span>
            </el-form-item>
          </el-col>
					<el-col :span="24">
						<el-form-item label="优惠包形式">
							<el-tag type="success" v-if="info.type === 1">赠送可消费余额</el-tag>
							<el-tag type="success" v-else-if="info.type === 2">赠送优惠券</el-tag>
						</el-form-item>
					</el-col>
          <template v-if="info.type === 1">
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
              <el-form-item label="金额">
                <span>{{ formatAmount(info.quota.amount) }}元</span>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" v-if="info.quota.canWithdraw">
              <el-form-item label="是否可提现">
                <span>可以</span>
              </el-form-item>
            </el-col>
          </template>
          <template v-else-if="info.type === 2">
            <template v-if="info.quota.type === 1">
              <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                <el-form-item label="充值起送线">
                  <span>{{ formatAmount(info.quota.bottom) }}元</span>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                <el-form-item label="赠送金额">
                  <span>{{ formatAmount(info.quota.send || 0) }}元</span>
                </el-form-item>
              </el-col>
            </template>
            <template v-else-if="info.quota.type === 2">
              <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                <el-form-item label="消费起抵线">
                  <span>{{ formatAmount(info.quota.bottom) }}元</span>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                <el-form-item label="抵扣金额">
                  <span>{{ formatAmount(info.quota.discount || 0) }}元</span>
                </el-form-item>
              </el-col>
            </template>
            <template v-else-if="info.quota.type === 3">
              <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                <el-form-item label="消费起抵线">
                  <span>{{ formatAmount(info.quota.bottom) }}元</span>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
                <el-form-item label="打折">
                  <span>{{ info.quota.pre }}%</span>
                </el-form-item>
              </el-col>
            </template>
          </template>
          <el-col>
            <el-form-item label="个人购买条件">
              <el-tag type="warning" v-if="info.needCompany || info.needCompanyFaceCheck || info.needCompanyOrigin">个人不能购买</el-tag>
              <el-tag type="success" v-else-if="info.needUserFace">需要人脸核验</el-tag>
              <el-tag type="success" v-else-if="info.needUserOrigin">需要身份信息上传核验</el-tag>
              <el-tag type="success" v-else-if="info.needVerify">需要实名核验</el-tag>
              <el-tag type="success" v-else>无限制</el-tag>
            </el-form-item>
          </el-col>
          <el-col>
            <el-form-item label="企业购买条件（使用人）">
              <el-tag type="success" v-if="info.needUserFace">需要人脸核验</el-tag>
              <el-tag type="success" v-else-if="info.needUserOrigin">需要身份信息上传核验</el-tag>
              <el-tag type="success" v-else-if="info.needVerify">需要实名核验</el-tag>
              <el-tag type="success" v-else>无限制</el-tag>
            </el-form-item>
          </el-col>
          <el-col>
            <el-form-item label="企业购买条件（企业）">
              <el-tag type="success" v-if="info.needCompanyFaceCheck">需要法人人脸核验</el-tag>
              <el-tag type="success" v-else-if="info.needCompanyOrigin">需要企业证件上传核验</el-tag>
              <el-tag type="success" v-else-if="info.needCompany">无限制</el-tag>
              <el-tag type="success" v-else>企业和个人均可</el-tag>
            </el-form-item>
          </el-col>
				</el-row>
			</el-form>
		</el-dialog>
	</div>
</template>
<script setup lang="ts" name="invoiceInfo">
import { reactive, ref } from 'vue';
import type { formTypes } from '/@/api/cardRoll/discountPackage/types';
import { formatAmount } from '/@/utils/formatAmount';

const info = ref<formTypes>({
	name: '',
	describe: '',
	shortDescribe: '',
	type: 1,
	quota: {
		type: 1,
		bottom: 1,
		send: 1,
		discount: 1,
		pre: 9,
	},
	dayLimit: 1,
	monthLimit: 5,
	yearLimit: 15,
	limit: 30,
	needVerify: true,
	needCompany: true,
	needUserOrigin: true,
	needCompanyOrigin: true,
	needUserFace: true,
	needCompanyFace: true,
	show: true,
});
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'look',
	title: '优惠包详情',
	submitTxt: '查看',
});
const openDialog = (row: any) => {
	info.value = row;
	dialog.isShowDialog = true;
};
//重置
const closeDialog = () => {
	dialog.isShowDialog = false;
};
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style scoped lang="scss"></style>
