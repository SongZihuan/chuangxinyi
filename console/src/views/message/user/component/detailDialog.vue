<template>
	<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" destroy-on-close @close="closeDialog">
		<!--文章详情页面，包含标题、内容、发送人、发送link，发布时间-->
		<div class="detail-dialog">
			<div class="detail-dialog__title">
				<span class="detail-dialog__title__text">{{ data.title }}</span>
			</div>
			<div class="detail-dialog__content">
				<div class="detail-dialog__content__text" v-html="data.content"></div>
			</div>
			<div class="detail-dialog__sender">
				<span class="detail-dialog__sender__text" :class="{link:data.senderLink}" @click="onSkip({ link: data.senderLink })">发送人：{{ data.sender }}</span>
			</div>
			<div class="detail-dialog__time">
				<span class="detail-dialog__time__text">发布时间：{{ dayjs.unix(data.createAt).format('YYYY-MM-DD HH:mm:ss') }}</span>
			</div>
		</div>
	</el-dialog>
</template>

<script setup lang="ts" name="DetailDialog">
import { reactive } from 'vue';
import dayjs from 'dayjs';
import { messageUserDataType } from '/@/views/message/user/types';
import { useMessageUserApi } from '/@/api/message/user';
import useSkip from "/@/hooks/useSkip";

const emit = defineEmits(['refresh']);

let dialog = reactive<dialogParams>({
	isShowDialog: false,
	type: 'add',
	title: '',
	submitTxt: '',
});
const { onSkip } = useSkip();
let data = reactive<messageUserDataType>({
	id: '',
	title: '',
	content: '',
	sender: '',
	createAt: 0,
	readAt: 0,
	senderLink: '',
});
// 打开弹窗
const openDialog = (row: messageUserDataType) => {
	dialog.isShowDialog = true;
	dialog.title = '';
	dialog.submitTxt = '查看';
	data = row;
	read(row.id);
};
// 已读
const read = (id: string) => {
	useMessageUserApi()
		.readMessage({ id })
		.then((res: any) => {
			if (res) {
				emit('refresh');
			}
		});
};
const closeDialog = () => {
	dialog.isShowDialog = false;
};
defineExpose({
	openDialog,
});
</script>

<style scoped lang="scss">
.detail-dialog {
	min-height: 200px;
	&__title {
		text-align: center;
		font-size: 20px;
		font-weight: bold;
		margin-bottom: 20px;
		&__text {
			color: #333;
		}
	}
	&__content {
		&__text {
			color: #333;
			line-height: 1.5;
			font-size: 14px;
			& p {
				margin-bottom: 10px;
			}
		}
	}
	&__sender {
		margin-top: 20px;
		&__text {
			color: #333;
			font-size: 14px;
			margin-right: 10px;
		}
    .link {
      color: #409eff;
      cursor: pointer;
    }
	}
	&__time {
		margin-top: 20px;
		&__text {
			color: #333;
			font-size: 14px;
		}
	}
}
</style>
