/**
 * Content string `json:"content"`
 *     From string `json:"from"`
 *     CreateAt int64 `json:"createAt"`
 */
export interface auditDataType {
	content: string;
	from: string;
	createAt: number;
}

export interface auditTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number | string;
	range: Date[];
	src: string;
	uid?: string;
	fromID?: number,
}
export type tableType = {
	data: auditDataType[];
	total: number;
	loading: boolean;
	param: auditTypes;
};
export interface auditStateType {
	tableData: tableType;
}

export interface invoiceAdminState {
	tableData: auditStateType;
}
