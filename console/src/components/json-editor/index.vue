<script setup lang="ts" name="JsonEditor">
import { reactive, watch } from 'vue';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
const props = defineProps({
	jsonData: String,
});
const state = reactive({
	val: JSON.stringify(props.jsonData),
	data: props.jsonData,
	showLine: true,
	showLineNumber: true,
	showDoubleQuotes: true,
	showLength: true,
	editable: false,
	showIcon: true,
	editableTrigger: 'click',
	deep: 3,
});

watch(
	() => state.val,
	(newVal) => {
		try {
			state.data = JSON.parse(newVal);
		} catch (err) {
			// console.log('JSON ERROR');
		}
	}
);

watch(
	() => state.data,
	(newVal) => {
		try {
			state.val = JSON.stringify(newVal);
		} catch (err) {
			// console.log('JSON ERROR');
		}
	}
);
</script>

<template>
	<el-card shadow="never">
		<template #header>
			<div class="card-header">json数据</div>
		</template>
		<vue-json-pretty
			v-model:data="state.data"
			:deep="state.deep"
			:show-double-quotes="state.showDoubleQuotes"
			:show-line="state.showLine"
			:show-length="state.showLength"
			:show-icon="state.showIcon"
			:show-line-number="state.showLineNumber"
			:editable="state.editable"
			:editable-trigger="(state.editableTrigger as any)"
		/>
	</el-card>
</template>
