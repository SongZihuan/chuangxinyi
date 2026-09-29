import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';
import pinia from '/@/stores/index';
import { storeToRefs } from 'pinia';
import { useKeepALiveNames } from '/@/stores/keepAliveNames';
import { useRoutesList } from '/@/stores/routesList';
import { useThemeConfig } from '/@/stores/themeConfig';
import { Local, Session } from '/@/utils/storage';
import { staticRoutes, notFoundAndNoPower, dynamicRoutes } from '/@/router/route';
import { initBackEndControlRoutes } from '/@/router/backEnd';
import { useSocket } from '/@/hooks/useWebstocket';
import { ElMessage } from 'element-plus';
import { NextLoading } from '/@/utils/loading';
/**
 * 1、前端控制路由时：isRequestRoutes 为 false，需要写 roles，需要走 setFilterRoute 方法。
 * 2、后端控制路由时：isRequestRoutes 为 true，不需要写 roles，不需要走 setFilterRoute 方法），
 * 相关方法已拆解到对应的 `backEnd.ts` 与 `frontEnd.ts`（他们互不影响，不需要同时改 2 个文件）。
 * 特别说明：
 * 1、前端控制：路由菜单由前端去写（无菜单管理界面，有角色管理界面），角色管理中有 roles 属性，需返回到 userInfo 中。
 * 2、后端控制：路由菜单由后端返回（有菜单管理界面、有角色管理界面）
 */

// 读取 `/src/stores/themeConfig.ts` 是否开启后端控制路由配置
const storesThemeConfig = useThemeConfig(pinia);
const { themeConfig } = storeToRefs(storesThemeConfig);
/**
 * 创建一个可以被 Vue 应用程序使用的路由实例
 * @method createRouter(options: RouterOptions): Router
 * @link 参考：https://next.router.vuejs.org/zh/api/#createrouter
 */
export const router = createRouter({
	history: createWebHistory(import.meta.env.VITE_PUBLIC_PATH),
	/**
	 * 说明：
	 * 1、notFoundAndNoPower 默认添加 404、401 界面，防止一直提示 No match found for location with path 'xxx'
	 * 2、backEnd.ts(后端控制路由)、frontEnd.ts(前端控制路由) 中也需要加 notFoundAndNoPower 404、401 界面。
	 *    防止 404、401 不在 layout 布局中，不设置的话，404、401 界面将全屏显示
	 */
	routes: [...notFoundAndNoPower, ...staticRoutes],
});

export const resetMatcher = () => {
	let tmp = createRouter({
		history: createWebHistory(import.meta.env.VITE_PUBLIC_PATH),
		/**
		 * 说明：
		 * 1、notFoundAndNoPower 默认添加 404、401 界面，防止一直提示 No match found for location with path 'xxx'
		 * 2、backEnd.ts(后端控制路由)、frontEnd.ts(前端控制路由) 中也需要加 notFoundAndNoPower 404、401 界面。
		 *    防止 404、401 不在 layout 布局中，不设置的话，404、401 界面将全屏显示
		 */
		routes: [...notFoundAndNoPower, ...staticRoutes],
	}) as any;
	(router as any).matcher = tmp.matcher;
};

/**
 * 路由多级嵌套数组处理成一维数组
 * @param arr 传入路由菜单数据数组
 * @returns 返回处理后的一维路由菜单数组
 */
export function formatFlatteningRoutes(arr: Array<RouteRecordRaw>): Array<RouteRecordRaw> {
	if (arr.length <= 0) return [];
	for (let i = 0; i < arr.length; i++) {
		if (arr[i].children) {
			arr = arr.slice(0, i + 1).concat(arr[i].children as Array<RouteRecordRaw>, arr.slice(i + 1));
		}
	}
	return arr;
}

/**
 * 一维数组处理成多级嵌套数组（只保留二级：也就是二级以上全部处理成只有二级，keep-alive 支持二级缓存）
 * @description isKeepAlive 处理 `name` 值，进行缓存。顶级关闭，全部不缓存
 * @link 参考：https://v3.cn.vuejs.org/api/built-in-components.html#keep-alive
 * @param arr 处理后的一维路由菜单数组
 * @returns 返回将一维数组重新处理成 `定义动态路由（dynamicRoutes）` 的格式
 */
export function formatTwoStageRoutes(arr: Array<RouteRecordRaw>): Array<RouteRecordRaw> {
	if (arr.length <= 0) return [];

	const newArr: any[] = [];
	const cacheList: Array<string> = [];

	arr.forEach((v: any) => {
		if (v.path === '/') {
			newArr.push({ component: v.component, name: v.name, path: v.path, redirect: v.redirect, meta: v.meta, children: [] });
		}
	});

	if (newArr.length !== 1) {
		return [];
	}

	arr.forEach((v: any) => {
		if (v.path !== '/') {
			// 判断是否是动态路由（xx/:id/:name），用于 tagsView 等中使用
			// 修复：https://gitee.com/lyt-top/vue-next-admin/issues/I3YX6G
			if (v.path.indexOf('/:') > -1) {
				v.meta['isDynamic'] = true;
				v.meta['isDynamicPath'] = v.path;
			}
			newArr[0].children.push({ ...v });
			// 存 name 值，keep-alive 中 include 使用，实现路由的缓存
			// 路径：/@/layout/routerView/parent.vue
			if (newArr[0].meta.isKeepAlive && v.meta.isKeepAlive) {
				cacheList.push(v.name);
				const stores = useKeepALiveNames(pinia);
				stores.setCacheKeepAlive(cacheList);
			}
		}
	});

	return newArr;
}
//路由白名单
const onlyNotTokenList = ['/register', '/login', '/fuwuhao', '/forgetPassword'];

const tokenWhiteList = [
	'/weixinloginsuccess',
	'/weixinbindsuccess',
	'/oauth2',
	'/oauth2/logout',
	'/externalLinkSkip',
	'/paysuccess',
	'/facesuccess',
	'/visitor',
	'/pay',
	'product',
	'/purchase',
];

// 路由加载前
router.beforeEach(async (to, from, next) => {
	NProgress.configure({ showSpinner: false });
	if (to.meta.title) NProgress.start();
	const token = Session.get('token');
	if (token) {
		useSocket();
	}

	if (tokenWhiteList.includes(to.path) && !onlyNotTokenList.includes(to.path)) {
		next();
	} else if (!token) {
		Session.clear();
		Local.clear();
		if (onlyNotTokenList.includes(to.path)) {
			next();
		} else {
			next(`/login?redirect=${to.path}&params=${JSON.stringify(to.query ? to.query : to.params)}`);
		}
	} else {
		if (onlyNotTokenList.includes(to.path)) {
			if (Number(to.query.forceLogin) === 1) {
				Session.clear();
				Local.clear();
				next();
			} else {
				next('/home');
			}
		} else {
			const storesRoutesList = useRoutesList(pinia);
			const { routesList } = storeToRefs(storesRoutesList);
			if (routesList.value.length === 0) {
				let res = await initBackEndControlRoutes();
				if (!res) {
					ElMessage.error('没有登录权限');
					Session.clear();
					Local.clear();
					next(`/login?redirect=${to.path}&params=${JSON.stringify(to.query ? to.query : to.params)}`);
				} else {
					// 解决刷新时，一直跳 404 页面问题，关联问题 No match found for location with path 'xxx'
					// to.query 防止页面刷新时，普通路由带参数时，参数丢失。动态路由（xxx/:id/:name"）isDynamic 无需处理
					next({ path: to.path, query: to.query });
				}
			} else {
				next();
			}
		}
	}
	NProgress.done();
});

// 路由加载后
router.afterEach(() => {
	NProgress.done();
	NextLoading.done();
});

// 导出路由
export default router;
