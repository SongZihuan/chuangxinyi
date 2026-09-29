export interface rechargeStateTypes {
	tableData: {
		data: any[];
		total: number;
		loading: boolean;
		param: {
			page: number;
			pagesize: number;
			timetype: string | number;
			range: any[];
			starttime?: number | null;
			endtime?: number | null;
			uid?: string;
		};
	};
}

export interface formTypes {
	cny: number;
	payType: number;
	couponsID?: string | null;
	payMode?: number;
}
