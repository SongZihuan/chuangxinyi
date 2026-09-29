import { RouteRecordRaw } from 'vue-router';

/**
 * 建议：路由 path 路径与文件夹名称相同，找文件可浏览器地址找，方便定位文件位置
 *
 * 路由meta对象参数说明
 * meta: {
 *      title:          菜单栏及 tagsView 栏、菜单搜索名称（国际化）
 *      isLink：        是否超链接菜单，开启外链条件，`1、isLink: 链接地址不为空 2、isIframe:false`
 *      isHide：        是否隐藏此路由
 *      isKeepAlive：   是否缓存组件状态
 *      isAffix：       是否固定在 tagsView 栏上
 *      isIframe：      是否内嵌窗口，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
 *      roles：         当前路由权限标识，取角色管理。控制路由显示、隐藏。超级管理员：admin 普通角色：common
 *      icon：          菜单、tagsView 图标，阿里：加 `iconfont xxx`，fontawesome：加 `fa xxx`
 * }
 */

// 扩展 RouteMeta 接口
declare module 'vue-router' {
	interface RouteMeta {
		title?: string;
		isLink?: string;
		isHide?: boolean;
		isKeepAlive?: boolean;
		isAffix?: boolean;
		isIframe?: boolean;
		roles?: string[];
		icon?: string;
	}
}

/**
 * 定义动态路由
 * 前端添加路由，请在顶级节点的 `children 数组` 里添加
 * @description 未开启 isRequestRoutes 为 true 时使用（前端控制路由），开启时第一个顶级 children 的路由将被替换成接口请求回来的路由数据
 * @description 各字段请查看 `/@/views/system/menu/main/addMenu.vue 下的 ruleForm`
 * @returns 返回路由菜单数据
 */
export const dynamicRoutes: Array<RouteRecordRaw> = [
	{
		path: '/',
		name: '/',
		component: () => import('/@/layout/index.vue'),
		redirect: '/home',
		meta: {
			isKeepAlive: true,
		},
		children: [],
	},
];

/**
 * 定义404、401界面
 * @link 参考：https://next.router.vuejs.org/zh/guide/essentials/history-mode.html#netlify
 */
export const notFoundAndNoPower = [
	{
		path: '/:path(.*)*',
		name: 'notFound',
		component: () => import('/@/views/error/404.vue'),
		meta: {
			title: '找不到此页面',
			isHide: true,
		},
	},
	{
		path: '/401',
		name: 'noPower',
		component: () => import('/@/views/error/401.vue'),
		meta: {
			title: '没有权限',
			isHide: true,
		},
	},
];

/**
 * 定义静态路由（默认路由）
 * 此路由不要动，前端添加路由的话，请在 `dynamicRoutes 数组` 中添加
 * @description 前端控制直接改 dynamicRoutes 中的路由，后端控制不需要修改，请求接口路由数据时，会覆盖 dynamicRoutes 第一个顶级 children 的内容（全屏，不包含 layout 中的路由出口）
 * @returns 返回路由菜单数据
 */
export const staticRoutes: Array<RouteRecordRaw> = [
	{
		path: '/login',
		name: 'login',
		component: () => import('/@/views/login/index.vue'),
		meta: {
			title: '登录',
		},
	},
	{
		path: '/register',
		name: 'register',
		component: () => import('/@/views/register/index.vue'),
		meta: {
			title: '注册',
		},
	},
	{
		path: '/weixinbindsuccess',
		name: 'weixinbindsuccess',
		component: () => import('/@/views/weixin/bindSuccess.vue'),
		meta: {
			title: '微信绑定',
		},
	},
	{
		path: '/weixinloginsuccess',
		name: 'weixinloginsuccess',
		component: () => import('/@/views/weixin/loginSuccess.vue'),
		meta: {
			title: '微信登录',
		},
	},
	{
		path: '/oauth2',
		name: 'oauth2',
		component: () => import('/@/views/oauth2/index.vue'),
		meta: {
			title: '授权登录',
		},
	},
	{
		path: '/oauth2/open',
		name: 'oauth2open',
		component: () => import('/@/views/oauth2/open.vue'),
		meta: {
			title: '功能开通',
		},
	},
	{
		path: '/product',
		name: 'product',
		component: () => import('/@/views/product/index.vue'),
		meta: {
			title: '套餐购买',
		},
	},
	{
		path: '/oauth2/logout',
		name: 'oauth2logout',
		component: () => import('/@/views/oauth2/logout.vue'),
		meta: {
			title: '退出登录',
		},
	},
	{
		path: '/forgetPassword',
		name: 'forgetPassword',
		component: () => import('/@/views/forgetPassword/index.vue'),
		meta: {
			title: '忘记密码',
		},
	},
	{
		path: '/externalLinkSkip',
		name: 'externalLinkSkip',
		component: () => import('/@/views/externalLinkSkip/index.vue'),
		meta: {
			title: '外链跳转',
		},
	},
	{
		path: '/paysuccess',
		name: 'paysuccess',
		component: () => import('/@/views/success/paySuccess.vue'),
		meta: {
			title: '支付成功',
		},
	},
	{
		path: '/facesuccess',
		name: 'facesuccess',
		component: () => import('/@/views/success/faceSuccess.vue'),
		meta: {
			title: '人脸识别成功',
		},
	},
	{
		path: '/fuwuhao',
		name: 'fuwuhao',
		component: () => import('/@/views/fuwuhao/index.vue'),
		meta: {
			title: '服务号登录',
		},
	},
	{
		path: '/pay',
		name: 'pay',
		component: () => import('/@/views/pay/index.vue'),
		meta: {
			title: '支付',
		},
	},
	{
		path: '/visitor',
		name: 'visitor',
		component: () => import('/@/views/userCenter/homePage/visiable.vue'),
		meta: {
			title: '访客',
		},
	},
	{
		path: '/purchase',
		name: 'purchase',
		component: () => import('/@/views/purchase/index.vue'),
		meta: {
			title: '订单充值',
		},
	},
	{
		path: '/404',
		name: '404',
		component: () => import('/@/views/error/404.vue'),
		meta: {
			title: '页面未找到',
		},
	},
];
