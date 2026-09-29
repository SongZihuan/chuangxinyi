export interface agreementStateType {
	tableData: {
		data: any[];
		total: number;
		loading: boolean;
		param: pageTypes;
	};
}

export interface fromType {
	aid: string;
	content: HtmlType | any;
}
