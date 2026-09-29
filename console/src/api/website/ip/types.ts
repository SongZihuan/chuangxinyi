export interface params extends pageTypes {
	websiteid?: string;
}
export interface ipListTypes {
	tableData: {
		data: any[];
		total: number;
		loading: boolean;
		param: params;
	};
}

export interface createType {
	websiteID?: number | string;
	ip: string;
	domain?: any;
}
