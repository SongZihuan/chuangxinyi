<template>
	<div class="search-header mb15">
		<el-form :inline="true" :model="tableData.param">
			<el-form-item>
				<el-input v-model="tableData.param.phone" placeholder="请输入手机号"></el-input>
			</el-form-item>
			<el-form-item>
				<el-date-picker
					v-model="tableData.param.range"
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
				<el-button type="success" class="ml10" @click="onSend">
					<SvgIcon name="my-send"></SvgIcon>
					发送短信
				</el-button>
			</el-form-item>
		</el-form>
	</div>
	<el-table :data="tableData.data" v-loading="tableData.loading" style="width: 100%">
		<el-table-column prop="phone" label="手机号" width="120" show-overflow-tooltip align="left"></el-table-column>
		<el-table-column prop="sig" label="签名" width="120" show-overflow-tooltip align="center"></el-table-column>
		<el-table-column prop="template" label="模板" width="180" show-overflow-tooltip align="center"></el-table-column>
		<el-table-column prop="errorMsg" label="错误信息" min-width="180" show-overflow-tooltip align="center">
			<template #default="scope">
				<el-tag v-if="scope.row.errorMsg" type="danger">{{ scope.row.errorMsg }}</el-tag>
				<el-tag v-else type="info">无</el-tag>
			</template>
		</el-table-column>
		<el-table-column prop="templateParam" label="模板参数" width="180" show-overflow-tooltip align="center">
			<template #default="scope">
				<template v-if="scope.row.templateParam">
					<span v-for="item in scope.row.templateParam" :key="item" type="primary">
						<el-text tag="b">{{ item.label }}:</el-text>:{{ item.value }}
					</span>
				</template>
				<el-tag v-else type="danger">无</el-tag>
			</template>
		</el-table-column>
		<el-table-column prop="success" label="是否成功" width="80" show-overflow-tooltip align="center">
			<template #default="scope">
				<el-tag v-if="scope.row.success" type="success">成功</el-tag>
				<el-tag v-else type="danger">失败</el-tag>
			</template>
		</el-table-column>
		<el-table-column prop="createAt" label="创建时间" width="180" show-overflow-tooltip align="right">
			<template #default="scope">
				<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
			</template>
		</el-table-column>
	</el-table>
	<!-- 发送短信弹窗 -->
	<el-dialog v-model="dialogVisible" title="发送短信" width="850px">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="模版ID" prop="templateID">
						<el-input v-model="ruleForm.templateID" placeholder="请输入模版ID" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="用户ID" prop="id">
						<el-input v-model="ruleForm.id" placeholder="请输入用户ID" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="短信签名" prop="sig">
						<el-input v-model="ruleForm.sig" placeholder="请输入短信签名" clearable type="textarea" :autosize="{ minRows: 2, maxRows: 6 }"></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="标签列表" prop="data" v-for="(item, index) in ruleForm.data" :key="index">
						<el-input v-model="item.label" placeholder="标签" style="width: 35%"></el-input>
						<el-input v-model="item.value" placeholder="标签值" style="width: 35%" class="ml10"></el-input>
						<el-button type="primary" @click="handleAdd" class="ml5" text v-if="index === 0">新增</el-button>
						<el-button type="danger" @click="handleDel(index)" text v-else>删除</el-button>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button type="primary" @click="onSubmit(formRef)" size="default">发送</el-button>
			</span>
		</template></el-dialog
	>
</template>
<script setup lang="ts">
import dayjs from 'dayjs';
import { reactive, ref } from 'vue';
import { messageAdminState } from '/@/views/message/admin/types';
import { type FormRules, type FormInstance, ElMessage } from 'element-plus';
import { useMessageAdminApi } from '/@/api/message/admin';
interface Props {
	data: messageAdminState;
}
const props = withDefaults(defineProps<Props>(), {});
const emit = defineEmits(['search', 'update']);
const tableData = reactive(props.data.tableData);
const handleSearch = () => {
	emit('search', tableData.param);
};
let ruleForm = ref({
	templateID: '',
	sig: '',
	data: [{ label: '', value: '' }],
	id: null,
});
const formRef = ref();
const rules = reactive<FormRules>({
	templateID: [
		{
			required: true,
			message: '请输入模版ID',
			trigger: 'blur',
		},
	],
	sig: [
		{
			required: true,
			message: '请输入短信',
			trigger: 'blur',
		},
	],
	id: [
		{
			required: true,
			message: '请输入用户ID',
			trigger: 'blur',
		},
	],
});

const dialogVisible = ref<boolean>(false);
const closeDialog = () => {
	dialogVisible.value = false;
};
const onSend = () => {
	ruleForm.value.templateID = '';
	ruleForm.value.sig = '';
	ruleForm.value.id = null;
	ruleForm.value.data = [{ label: '', value: '' }];
	dialogVisible.value = true;
};
const handleAdd = () => {
	ruleForm.value.data.push({ value: '', label: '' });
};
const handleDel = (index: number) => {
	ruleForm.value.data.splice(index, 1);
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useMessageAdminApi()
				.sendSms({ ...ruleForm.value, id: Number(ruleForm.value.id) })
				.then((res: any) => {
					if (res.code === "SUCCESS" && res.data.success) {
						ElMessage({
							type: 'success',
							message: '发送短信成功!',
						});
						dialogVisible.value = false;
            emit('update');
					} else {
						ElMessage({
							type: 'error',
							message: '发送失败',
						});
					}
				});
		} else {
			return false;
		}
	});
};
</script>
<style scoped lang="scss"></style>
