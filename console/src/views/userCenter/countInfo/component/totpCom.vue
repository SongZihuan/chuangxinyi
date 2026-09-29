<template>
	<!-- 2fa -->
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<div class="login-scan-container">
			<div v-if="qrcodeUrl">
				<div ref="qrcodeRef" class="mt-5"></div>
				<div class="font12 mt20 login-msg">
					<i class="iconfont icon-saoyisao mr5"></i>
					<div class="scan-content mt-5">
						<el-text size="large" class="copy-text" v-for="(item, index) in fomat(ruleForm.secret)" :key="index">
							{{ item }}
						</el-text>
					</div>
					<el-button-group class="mt-5">
            <el-button type="danger" class="copy-tip" @click="handleCopy">复制密钥</el-button>
            <el-button type="success" class="copy-tip" @click="showCodeDialog">确认绑定</el-button>
          </el-button-group>
				</div>
			</div>
			<!-- 验证码弹窗 -->
			<el-dialog v-model="dialogVisible" title="输入2FA验证码" width="340px" center destroy-on-close>
				<el-input v-model="code" placeholder="2FA验证码" />
				<template #footer>
					<span class="dialog-footer">
						<el-button @click="dialogVisible = false">取消</el-button>
						<el-button type="primary" @click="bind"> 确认 </el-button>
					</span>
				</template>
			</el-dialog>
		</div>
	</el-dialog>
</template>

<script setup lang="ts" name="userTotpCom">
import { ref, nextTick, reactive } from 'vue';
import commonFunction from '/@/utils/commonFunction';
import QRCode from 'qrcodejs2-fixes';
import { ElMessage } from 'element-plus';
import { useGoogleAuthApi } from '/@/api/Totp/index';
import type { totpConfigTypes } from '/@/api/Totp/types';

const { copyText } = commonFunction();
const useGoogleAuthApiCollect = useGoogleAuthApi();
const emit = defineEmits(['refresh']);

interface Props {
	totpForm: totpConfigTypes;
}

const props = withDefaults(defineProps<Props>(), {});
const ruleForm = ref<totpConfigTypes>(props.totpForm);
const qrcodeRef = ref<HTMLElement | null>(null);
const qrcodeUrl = ref<string>('');
const dialogVisible = ref<boolean>(false);
const code = ref<string>();
const fomat = (key: string) => {
	return key.match(/(.{16})/g) || [];
};
const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '绑定的2FA',
	submitTxt: '确认',
});
const openDialog = () => {
	dialog.isShowDialog = true;
	reset();
	generateCode();
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
const reset = () => {
	ruleForm.value.secret = '';
	qrcodeUrl.value = '';
	dialogVisible.value = false;
};
const showCodeDialog = () => {
	dialogVisible.value = true;
};
const handleCopy = () => {
	copyText(ruleForm.value.secret);
};
//生成用户绑定二维码
const generateCode = async () => {
	if (!ruleForm.value.phone) {
		ElMessage.error('绑定2FA失败');
		return;
	}
	await useGoogleAuthApiCollect.totpUrl(ruleForm.value).then((res: any) => {
		if (res.code === "SUCCESS") {
			qrcodeUrl.value = res.data.url;
			ruleForm.value.secret = res.data.secret;
			nextTick(() => {
				(<HTMLElement>qrcodeRef.value).innerHTML = '';
				new QRCode(qrcodeRef.value, {
					text: qrcodeUrl.value,
					width: 200,
					height: 200,
					colorDark: '#000000',
					colorLight: '#ffffff',
				});
			});
		}
	});
};
//用户输入验证码
const bind = async () => {
	await useGoogleAuthApiCollect.bindTotp({ code: code.value, secret: ruleForm.value.secret }).then((res: any) => {
		if (res.code === "SUCCESS") {
			ElMessage.success('绑定双因素验证器成功');
			dialogVisible.value = false;
			closeDialog();
			emit('refresh');
		}
	});
};
defineExpose({
	openDialog,
});
</script>

<style lang="scss" scoped>
.login-content-submit {
	width: 100%;
	letter-spacing: 2px;
	font-weight: 300;
	margin-top: 15px;
}

.scan-tip {
	color: var(--el-color-primary);
	cursor: pointer;
	margin-bottom: 10px;
}

.copy {
	margin-top: 15px;
	min-height: 180px;
}
.scan-content {
	display: flex;
	flex-wrap: wrap;
	justify-content: center;
	font-family: 'Courier New', Courier, monospace;
	.copy-text {
		min-width: 180px;
		margin: 10px 5px 0 5px;
		text-align: left;
		cursor: pointer;
		font-weight: bold;
	}
}

.copy-tip {
	margin-top: 15px;
}

.login-scan-container {
	padding: 0 20px 20px;
	display: flex;
	flex-direction: column;
	text-align: center;
	min-height: 430px;
	animation-delay: 0.1s;

	:deep(img) {
		margin: auto;
	}

	.login-msg {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		color: var(--el-text-color-placeholder);

		animation-delay: 0.2s;
	}
}
</style>
