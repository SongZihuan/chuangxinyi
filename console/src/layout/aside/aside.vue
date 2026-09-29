<template>
	<div class="h100" v-show="!isTagsViewCurrenFull">
		<el-drawer
			v-model="themeConfig.isCollapse"
			direction="ltr"
			class="mt50 drawer"
			size="375px"
			lock-scroll
			open-delay="300"
			title="菜 单"
			close-on-click-modal
			close-on-press-escape
			modal-class="my-voerlay"
		>
			<div class="aside">
				<div class="left-aside">
					<leftSide :menuList="state.menuList" />
				</div>
				<!-- <Vertical :menuList="state.menuList" class="flex1" /> -->
			</div>
		</el-drawer>
	</div>
</template>

<script setup lang="ts" name="layoutAside">
import { defineAsyncComponent, reactive, watch, onBeforeMount, ref, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useRoutesList } from '/@/stores/routesList';
import { useThemeConfig } from '/@/stores/themeConfig';
import { useTagsViewRoutes } from '/@/stores/tagsViewRoutes';

import mittBus from '/@/utils/mitt';
import {Logger} from "sass";

// const Vertical = defineAsyncComponent(() => import('/@/layout/navMenu/vertical.vue'));
const leftSide = defineAsyncComponent(() => import('/@/layout/navMenu/leftSide.vue'));

// 定义变量内容
const layoutAsideScrollbarRef = ref();
const stores = useRoutesList();
const storesThemeConfig = useThemeConfig();
const storesTagsViewRoutes = useTagsViewRoutes();
const { routesList } = storeToRefs(stores);
const { themeConfig } = storeToRefs(storesThemeConfig);
const { isTagsViewCurrenFull } = storeToRefs(storesTagsViewRoutes);
const state = reactive<AsideState>({
	menuList: [],
	clientWidth: 0,
});

onMounted(() => {
	themeConfig.value.isCollapse = false;
});


// 关闭移动端蒙版
const closeLayoutAsideMobileMode = () => {
  const el = document.querySelector('.layout-aside-mobile-mode');
  el?.setAttribute('style', 'animation: error-img-two 0.3s');
  setTimeout(() => {
    el?.parentNode?.removeChild(el);
  }, 300);
  const clientWidth = document.body.clientWidth;
  if (clientWidth < 1000) themeConfig.value.isCollapse = false;
  document.body.setAttribute('class', '');
};
// 设置/过滤路由（非静态路由/是否显示在菜单中）
const setFilterRoutes = () => {
  state.menuList = filterRoutesFun(routesList.value);
};
// 路由过滤递归函数
const filterRoutesFun = <T extends RouteItem>(arr: T[]): T[] => {
  let newArr = arr
      .filter((item: T) => !item.meta?.isHide)
      .map((item: T) => {
        item = Object.assign({}, item);
        if (item.children) item.children = filterRoutesFun(item.children);
        return item;
      });

  return newArr
};
// 设置菜单导航是否固定（移动端）
const initMenuFixed = (clientWidth: number) => {
  state.clientWidth = clientWidth;
};

// 页面加载前
onBeforeMount(() => {
  initMenuFixed(document.body.clientWidth);
  setFilterRoutes();
  // 此界面不需要取消监听(mittBus.off('setSendColumnsChildren))
  // 因为切换布局时有的监听需要使用，取消了监听，某些操作将不生效
  mittBus.on('setSendColumnsChildren', (res: MittMenu) => {
    let newData: any = [];
    res.children.forEach((item: any) => {
      if (item.children.length == 1) {
        newData.push(item.children[0])
      } else {
        newData.push(item)
      }
    })
    state.menuList = newData;
  });
  // 开启经典布局分割菜单时，设置菜单数据
  mittBus.on('setSendClassicChildren', () => {
    // 开启经典布局分割菜单时，重新处理菜单数据
    mittBus.on('getBreadcrumbIndexSetFilterRoutes', () => {
      setFilterRoutes();
    });
    // 监听窗口大小改变时(适配移动端)
    mittBus.on('layoutMobileResize', (res: LayoutMobileResize) => {
      initMenuFixed(res.clientWidth);
      closeLayoutAsideMobileMode();
    });
  })
});
// 监听 pinia 值的变化，动态赋值给菜单中
watch(
    () => [themeConfig.value.isShowLogoChange, themeConfig.value.isShowLogo, themeConfig.value.isClassicSplitMenu],
    ([isShowLogoChange, isShowLogo]) => {
      if (isShowLogoChange !== isShowLogo) {
        if (layoutAsideScrollbarRef.value) layoutAsideScrollbarRef.value.update();
      }
      setFilterRoutes();
    },
    {
      deep: true,
    }
);
// 监听用户权限切换，用于演示 `权限管理 -> 前端控制 -> 页面权限` 权限切换不生效
watch(
    () => routesList.value,
    () => {
      setFilterRoutes();
    }
);

</script>
<style>
.drawer {
	margin-top: 50px;
}
.aside {
	display: flex;
	flex-direction: row;
	height: 100%;
	.left-aside {
		width: 100%;
		height: 100%;
		padding-right: 10px;
		margin-left: 15px;
	}
	.flex1 {
		flex: 1;
		margin-left: 10px;
	}
}
.my-voerlay {
	z-index: 80 !important;
}
</style>
