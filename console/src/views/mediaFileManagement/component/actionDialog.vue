<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="850px" destroy-on-close @close="closeDialog">
		<el-form ref="formRef" :model="ruleForm" size="default" label-width="80px" :rules="rules">
			<el-row :gutter="35">
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
          <el-input v-model="ruleForm.fid" placeholder="输入名字" :disabled="nameDisable"/>
        </el-col>
				<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
					<el-form-item label="文件上传">
						<el-upload
							ref="upload"
							class="upload-demo"
							:show-file-list="true"
							:limit="1"
							:on-change="handleChange"
							:on-exceed="handleExceed"
							:auto-upload="false"
						>
							<template #trigger>
								<el-button type="primary">选择</el-button>
							</template>
							<template #tip>
								<div class="el-upload__tip text-red">限制1个文件，新文件将覆盖旧文件</div>
							</template>
						</el-upload>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="closeDialog" size="default">取 消</el-button>
				<el-button v-if="dialog.type==='edit'" type="primary" @click="edit(formRef)" size="default">{{ dialog.submitTxt }}</el-button>
				<el-button v-else type="primary" @click="onSubmit(formRef)" size="default">{{ dialog.submitTxt }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script setup lang="ts" name="invoiceDialog">
import { reactive, ref, onMounted } from 'vue';
import {
  type FormRules,
  type FormInstance,
  UploadRawFile,
  UploadProps,
  UploadInstance,
  genFileId,
  ElMessageBox
} from 'element-plus';
import { ElMessage } from 'element-plus';
import { useMediaFileApi } from '/@/api/mediaFileManagement';

const nameDisable = ref(false)
const formRef = ref();
const emit = defineEmits(['refresh']);
let ruleForm = ref<any>({
	fid: '',
	file: '',
});

const rules = reactive<FormRules>({});

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '新增文件',
	submitTxt: '新增',
});

const openDialog = (row:any,type:string) => {
  dialog.type = type;
  if(type === 'edit') {
    dialog.title = '覆盖文件';
    dialog.submitTxt = '覆盖';
    ruleForm.value = {
      fid: row.fid,
      file: '',
    };
    nameDisable.value = true
  } else {
    dialog.title = '新增文件';
    dialog.submitTxt = '新增';
    reset();
  }
  dialog.isShowDialog = true;
  nameDisable.value = false
};
const upload = ref<UploadInstance>();

const handleChange: UploadProps['onChange'] = (file) => {
  // 去除上传文件的后缀名
  file.name = file.name.split('.')[0];
  if(dialog.type === 'edit') {
    ruleForm.value.file = file.raw;
    return;
  }
	ruleForm.value.file = file.raw;
};
const handleExceed: UploadProps['onExceed'] = (files) => {
	upload.value!.clearFiles();
	const file = files[0] as UploadRawFile;
	file.uid = genFileId();
	upload.value!.handleStart(file);
};
//重置
const reset = () => {
	ruleForm.value = {
		fid: '',
		file: '',
	};
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const edit = (formEl: FormInstance | undefined) => {
  if (ruleForm.value.file === '' || ruleForm.value.fid === '') {
    ElMessage({
      type: 'error',
      message: '请上传文件',
    });
    return;
  }
  ElMessageBox.confirm(`此操作将覆盖文件, 是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      onSubmit(formEl);
    })
    .catch(() => {
      ElMessage.info('已取消覆盖');
    });
};
const onSubmit = (formEl: FormInstance | undefined) => {
	if (ruleForm.value.file === '' || ruleForm.value.fid === '') {
		ElMessage({
			type: 'error',
			message: '请上传文件',
		});
		return;
	}
	if (!formEl) return;
	let fromData = new FormData();
	fromData.append('file', ruleForm.value.file);
	fromData.append('fid', ruleForm.value.fid);
	formEl.validate((valid) => {
		if (valid) {
			useMediaFileApi()
				.addMediaFile(ruleForm.value)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						closeDialog();
						ElMessage({
							type: 'success',
							message: '新增文件成功',
						});
						emit('refresh');
					}
				});
		} else {
			return false;
		}
	});
};

onMounted(() => {});
defineExpose({
	openDialog,
	closeDialog,
});
</script>

<style lang="scss" scoped>
.ifr {
	width: 100%;
	height: 560px;
}

.tip {
	display: flex;
	padding: 10px 10px;
}
</style>
