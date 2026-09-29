<template>
	<el-config-provider :size="getGlobalComponentSize" :locale="zhCn">
		<router-view v-show="setLockScreen" />
		<CloseFull v-if="!themeConfig.isLockScreen" />
		<Upgrade v-if="checkVersion" />
	</el-config-provider>
  <Message></Message>
</template>

<script setup lang="ts" name="app">
import { defineAsyncComponent, computed, ref, onBeforeMount, onMounted, onUnmounted, nextTick, watch, onBeforeUnmount } from 'vue';
import { useRoute } from 'vue-router';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import { storeToRefs } from 'pinia';
import { useTagsViewRoutes } from '/@/stores/tagsViewRoutes';
import { useThemeConfig } from '/@/stores/themeConfig';
import other from '/@/utils/other';
import { Local, Session } from '/@/utils/storage';
import mittBus from '/@/utils/mitt';
import setIntroduction from '/@/utils/setIconfont';
import Message from '/@/components/message/message.vue'

// 引入组件
const CloseFull = defineAsyncComponent(() => import('/@/layout/navBars/topBar/closeFull.vue'));
const Upgrade = defineAsyncComponent(() => import('/@/layout/upgrade/index.vue'));
// 定义变量内容
const setingsRef = ref();
const route = useRoute();
const stores = useTagsViewRoutes();
const storesThemeConfig = useThemeConfig();
const { themeConfig } = storeToRefs(storesThemeConfig);
import { major, valid } from "semver";

const setLockScreen = computed(() => {
  // 防止锁屏后，刷新出现不相关界面
  // https://gitee.com/lyt-top/vue-next-admin/issues/I6AF8P
  return themeConfig.value.isLockScreen ? themeConfig.value.lockScreenTime > 1 : themeConfig.value.lockScreenTime >= 0;
});

// 获取版本号
const checkVersion = computed((): boolean => {
  if (route.path === '/login' || route.path === '/fuwuhao') {
    return false
  }

  let localVersion = Local.get('version')

  if (!localVersion) {
    return false
  }

  if (!valid(localVersion)) {
    return true
  }

  // @ts-ignore
  if (!valid(__NEXT_VERSION__)) {
    return false
  }

  // @ts-ignore
  return major(localVersion) < major(__NEXT_VERSION__);
});

// 获取全局组件大小
const getGlobalComponentSize = computed(() => {
	return other.globalComponentSize();
});
// 设置初始化，防止刷新时恢复默认
onBeforeMount(() => {
	// 设置批量第三方 icon 图标
	setIntroduction.cssCdn();
	// 设置批量第三方 js
	setIntroduction.jsCdn();
});
onBeforeUnmount(() => {
	Local.remove('redirect_uri');
	Local.remove('back_domain');
	Local.remove('window_url');
});
// 页面加载时
onMounted(async () => {
	nextTick(() => {
		// 监听布局配'置弹窗点击打开
		mittBus.on('openSetingsDrawer', () => {
			setingsRef.value.openDrawer();
		});
		// 获取缓存中的布局配置
		if (Local.get('themeConfig')) {
			storesThemeConfig.setThemeConfig({ themeConfig: Local.get('themeConfig') });
			document.documentElement.style.cssText = Local.get('themeConfigStyle');
		}
		// 获取缓存中的全屏配置
		if (Session.get('isTagsViewCurrenFull')) {
			stores.setCurrenFullscreen(Session.get('isTagsViewCurrenFull'));
		}
	});
});
// 页面销毁时，关闭监听布局配置/i18n监听
onUnmounted(() => {
	mittBus.off('openSetingsDrawer', () => {});
});
// 监听路由的变化，设置网站标题
watch(
	() => route.path,
	() => {
		other.useTitle();
	},
	{
		deep: true,
	}
);
</script>
