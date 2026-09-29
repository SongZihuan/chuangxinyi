<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="用户ID" prop="id">
						<el-input v-model="ruleForm.id" placeholder="请输入用户ID" clearable disabled></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="标题" prop="subject">
						<el-input v-model="ruleForm.subject" placeholder="请输入标题" clearable></el-input>
					</el-form-item>
				</el-col>
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
          <el-form-item label="可退款期限" prop="price">
            <el-input-number
                v-model="ruleForm.returnDayLimit"
                placeholder="请输入可退款期限"
                clearable
                :min="0"
                style="width: 100%"
                :step="1"
                step-strictly
                controls-position="right"
            ></el-input-number>
          </el-form-item>
        </el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="价格" prop="price">
						<template #label>价格( <el-text class="mx-1" type="danger" tag="b">分</el-text>)</template>
						<el-input-number
							v-model="ruleForm.price"
							placeholder="请输入价格"
							clearable
							:min="0"
							style="width: 100%"
							:step="1"
							step-strictly
							controls-position="right"
						></el-input-number>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="是否分销" prop="distribution">
						<el-radio-group v-model="ruleForm.distribution">
							<el-radio :label="true">是</el-radio>
							<el-radio :label="false">否</el-radio>
						</el-radio-group>
					</el-form-item>
				</el-col>
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
          <el-form-item label="分销是否可提现" prop="canWithdraw">
            <el-radio-group v-model="ruleForm.canWithdraw">
              <el-radio :label="true">是</el-radio>
              <el-radio :label="false">否</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="描述" prop="describe">
						<el-input v-model="ruleForm.describe" type="textarea" placeholder="请输入简短描述"></el-input>
					</el-form-item>
				</el-col>
			</el-row> </el-form
		><template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">确认</el-button>
			</span>
		</template>
		<checkPhone ref="checkPhoneRef" @checkPhoneSuccess="checkPhoneSuccess" />
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { type FormRules, type FormInstance, ElMessage } from 'element-plus';
import checkPhone from '/@/components/checkPhone/index.vue';
import { useDefraytApi } from '/@/api/defray';
const formRef = ref();
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '添加用户消费',
	submitTxt: '',
});
const isShowDialog = ref(false);
const checkPhoneRef = ref();
const rules = reactive<FormRules>({
	subject: [
		{
			required: true,
			message: '请输入标题',
			trigger: 'blur',
		},
	],
	price: [
		{
			required: true,
			message: '请输入价格',
			trigger: 'blur',
		},
	],
	distribution: [
		{
			required: true,
			message: '请选择',
			trigger: 'blur',
		},
	],
  canWithdraw: [
    {
      required: true,
      message: '请选择',
      trigger: 'blur',
    },
  ],
	describe: [
		{
			required: true,
			message: '请输入描述',
			trigger: 'blur',
		},
	],
});
const ruleForm = ref({
	price: 1,
	subject: '',
	describe: '',
	distribution: true,
  canWithdraw: true,
  returnDayLimit: 0,
	id: null,
});
const reset = () => {
	ruleForm.value = {
		price: 1,
		subject: '',
		describe: '',
		distribution: true,
    canWithdraw: false,
    returnDayLimit: 0,
		id: null,
	};
};

const closeDialog = () => {
	dialog.isShowDialog = false;
	isShowDialog.value = false;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			checkPhoneRef.value.openDialog({}, '添加用户订单确认');
		} else {
			return false;
		}
	});
};
const checkPhoneSuccess = (phoneToken: string) => {
	useDefraytApi()
		.adddefary({ ...ruleForm.value, phoneToken: phoneToken })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				ElMessage({
					type: 'success',
					message: '添加用户消费成功',
				});
				closeDialog();
			}
		});
};

const openDialog = async (row?: any) => {
	reset();
	ruleForm.value.id = JSON.parse(JSON.stringify(row)).numberID;
	dialog.isShowDialog = true;
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>
