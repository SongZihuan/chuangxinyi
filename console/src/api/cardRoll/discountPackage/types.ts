type quota = {
	type: number;
	bottom: number;
	send?: number;
	discount?: number;
	pre?: number;
	amount?: number;
};

export interface formTypes {
	name: string;
	describe: string;
	shortDescribe: string;
	type: number;
	quota: quota;
	dayLimit: number;
	monthLimit: number;
	yearLimit: number;
	limit: number;
	needVerify: boolean;
	needCompany: boolean;
	needUserOrigin: boolean;
	needCompanyOrigin: boolean;
	show: boolean;
	needUserFace: boolean;
	needCompanyFace: boolean;
}
export interface formTypes {
	name: string;
	describe: string;
	shortDescribe: string;
	type: number;
	quota: quota;
	dayLimit: number;
	monthLimit: number;
	yearLimit: number;
	limit: number;
	needVerify: boolean;
	needCompany: boolean;
	needUserOrigin: boolean;
	needCompanyOrigin: boolean;
	show: boolean;
	needUserFace: boolean;
	needCompanyFace: boolean;
}

export interface stateTypes {
	tableData: {
		data: any[];
		total: number;
		loading: boolean;
		param: pageTypes;
	};
}

export interface allocationUserType {
	uid?: number | null;
	id?: number | null;
	discountID: number | null;
}

export interface adminUserType {
	id?: number | string | null;
	uid?: number | string | null;
	page: number;
	pagesize: number;
	starttime?: number | null;
	endtime?: number | null;
	timetype?: number;
	range: Date[];
}
