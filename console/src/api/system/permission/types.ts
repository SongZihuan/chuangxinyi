export interface permissionTypes {
	permissionSuperior?: any; // 上级菜单
	parentID?: number | any; //父级菜单ID
	permissionType: number; // 菜单类型
	permissionId?: number;
	name: string; // 路由名称
	permissionCategory: number; //菜单类别
	component: string; // 组件路径
	componentAlias: string; // 组件路径别名
	isLink: boolean; // 是否外链
	path: string; // 路由路径
	redirect?: string; // 路由重定向，有子集 children 时
	title: string; // 菜单名称
	icon: string; // 菜单图标
	isHide: boolean; // 是否隐藏
	isKeepAlive: boolean; // 是否缓存
	isAffix: boolean; // 是否固定
	metaIsLink?: string; // 外链/内嵌时链接地址（http:xxx.com），开启外链条件，`1、isLink: 链接地址不为空`
	isIframe: boolean; // 是否内嵌，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
	roles: number[]; // 权限标识，取角色管理
	btnPower: string; // 菜单类型为按钮时，权限标识
}

export interface searchTypes {
	name?: string;
	page: number;
	pagesize: number;
	total: number;
}
export interface permissionDeltypes {
	id: number;
}
export interface permissionMoveTypes {
	id: number;
	isUp: boolean;
}

export interface roleActionTypes {
	permissionID: number;
	roleID: number;
}
