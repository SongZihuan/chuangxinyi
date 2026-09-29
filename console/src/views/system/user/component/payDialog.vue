<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="900px" destroy-on-close @close="closeDialog">
		<div class="search-header mb15">
			<el-form :inline="true" :model="state.tableData.param">
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
					<el-button type="primary" @click="handleSearch">
						<el-icon>
							<ele-Search />
						</el-icon>
						查询
					</el-button>
				</el-form-item>
			</el-form>
		</div>
		<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
			<el-table-column prop="tradeID" label="订单ID" show-overflow-tooltip min-width="240" align="left"></el-table-column>
			<el-table-column prop="subject" label="订单名称" show-overflow-tooltip align="center" min-width="140"></el-table-column>
			<el-table-column prop="cny" label="订单金额" show-overflow-tooltip align="center" min-width="140">
				<template #default="scope">
					<el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.cny) }}元</el-text>
				</template>
			</el-table-column>
			<el-table-column prop="tradeStatus" label="支付状态" show-overflow-tooltip align="center" width="100">
				<template #default="scope">
					<el-tag type="primary" v-if="scope.row.tradeStatus === 1">等待支付</el-tag>
					<el-tag type="success" v-else-if="scope.row.tradeStatus === 2 || scope.row.tradeStatus === 3">支付成功</el-tag>
					<el-tag type="danger" v-else-if="scope.row.tradeStatus === 4">支付关闭</el-tag>
					<el-tag type="info" v-else-if="scope.row.tradeStatus === 5">等待退款</el-tag>
					<el-tag type="success" v-else-if="scope.row.tradeStatus === 6">退款成功</el-tag>
					<el-tag type="info" v-else-if="scope.row.tradeStatus === 7">退款失败</el-tag>
          <el-tag type="info" v-else-if="scope.row.tradeStatus === 8">单边退款成功</el-tag>
					<el-tag type="danger" v-else>支付未找到</el-tag>
				</template>
			</el-table-column>
			<el-table-column prop="get" label="实际获得金额" show-overflow-tooltip align="center" width="130">
				<template #default="scope">
					<el-text class="mx-1" type="danger" tag="b">{{ formatAmount(scope.row.get) }}元</el-text>
				</template>
			</el-table-column>
			<el-table-column prop="payWay" label="支付方式" show-overflow-tooltip align="center" width="130"></el-table-column>
			<el-table-column prop="createTime" label="支付时间" show-overflow-tooltip width="170" align="center">
				<template #default="scope">
					<span v-if="scope.row.payAt">{{ dayjs.unix(scope.row.payAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
				</template>
			</el-table-column>
			<el-table-column label="退款时间" show-overflow-tooltip width="170" align="center">
				<template #default="scope">
					<span v-if="scope.row.refundAt">{{ dayjs.unix(scope.row.refundAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
				</template>
			</el-table-column>
			<el-table-column prop="createTime" label="创建时间" show-overflow-tooltip width="170" align="center">
				<template #default="scope">
					<span>{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
				</template>
			</el-table-column>
			<el-table-column label="操作" width="200" fixed="right" align="center">
				<template #default="scope">
					<el-button text type="primary" @click="onCheck(scope.row)" v-if="isCheck(scope.row)">充值审核</el-button>
					<el-button
						text
						type="danger"
						@click="onRefund(scope.row)"
						v-if="
							(scope.row.tradeStatus === 2 || scope.row.tradeStatus === 3 || scope.row.tradeStatus === 7)
						"
						>
            退款
          </el-button>
          <el-button
              text
              type="danger"
              @click="onRefundInside(scope.row)"
              v-if="
							(scope.row.tradeStatus === 2 || scope.row.tradeStatus === 3 || scope.row.tradeStatus === 7)
						"
          >
            单边退款
          </el-button>
					<el-button text type="danger" @click="onCheckRefund(scope.row, true)" v-if="scope.row.tradeStatus === 5">审核通过</el-button>
          <el-button text type="danger" @click="onCheckRefund(scope.row, false)" v-if="scope.row.tradeStatus === 5">审核驳回</el-button>
				</template>
			</el-table-column>
		</el-table>
		<el-pagination hide-on-single-page
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
		<el-dialog title="充值审核" v-model="isShowDialog" width="600px" height="500px" center destroy-on-close @close="isShowDialog = false">
			<el-form ref="formRef" :model="ruleForm" size="default" label-width="140px" :rules="rules">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="优惠政策" prop="sex">
							<el-radio-group v-model="ruleForm.way">
								<el-radio :label="1">是</el-radio>
								<el-radio :label="2">否</el-radio>
							</el-radio-group>
						</el-form-item>
					</el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="是否成功" prop="sex">
              <el-radio-group v-model="ruleForm.success">
                <el-radio :label="true">是</el-radio>
                <el-radio :label="false">否</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" v-if="ruleForm.way === 1">
						<el-form-item label="充值金额(分)" prop="cny">
							<template #label>充值金额( <el-text class="mx-1" type="danger" tag="b">分</el-text>)</template>
							<el-input-number
								v-model="ruleForm.cny"
								placeholder="请输入充值金额"
								clearable
								:disabled="ruleForm.way === 1"
								controls-position="right"
								style="width: 250px"
							></el-input-number>
						</el-form-item>
					</el-col>

					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" v-if="ruleForm.way === 2">
						<el-form-item label="实际获得金额" prop="get">
							<template #label>实际获得金额( <el-text class="mx-1" type="danger" tag="b">分</el-text>)</template>
							<el-input-number
								v-model="ruleForm.get"
								controls-position="right"
								placeholder="请输入实际获得金额"
								clearable
								style="width: 250px"
								:min="0"
							></el-input-number>
						</el-form-item>
					</el-col>
				</el-row> </el-form
			><template #footer>
				<span class="dialog-footer">
					<el-button @click="isShowDialog = false" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit(formRef)" size="default">确认</el-button>
				</span>
			</template></el-dialog
		>
		<checkPhone ref="checkPhoneRef" @checkPhoneSuccess="checkPhoneSuccess" />
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import dayjs from 'dayjs';
import type { rechargeStateTypes } from '/@/api/recharge/types';
import { type FormRules, type FormInstance, ElMessage } from 'element-plus';
import checkPhone from '/@/components/checkPhone/index.vue';
import { formatAmount } from '/@/utils/formatAmount';
import { useDefraytApi } from '/@/api/defray';

const info = ref<any>({});
const formRef = ref();
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '用户充值',
	submitTxt: '',
});
const checkPhoneRef = ref();
const isShowDialog = ref(false);
const currentType = ref(1); //1 充值审核  2退款处理  3-单边退款
const dateTypeDict = ref([
	{ label: '创建时间', value: 1 },
	{ label: '支付时间', value: 2 },
]);
const rules = reactive<FormRules>({
	get: [
		{
			required: true,
			message: '请输入实际获得金额',
			trigger: 'blur',
		},
	],
	cny: [
		{
			required: true,
			message: '请输入充值金额',
			trigger: 'blur',
		},
	],
});
const ruleForm = ref({
	cny: '',
	get: -1,
	way: 1,
  success: true,
})

const state = reactive<rechargeStateTypes>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			range: [],
			page: 1,
			pagesize: 20,
			timetype: '',
			starttime: null,
			endtime: null,
			uid: '',
		},
	},
});
const payWay = ['支付宝电脑支付', '支付宝手机支付', '微信Native支付', '微信H5支付', '分销收益', '优惠'];
const isCheck = (row: any) => {
	return !payWay.some((item: string) => item === row.payWay) && row.tradeStatus == 1;
};
const handleSearch = () => {
	state.tableData.param.page = 1;
	getTableData();
};
const currentInfo = ref();
const currentSuccess = ref(false)
const onCheck = (row: any) => {
	isShowDialog.value = true;
	ruleForm.value.cny = row.cny;
	currentInfo.value = row;
	currentType.value = 1;
};

const onRefund = (row: any) => {
	currentInfo.value = row;
	checkPhoneRef.value.openDialog(row, '退款确认');
	currentType.value = 2;
};

const onRefundInside = (row: any) => {
  currentInfo.value = row;
  checkPhoneRef.value.openDialog(row, '单边退款确认');
  currentType.value = 3;
};

const onCheckRefund = (row: any, success: boolean) => {
  currentInfo.value = row;
  checkPhoneRef.value.openDialog(row, '退款审核');
  currentType.value = 4;
  currentSuccess.value = success
};

const checkPhoneSuccess = (phoneToken: string) => {
	if (currentType.value === 1) {
		let get: number = 0;
		ruleForm.value.way == 1 ? (get = -1) : (get = parseInt(ruleForm.value.get));
		useDefraytApi()
			.adminpayProcess({ id: currentInfo.value.tradeID, get: get, phoneToken: phoneToken, success: ruleForm.value.success, })
			.then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage({
						type: 'success',
						message: '充值成功',
					});
					isShowDialog.value = false;
					getTableData();
				}
			});
	} else if (currentType.value === 2) {
		useDefraytApi()
			.adminRefund({ tradeID: currentInfo.value.tradeID, success: true, phoneToken: phoneToken })
			.then((res: any) => {
				if (res.code === "SUCCESS") {
					ElMessage({
						type: 'success',
						message: '退款成功',
					});
          getTableData();
				}
			});
	} else if (currentType.value === 3) {
    useDefraytApi()
        .adminRefundInside({ tradeID: currentInfo.value.tradeID, success: true, phoneToken: phoneToken })
        .then((res: any) => {
          if (res.code === "SUCCESS") {
            ElMessage({
              type: 'success',
              message: '退款成功',
            });
            getTableData();
          }
        });
  } else if (currentType.value === 4) {
    useDefraytApi()
        .adminRefund({ tradeID: currentInfo.value.tradeID, success: currentSuccess.value, phoneToken: phoneToken })
        .then((res: any) => {
          if (res.code === "SUCCESS") {
            ElMessage({
              type: 'success',
              message: '退款成功',
            });
            getTableData();
          }
        });
  }
};
const getTableData = () => {
	state.tableData.loading = true;
	useDefraytApi()
		.adminPayList({ ...state.tableData.param, uid: info.value.id })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.pay;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
const closeDialog = () => {
	dialog.isShowDialog = false;
	isShowDialog.value = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			checkPhoneRef.value.openDialog({}, '充值审核确认');
		} else {
			return false;
		}
	});
};

const openDialog = async (row?: any) => {
	info.value = JSON.parse(JSON.stringify(row));
	await getTableData();
	dialog.isShowDialog = true;
};
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
	getTableData();
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>
