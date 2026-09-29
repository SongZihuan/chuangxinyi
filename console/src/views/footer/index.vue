<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<el-form ref="formRef" :model="ruleForm" label-width="80px" :rules="rules" label-position="left" @submit.native.prevent>
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="ICP备案号" prop="icp1" label-width="125px">
							<el-input v-model="ruleForm.icp1" placeholder="请输入ICP备案号" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="电信经营许可证" prop="icp2" label-width="125px">
							<el-input v-model="ruleForm.icp2" placeholder="请输入电信经营许可证" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="公安备案号" prop="gongan" label-width="125px">
							<el-input v-model="ruleForm.gongan" placeholder="请输入公安备案号" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
						<el-form-item label="版权人" prop="copyright" label-width="125px">
							<el-input v-model="ruleForm.copyright" placeholder="请输入版权人" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-button v-waves type="primary" @click.prevent="onSubmit(formRef)" size="default">确认</el-button>
					</el-col>
				</el-row>
			</el-form>
		</el-card>
	</div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue';
import { type FormRules, type FormInstance } from 'element-plus';
import { ElMessage } from 'element-plus';
import { useFooterApi } from '/@/api/footer/index';
const formRef = ref();
let ruleForm = ref({
	icp1: '',
	icp2: '',
	gongan: '',
	copyright: '',
});
const rules = reactive<FormRules>({});
const onSubmit = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate((valid) => {
		if (valid) {
			useFooterApi()
				.editFooter(ruleForm.value)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						ElMessage({
							type: 'success',
							message: '更新成功',
						});
					}
				});
		} else {
			return false;
		}
	});
};

const getFoot = () => {
	useFooterApi()
		.footer()
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				ruleForm.value = res.data;
			}
		});
};
onMounted(() => {
	getFoot();
});
</script>

<style scoped lang="scss"></style>
