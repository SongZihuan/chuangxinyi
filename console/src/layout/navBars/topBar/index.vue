<template>
	<div class="layout-navbars-breadcrumb-index" :class="{shadow: showShadow}">
		<Breadcrumb />
		<User />
	</div>
</template>

<script setup lang="ts" name="layoutBreadcrumbIndex">
import {defineAsyncComponent, computed, reactive, onMounted, onUnmounted, watch, ref} from 'vue';
import { storeToRefs } from 'pinia';
import { useRoutesList } from '/@/stores/routesList';
import { useThemeConfig } from '/@/stores/themeConfig';
import mittBus from '/@/utils/mitt';

// 引入组件
const Breadcrumb = defineAsyncComponent(() => import('/@/layout/navBars/topBar/breadcrumb.vue'));
const User = defineAsyncComponent(() => import('/@/layout/navBars/topBar/user.vue'));

// 定义变量内容
const stores = useRoutesList();
const storesThemeConfig = useThemeConfig();
const { themeConfig } = storeToRefs(storesThemeConfig);
const { routesList } = storeToRefs(stores);
const state = reactive({
	menuList: [] as RouteItems,
});

const showShadow = ref(false)

watch(themeConfig.value, () => {
  if (themeConfig.value.isCollapse) {
    setTimeout(()=>{
      showShadow.value = true
    }, 500)
  } else {
    showShadow.value = false
  }
}, {deep: true});

// 设置 logo 显示/隐藏
const setIsShowLogo = computed(() => {
	let { isShowLogo } = themeConfig.value;
	return isShowLogo
});

// 设置是否显示横向导航菜单
// 设置/过滤路由（非静态路由/是否显示在菜单中）
const setFilterRoutes = () => {
  state.menuList = filterRoutesFun(routesList.value);
};
// 设置了分割菜单时，删除底下 children
const delClassicChildren = <T extends ChilType>(arr: T[]): T[] => {
	arr.map((v: T) => {
		if (v.children) delete v.children;
	});
	return arr;
};
// 路由过滤递归函数
const filterRoutesFun = <T extends RouteItem>(arr: T[]): T[] => {
	return arr
		.filter((item: T) => !item.meta?.isHide)
		.map((item: T) => {
			item = Object.assign({}, item);
			if (item.children) item.children = filterRoutesFun(item.children);
			return item;
		});
};
// 传送当前子级数据到菜单中
const setSendClassicChildren = (path: string) => {
	const currentPathSplit = path.split('/');
	let currentData: MittMenu = { children: [] };
	filterRoutesFun(routesList.value).map((v: RouteItem, k: number) => {
		if (v.path === `/${currentPathSplit[1]}`) {
			v['k'] = k;
			currentData['item'] = { ...v };
			currentData['children'] = [{ ...v }];
			if (v.children) currentData['children'] = v.children;
		}
	});
	return currentData;
};
// 页面加载时
onMounted(() => {
	setFilterRoutes();
	mittBus.on('getBreadcrumbIndexSetFilterRoutes', () => {
		setFilterRoutes();
	});
});
// 页面卸载时
onUnmounted(() => {
	mittBus.off('getBreadcrumbIndexSetFilterRoutes', () => {});
});
</script>

<style scoped lang="scss">
.layout-navbars-breadcrumb-index {
	height: 50px;
	display: flex;
	align-items: center;
	background: var(--next-bg-topBar);
	border-bottom: 1px solid var(--next-border-color-light);
	position: relative;
	z-index: 1000;
}
</style>
