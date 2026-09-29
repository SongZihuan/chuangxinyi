export interface listParamType {
	page: number;
	pagesize: number;
	timetype: string | number;
	range: any[];
	starttime?: number | null;
	endtime?: number | null;
	uid?: string;
}

export interface fromType {
	aid: string;
	content: HtmlType | any;
}
