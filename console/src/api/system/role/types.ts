export interface roleTypes {
	name: string; //角色名称
	sign: string; //角色标识
	describe?: string; //角色描述
	status: number; //角色状态
	permission: any; //菜单权限数组
	notDelete: boolean; //不可删除
	notChangeSign: boolean; //不可编辑sign
	notChangePermissions: boolean; //不可编辑权限
	roleID?: number;
	id?: number;
	belong?: string;
}

export interface roleDelTypes {
	roleId: number;
}

export interface RowRoleType {
	name: string;
	sign: string;
	describe: string;
	sort: number;
	status: boolean;
	createTime: string;
}

interface SysRoleTableType extends TableType {
	data: RowRoleType[];
	param: {
		name: string;
		page: number;
		pagesize: number;
		total?: number;
	};
}

export interface SysRoleState {
	tableData: SysRoleTableType;
}

export interface roleMoveTypes {
	id: number;
	isUp: boolean;
}
