<template>
	<div class="search-header mb15">
		<el-form :inline="true" :model="tableData.param">
			<el-form-item>
				<el-input v-model="tableData.param.id" placeholder="请输入ID"></el-input>
			</el-form-item>
			<el-form-item>
				<el-input v-model="tableData.param.uid" placeholder="请输入用户ID"></el-input>
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
					发送站内信
				</el-button>
			</el-form-item>
		</el-form>
	</div>
	<el-table :data="tableData.data" v-loading="tableData.loading" style="width: 100%">
		<el-table-column prop="userID" label="用户ID" width="120" show-overflow-tooltip align="left"></el-table-column>
		<el-table-column prop="title" label="标题" width="180" show-overflow-tooltip align="center"></el-table-column>
		<el-table-column prop="content" label="内容" show-overflow-tooltip min-width="180" align="center"></el-table-column>
		<el-table-column prop="sender" label="发送者" width="180" show-overflow-tooltip align="center">
			<template #default="scope">
				<span :class="{link:scope.row.senderLink}" v-if="scope.row.sender" @click="onSkip({ link: scope.row.senderLink })">{{ scope.row.sender }}</span>
			</template>
		</el-table-column>
		<el-table-column prop="sender" label="发送链接" width="180" show-overflow-tooltip align="center">
			<template #default="scope">
				<span :class="{link:scope.row.senderLink}" @click="onSkip({ link: scope.row.senderLink })">{{ scope.row.senderLink }}</span>
			</template>
		</el-table-column>
		<el-table-column prop="readAt" label="阅读时间" width="180" show-overflow-tooltip align="center">
			<template #default="scope">
				<span v-if="scope.row.readAt">{{ dayjs.unix(scope.row.readAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
        <span v-else>未阅读</span>
			</template>
		</el-table-column>
		<el-table-column prop="createAt" label="发送时间" width="180" show-overflow-tooltip align="center">
			<template #default="scope">
				<span v-if="scope.row.createAt">{{ dayjs.unix(scope.row.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
			</template>
		</el-table-column>
	</el-table>
	<!-- 发送站内信弹窗 -->
	<el-dialog v-model="dialogVisible" title="发送站内信" width="850px">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="70px" :rules="rules">
			<el-row :gutter="35">
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="标题" prop="title">
						<el-input v-model="ruleForm.title" placeholder="请输入标题" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="消息" prop="message">
						<el-input v-model="ruleForm.message" placeholder="请输入消息" clearable type="textarea" :autosize="{ minRows: 2, maxRows: 6 }"></el-input>
					</el-form-item>
				</el-col>
				<!--        sendLink-->
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="链接" prop="senderLink">
						<el-input v-model="ruleForm.senderLink" placeholder="请输入链接" clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="用户ID" prop="id">
						<el-input v-model="ruleForm.id" placeholder="请输入用户ID" clearable></el-input>
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
import { reactive, ref, onMounted } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { useMessageAdminApi } from '/@/api/message/admin';
import { ElMessage } from 'element-plus';
import useSkip from '/@/hooks/useSkip';
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
	title: null,
	message: null,
	id: null,
	senderLink: null,
});
const formRef = ref();
const { onSkip } = useSkip();
const rules = reactive<FormRules>({
	title: [
		{
			required: true,
			message: '请输入标题',
			trigger: 'blur',
		},
	],
	message: [
		{
			required: true,
			message: '请输入消息',
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
	senderLink: [
		{
			required: true,
			message: '请输入发送链接',
			trigger: 'blur',
		},
	],
});

const dialogVisible = ref<boolean>(false);
const closeDialog = () => {
	dialogVisible.value = false;
};
const onSend = () => {
	ruleForm.value.title = null;
	ruleForm.value.message = null;
	ruleForm.value.id = null;
	dialogVisible.value = true;
};

const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useMessageAdminApi()
				.sendMsg({ ...ruleForm.value, id: Number(ruleForm.value.id) })
				.then((res: any) => {
					if (res.code === "SUCCESS" && res.data.success) {
						ElMessage({
							type: 'success',
							message: '发送站内信成功!',
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
onMounted(() => {});
</script>
<style scoped lang="scss">
.link {
	color: #409eff;
	cursor: pointer;
}
</style>
