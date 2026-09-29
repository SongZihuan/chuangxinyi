<template>
	<div class="container layout-padding">
		<el-card shadow="hover">
			<div class="my-wallet">
				<div class="header">
					<div class="title">我的钱包</div>
					<div class="icon">
						<el-icon size="40" color="#b5b8ba">
							<ele-CreditCard />
						</el-icon>
					</div>
				</div>
				<div class="balance">
					￥{{ formatAmount(financeInfo.balance + financeInfo.notWithdraw + financeInfo.waitBalance + financeInfo.waitWithdraw) }}
					<span class="withdraw">（仅提现：￥{{ formatAmount(financeInfo.notWithdraw + financeInfo.waitWithdraw) }}）</span>
				</div>
				<div class="hint">
					钱包组成部分
					<el-tooltip effect="dark" placement="right">
						<template #content>
							账户余额（仅消费）：￥{{ formatAmount(financeInfo.balance) }}
							<br />
							优惠返现（仅提现）：￥{{ formatAmount(financeInfo.notWithdraw) }}
							<br />
							待入账优惠返现（仅消费）：￥{{ formatAmount(financeInfo.waitBalance) }}
							<br />
							待入账优惠返现（仅提现）：￥{{ formatAmount(financeInfo.waitWithdraw) }}
							<br />
							已提现金额：￥{{ formatAmount(financeInfo.hasWithdraw) }}
							<br />
							可开票金额：￥{{ formatAmount(financeInfo.notBilled) }}
							<br />
							已开票金额：￥{{ formatAmount(financeInfo.hasBilled) }}
							<br />
							总充值金额：￥{{ formatAmount(financeInfo.cny) }}
						</template>
						<el-icon size="20" style="margin-left: 5px"><ele-QuestionFilled /></el-icon>
					</el-tooltip>
				</div>
				<div class="tool-btns">
					<el-button class="btn" type="primary" size="small" @click="onRechargeChange"
						><el-icon size="20"><ele-Coin /></el-icon>充值</el-button
					>
					<el-button class="btn" type="primary" size="small" @click="onWithdrawChange"
						><el-icon size="20"><ele-Wallet /></el-icon>提现</el-button
					>
				</div>
			</div>
		</el-card>
		<el-card shadow="hover" class="layout-padding-auto mt20">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<el-form-item>
						<el-input v-model="state.tableData.param.src" placeholder="请输入关键词" style="width: 100%"> </el-input>
					</el-form-item>
					<el-form-item>
						<el-input v-model="state.tableData.param.fundingID" placeholder="请输入资金ID" style="width: 100%"> </el-input>
					</el-form-item>
					<el-form-item>
						<el-select v-model="state.tableData.param.timetype" placeholder="请选择日期类型" style="width: 100%">
							<el-option v-for="(item, index) in dateTypeDict" :value="item.value" :label="item.label" :key="index" />
						</el-select>
					</el-form-item>
					<el-form-item>
						<el-date-picker
							v-model="state.tableData.param.range"
							type="datetimerange"
							range-separator="至"
							start-placeholder="开始日期"
							end-placeholder="结束日期"
							style="width: 100%"
						/>
					</el-form-item>
					<el-form-item>
						<el-button type="primary" class="ml10" @click="handleSearch">
							<el-icon>
								<ele-Search />
							</el-icon>
							查询
						</el-button>
					</el-form-item>
				</el-form>
			</div>
			<price />
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column prop="fundingID" label="资金ID" show-overflow-tooltip min-width="150px" align="center"></el-table-column>
				<el-table-column prop="user.id" label="操作人ID" show-overflow-tooltip min-width="150px" align="center"></el-table-column>
				<el-table-column prop="user.phone" label="操作人" show-overflow-tooltip min-width="150px" align="center"></el-table-column>
				<el-table-column prop="type" label="类型" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.type === 1" type="success">充值</el-tag>
						<el-tag v-else-if="scope.row.type === 2" type="success">返现</el-tag>
						<el-tag v-else-if="scope.row.type === 3" type="success">消费</el-tag>
						<el-tag v-else-if="scope.row.type === 4" type="success">发票</el-tag>
						<el-tag v-else-if="scope.row.type === 5" type="warning">系统操作</el-tag>
						<el-tag v-else-if="scope.row.type === 7" type="success">提现</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="reason" label="原因" show-overflow-tooltip min-width="150px" align="center"></el-table-column>
				<el-table-column prop="balance" label="余额" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeBalance) }}元 -> {{ formatAmount(scope.row.balance) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="待入账余额" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeWaitBalance) }}元 -> {{ formatAmount(scope.row.waitBalance) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="充值总金额" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeCny) }}元 -> {{ formatAmount(scope.row.cny) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="未开票额度" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeNotBilled) }}元 -> {{ formatAmount(scope.row.notBilled) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="已开票额度" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeHasBilled) }}元 -> {{ formatAmount(scope.row.hasBilled) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="可开票额度" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeBilled) }}元 -> {{ formatAmount(scope.row.billed) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="未提现额度" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeNotWithdraw) }}元 -> {{ formatAmount(scope.row.notWithdraw) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="已提现额度" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeHasWithdraw) }}元 -> {{ formatAmount(scope.row.hasWithdraw) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="可提现额度" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeWithdraw) }}元 -> {{ formatAmount(scope.row.withdraw) }}元</span>
					</template>
				</el-table-column>
				<el-table-column prop="balance" label="待入账可提现额度" show-overflow-tooltip min-width="150px" align="center">
					<template #default="scope">
						<span>{{ formatAmount(scope.row.beforeWaitWithdraw) }}元 -> {{ formatAmount(scope.row.waitWithdraw) }}元</span>
					</template>
				</el-table-column>
				<el-table-column label="变更时间" show-overflow-tooltip width="170" align="center">
					<template #default="scope">
						<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
						<span v-else>-</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="100" fixed="right" align="center">
					<template #default="scope">
						<el-button text type="primary" @click="onCopy(scope.row.fundingID)">复制资金ID</el-button>
					</template>
				</el-table-column>
			</el-table>
			<el-pagination
				hide-on-single-page
				@size-change="onHandleSizeChange"
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				:page-sizes="[10, 20, 30]"
				v-model:current-page="state.tableData.param.page"
				background
				v-model:page-size="state.tableData.param.pagesize"
				layout="total, sizes, prev, pager, next, jumper"
				:total="state.tableData.total"
			>
			</el-pagination>
		</el-card>
		<rechargeDialog ref="rechargeDialogRef" @refresh="getUserInfo" />
		<withdrawDialog ref="withdrawDialogRef" @refresh="getTableData" />
	</div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { formatAmount } from '/@/utils/formatAmount';
import dayjs from 'dayjs';
import { useWalletRecordApi } from '/@/api/wallet/record';
import commonFunction from '/@/utils/commonFunction';
import { Session } from '/@/utils/storage';
import { useUserInfo } from '/@/stores/userInfo';
import { useUserApi } from '/@/api/user/user';
import { financeInfoType } from '/@/views/wallet/types';
import rechargeDialog from '../recharge/index.vue';
import withdrawDialog from '../withdraw/index.vue';
const state = reactive<any>({
	tableData: {
		data: [],
		loading: false,
		param: {
			range: [],
			page: 1,
			pagesize: 20,
			timetype: '',
			starttime: null,
			endtime: null,
			status: '',
			src: '',
			fundingID: '',
		},
	},
});
const financeInfo = ref<financeInfoType>({
	walletID: 0,
	balance: 0,
	notBilled: 0,
	billed: 0,
	hasBilled: 0,
	cny: 0,
	waitBalance: 0,
	withdraw: 0,
	waitWithdraw: 0,
	notWithdraw: 0,
	hasWithdraw: 0,
});
watch(
	() => state.tableData.param.range,
	() => {
		if (state.tableData.param.range.length > 0) {
			state.tableData.param.starttime = dayjs(state.tableData.param.range[0]).unix();
			state.tableData.param.endtime = dayjs(state.tableData.param.range[1]).unix();
		}
	},
	{
		deep: true,
	}
);
const rechargeDialogRef = ref();
const withdrawDialogRef = ref();
const onRechargeChange = () => {
	rechargeDialogRef.value.openDialog();
};
const onWithdrawChange = () => {
	withdrawDialogRef.value.openDialog(1);
};
const { copyText } = commonFunction();
const onCopy = (row: { id: string }) => {
	copyText(row.id);
};
const statusDict = [
	{ label: '等待提现', value: 1 },
	{ label: '提现成功', value: 2 },
	{ label: '提现失败', value: 3 },
];
const dateTypeDict = [
	{ label: '创建时间', value: 1 },
	{ label: '支付时间', value: 2 },
	{ label: '提现时间', value: 23 },
];
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
const getTableData = () => {
	state.tableData.loading = true;
	useWalletRecordApi()
		.userWalletList(state.tableData.param)
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				state.tableData.data = res.data.record;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
// 分页改变
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
	getTableData();
};
const stores = useUserInfo();
const useUserApiCollect = useUserApi();
const getUserInfo = async () => {
	await useUserApiCollect.userInfo().then((res: any) => {
		if (res.code === 'SUCCESS') {
			if (res.data?.balance) {
				financeInfo.value = res.data.balance;
			}
			stores.setUserData(res.data);
			Session.set('userInfo', res.data);
			Session.set('userData', res.data);
			stores.setUserInfos();
		}
	});
};
onMounted(() => {
	getTableData();
	getUserInfo();
});
</script>

<style scoped lang="scss">
.my-wallet {
	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;

		.title {
			font-size: 20px;
			color: #6c757d;
		}
		.icon {
			width: 50px;
			height: 50px;
			display: flex;
			justify-content: center;
			align-items: center;
			color: #b5b8ba;
		}
	}
	.balance {
		font-size: 60px;
		color: #000;
		span {
			margin-right: 10px;
			font-size: 20px;
			font-weight: normal;
			color: #6c757d;
		}
	}
	.hint {
		font-size: 16px;
		color: #6c757d;
		display: flex;
		align-items: center;
		margin-top: 10px;
	}
	.tool-btns {
		margin-top: 20px;
		.btn {
			padding: 15px 18px;
			border-radius: 4px;
			background: #0665d0;
			font-size: 14px;
		}
		el-button {
			margin-right: 10px;
		}
	}
}
</style>
