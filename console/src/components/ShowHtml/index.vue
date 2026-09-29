<template>
	<el-dialog :title="props.title" v-model="dialog.isShowDialog" width="769px" height="500px" center>
		<div class="dialog-content">
			<div class="subtitle" v-if="props.subtitle">{{ props.subtitle }}</div>
			<div v-html="props.contentHtml"></div>
		</div>
	</el-dialog>
</template>

<script setup lang="ts">
import { reactive, toRefs } from 'vue';
interface Props {
	title: string;
	contentHtml: HtmlType;
	subtitle: string;
}
const props = withDefaults(defineProps<Props>(), {
	title: '',
	contentHtml: '',
	subtitle: '',
});

const state = reactive({
	dialog: {
		isShowDialog: false,
	},
});
const { dialog } = toRefs(state);
const openDialog = () => {
	dialog.value.isShowDialog = true;
};

defineExpose({
	openDialog,
});
</script>

<style lang="scss" scoped>
.subtitle {
	color: var(--el-text-color-secondary);
	text-align: center;
	padding-bottom: 20px;
}
.dialog-content {
	min-height: 200px;
}
</style>
