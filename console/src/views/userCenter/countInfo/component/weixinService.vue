<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="500px" destroy-on-close @close="closeDialog">
		<div class="wx-service-container">
			<div class="wx-service-box">
				<img class="img" :src="fuwuhao" alt="" />
				<div class="success" v-if="(props.wxStatus !== 2 && props.wxStatus !== 4) && !isTimeout">
					<div class="success-img"></div>
					<div class="success-title">扫码成功</div>
					<div class="count-down">扫码成功，本窗口将在{{ autoCloseTime }}s之后，自动关闭</div>
				</div>
				<div class="error" v-else-if="(props.wxStatus === 2 || props.wxStatus === 4) && isTimeout">
					<div class="error-img"></div>
					<div class="error-title">扫码超时</div>
					<div class="refresh" @click="refreshCode"><SvgIcon name="ele-RefreshRight"></SvgIcon>重新扫码</div>
				</div>
			</div>
      <div class="hint">请打开微信扫码关注服务号，若已关注请点击按钮绑定（账号>绑定账号）</div>
		</div>
	</el-dialog>
</template>
<script setup lang="ts">
import {onUnmounted, reactive, ref, watch} from 'vue';
import useFile from "/@/hooks/useFile";

const fuwuhao = useFile().getFile("fuwuhao");

interface Props {
	wxStatus: number;
}

const props = withDefaults(defineProps<Props>(), {});
const emit = defineEmits(['refresh']);

const dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: '',
	title: '绑定微信服务号',
	submitTxt: '确认',
});
// 有效时间 15分钟
const validTime = ref<number>(15 * 60 * 1000);
// 是否超时
const isTimeout = ref<boolean>(false);
// 自动退出时间
const autoCloseTime = ref<number>(3);
const openDialog = () => {
	dialog.isShowDialog = true;
	refreshCode();
};

const closeDialog = () => {
	dialog.isShowDialog = false;
};
// 刷新
const intervalId = ref<any>(null);
watch(
    () => props.wxStatus,
    (val) => {
      if (val !== 2 && val !== 4) {
        clearInterval(intervalId.value);
        const timer = setInterval(() => {
          autoCloseTime.value -= 1;
          if (autoCloseTime.value <= 0) {
            clearInterval(timer);
            closeDialog();
          }
        }, 1000);
      }
    }
)
// 检查是否扫码成功
const checkWxStatus = () => {
	intervalId.value = setInterval(() => {
		validTime.value -= 3000;
		if (validTime.value <= 0) {
			clearInterval(intervalId.value);
			isTimeout.value = true;
		}
		emit('refresh');
	}, 3000);
};
const refreshCode = () => {
	validTime.value = 15 *60 * 1000;
	isTimeout.value = false;
	checkWxStatus();
};
onUnmounted(() => {
	clearInterval(intervalId.value);
});
defineExpose({
	openDialog,
});
</script>
<style scoped lang="scss">
.wx-service-container {
	display: flex;
  flex-direction: column;
	justify-content: center;
	align-items: center;
	.wx-service-box {
		position: relative;
		left: 0;
		top: 0;
		width: 300px;
		height: 300px;
		.img {
			width: 100%;
			height: 100%;
		}
		.success,
		.error {
			position: absolute;
			left: 0;
			top: 0;
			width: 100%;
			height: 100%;
			background: rgba(255, 255, 255, 0.98);
			background-size: 100% 100%;
			.success-img {
				width: 120px;
				height: 120px;
				margin: 60px auto 0;
				background: url('/@/assets/check-circle.png') no-repeat center;
				background-size: 100% 100%;
			}
			.refresh {
				display: flex;
				justify-content: center;
				align-items: center;
				font-size: 14px;
				color: #999999;
				margin-top: 10px;
				cursor: pointer;
				user-select: none;
			}
			.error-img {
				width: 120px;
				height: 120px;
				margin: 60px auto 0;
				background: url('/@/assets/x-circle.png') no-repeat center;
				background-size: 100% 100%;
			}
			.error-title,
			.success-title {
				width: 100%;
				text-align: center;
				font-size: 16px;
				color: #333333;
				margin-top: 10px;
			}
			.count-down {
				width: 100%;
				text-align: center;
				font-size: 14px;
				color: #999999;
				margin-top: 10px;
			}
		}
	}
  .hint{
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 14px;
    color: #999999;
    margin-bottom: 20px;
  }
}
</style>
