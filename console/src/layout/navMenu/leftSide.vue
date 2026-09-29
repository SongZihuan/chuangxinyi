<template>
	<div class="menu-container">
		<div class="menu-left">
			<div v-for="val in menuLists" :key="val.name" class="sub-menu" @click="menuClick(val)" @mouseenter="menuActive(val)">
				<div class="menu-left-item" :class="{ active: val.path == activePath }">
					<SvgIcon :name="val.meta.icon" class="svg-class" size="18" />
					<div class="menu-title">{{ val.meta.title }}</div>
				</div>
			</div>
		</div>
		<div class="menu-right" v-if="curentMenu && curentMenu.length > 0">
			<div v-for="val in curentMenu" :key="val.name">
				<div class="menu-left-item" v-if="val.meta" @click="go(val)" :class="{ active: !val.noactive && val.path == route.path }">
					<SvgIcon :name="val.meta.icon" class="svg-class" size="18" />
					<div class="menu-title">{{ val.meta.title }}</div>
				</div>
				<div class="menu-left-item" v-else @click="go(val)" :class="{ active: !val.noactive && val.path == route.path }">
					<SvgIcon :name="val.icon" class="svg-class" size="18" />
					<div class="menu-title">{{ val.title }}</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts" name="navMenuVertical">
import { reactive, computed, onMounted, watch, ref } from 'vue';
import { useRoute, onBeforeRouteUpdate, RouteRecordRaw, useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useThemeConfig } from '/@/stores/themeConfig';

// 定义父组件传过来的值
const props = defineProps({
	// 菜单列表
	menuList: {
		type: Array<RouteRecordRaw>,
		default: () => [],
	},
});
const activePath = ref('');
const curentMenu = ref<any>([]);
const router = useRouter();
const go = (val: any) => {
	router.push({ path: val.path });
	themeConfig.value.isCollapse = false;
};
const menuClick = (val: any) => {
	if (!val.children || val.children.length === 0) {
		router.push({ path: val.path });
		themeConfig.value.isCollapse = false;
	}
};

const menuActive = (val: any) => {
	activePath.value = val.path;
	if (val.children && val.children.length > 0) {
		curentMenu.value = [...val.children];
	} else {
		curentMenu.value = [{ ...val, noactive: true }];
	}
};

// 定义变量内容
const storesThemeConfig = useThemeConfig();
const { themeConfig } = storeToRefs(storesThemeConfig);
const route = useRoute();
const state = reactive({
	// 修复：https://gitee.com/lyt-top/vue-next-admin/issues/I3YX6G
	defaultActive: route.meta.isDynamic ? route.meta.isDynamicPath : route.path,
	isCollapse: false,
});

// 获取父级菜单数据
const menuLists = computed(() => {
	return <RouteItems>props.menuList;
});

// 菜单高亮（详情时，父级高亮）
const setParentHighlight = (currentRoute: RouteToFrom) => {
	const { path, meta } = currentRoute;
	const pathSplit = meta?.isDynamic ? meta.isDynamicPath!.split('/') : path!.split('/');
	if (pathSplit.length >= 4 && meta?.isHide) return pathSplit.splice(0, 3).join('/');
	else return path;
};
// 页面加载时
onMounted(() => {
	state.defaultActive = setParentHighlight(route);
});
// 路由更新时
onBeforeRouteUpdate((to) => {
	// 修复：https://gitee.com/lyt-top/vue-next-admin/issues/I3YX6G
	state.defaultActive = setParentHighlight(to);
	const clientWidth = document.body.clientWidth;
	if (clientWidth < 0) themeConfig.value.isCollapse = false;
});
// 设置菜单的收起/展开
watch(
	() => themeConfig.value.isCollapse,
	(isCollapse) => {
		document.body.clientWidth <= 1000 ? (state.isCollapse = false) : (state.isCollapse = isCollapse);
	},
	{
		immediate: true,
	}
);

curentMenu.value = [...menuLists.value[0].children];

if (curentMenu.value && curentMenu.value.length === 0) {
	curentMenu.value = [{ ...menuLists.value[0], noactive: true }];
}
</script>
<style lang="scss" scoped>
.menu-title {
	margin: 0px;
	padding: 0px 8px;
	height: 32px;
	line-height: 32px;
	font-size: 14px;
	font-weight: 600;
	white-space: nowrap;
	color: #333;
}
.menu-subtitle {
	display: block;
	border-width: 1px;
	border-style: solid;
	border-image: initial;
	border-radius: 2px;
	width: 100%;
	max-width: 100%;
	cursor: pointer;
	font-size: 14px;
	text-align: left;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	color: var(--cb-color-button-tertiary-text, #666);
	background-color: var(--cb-color-button-tertiary-bg, transparent);
	border-color: var(--cb-color-button-tertiary-border, transparent);
	padding: 0px 8px;
	height: 32px;
	line-height: 30px;
	font-size: 12px;
}
.flex-col {
	flex-direction: column;
}
.menu-container {
	display: flex;
	flex-direction: row;
	font-size: 14px;
	.sub-menu {
		width: 100%;
		display: flex;
		flex-direction: column;
	}
	.menu-left {
		width: 120px;
		display: flex;
		flex-direction: column;
		border-right: 1px solid var(--next-border-color-light) !important;
		.menu-left-item {
			display: flex;
			flex-direction: row;
			align-items: center;
			margin-top: 10px;
			cursor: pointer;
			width: 100px;
			padding: 0px 10px;
			border-radius: 4px;
			&:hover {
				box-shadow: 0 2px 12px 0 rgb(0 0 0 / 10%);
				cursor: pointer;
				color: #fff;
				background-color: var(--el-color-primary);
			}
			&:hover .menu-title {
				color: #fff;
				background-color: var(--el-color-primary);
			}
		}
	}
	.active {
		box-shadow: 0 2px 12px 0 rgb(0 0 0 / 10%);
		transition: all 0.3s ease;
		cursor: pointer;
		color: #fff;
		background-color: var(--el-color-primary);
		border-radius: 4px;

		.menu-title {
			color: #fff;
			background-color: var(--el-color-primary);
		}
	}
	.menu-right {
		margin-left: 20px;
		flex: 1;
		width: 140px;
		display: flex;
		flex-direction: column;

		.menu-left-item {
			display: flex;
			flex-direction: row;
			align-items: center;
			padding: 0px 10px;
			margin-top: 10px;
			cursor: pointer;
			&:hover {
				box-shadow: 0 2px 12px 0 rgb(0 0 0 / 10%);
				transition: all 0.3s ease;
				cursor: pointer;
				color: var(--el-color-primary);
			}
		}
	}
}
.svg-class {
	font-size: 20px;
}
</style>
