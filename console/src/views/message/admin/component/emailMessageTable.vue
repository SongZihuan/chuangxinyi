<template>
	<div class="search-header mb15">
		<el-form :inline="true" :model="tableData.param">
			<el-form-item label="邮箱">
				<el-input v-model="tableData.param.email" placeholder="请输入邮箱"></el-input>
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
				<el-button type="primary" class="ml10" @click="handleSearch">
					<el-icon>
						<ele-Search />
					</el-icon>
					查询
				</el-button>
				<el-button type="success" class="ml10" @click="onSend">
					<SvgIcon name="my-send"></SvgIcon>
					发送邮件
				</el-button>
			</el-form-item>
		</el-form>
	</div>
	<el-table :data="tableData.data" v-loading="tableData.loading" style="width: 100%">
		<el-table-column prop="email" label="邮箱" width="180" show-overflow-tooltip align="left"></el-table-column>
		<el-table-column prop="subject" label="主题" width="180" show-overflow-tooltip align="center"></el-table-column>
		<el-table-column prop="content" label="内容" min-width="180" show-overflow-tooltip align="center"></el-table-column>
		<el-table-column prop="success" label="是否成功" width="80" show-overflow-tooltip align="center">
			<template #default="scope">
				<el-tag v-if="scope.row.success" type="success">成功</el-tag>
				<el-tag v-else type="danger">失败</el-tag>
			</template>
		</el-table-column>
		<el-table-column prop="errorMsg" label="错误信息" width="180" show-overflow-tooltip align="center">
			<template #default="scope">
				<el-tag v-if="scope.row.errorMsg" type="danger">{{ scope.row.errorMsg }}</el-tag>
				<el-tag v-else type="info">无</el-tag>
			</template>
		</el-table-column>
		<el-table-column prop="createAt" label="创建时间" width="180" show-overflow-tooltip align="right">
			<template #default="scope">
				<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
			</template>
		</el-table-column>
	</el-table>
	<el-dialog v-model="dialogVisible" title="发送邮件" width="850px">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="主题" prop="subject">
						<el-input v-model="ruleForm.subject" placeholder="请输入主题" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="用户ID" prop="id">
						<el-input v-model="ruleForm.id" placeholder="请输入用户ID" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="内容" prop="content">
						<el-input v-model="ruleForm.content" placeholder="请输入内容" clearable type="textarea" :autosize="{ minRows: 2, maxRows: 6 }"></el-input>
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
import { messageAdminState } from '/@/views/message/admin/types';
import { type FormRules, type FormInstance, ElMessage } from 'element-plus';
import { useMessageAdminApi } from '/@/api/message/admin';
import { reactive, ref } from 'vue';
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
	subject: '',
	content: '',
	id: null,
});
const formRef = ref();
const rules = reactive<FormRules>({
	subject: [
		{
			required: true,
			message: '请输入主题',
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
	content: [
		{
			required: true,
			message: '请输入内容',
			trigger: 'blur',
		},
	],
});
const dialogVisible = ref<boolean>(false);
const closeDialog = () => {
	dialogVisible.value = false;
};
const onSend = () => {
	ruleForm.value.subject = '';
	ruleForm.value.content = '';
	ruleForm.value.id = null;
	dialogVisible.value = true;
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useMessageAdminApi()
				.sendEmail({ ...ruleForm.value, id: Number(ruleForm.value.id) })
				.then((res: any) => {
					if (res.code === "SUCCESS" && !res.data.have) {
						ElMessage({
							type: 'error',
							message: '当前用户暂未绑定邮箱',
						});
						return;
					}
					if (res.code === "SUCCESS" && res.data.success) {
						ElMessage({
							type: 'success',
							message: '发送邮件成功!',
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
