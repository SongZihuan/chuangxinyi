interface params extends pageTypes {
	domain?: string;
}
//站点列表类型
export interface websiteListTypes {
	tableData: {
		data: any[];
		total: number;
		loading: boolean;
		param: params;
	};
}

export interface websiteFromType {
	name: string;
	describe?: string;
	domain: string;
	agreement: string;
	keyMap: dictTypes[];
	policy: any[];
	pubkey: string;
	status: number;
}
