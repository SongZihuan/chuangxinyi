/**
 * Fid string `json:"fid"`
 * Type int64 `json:"type"`
 * MediaType string `json:"mediaType"`
 */
export interface fileTypes {
	fid: string;
	type: number;
	mediaType: string;
}
export interface fileParamTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
	range: Date[];
}
interface tableTypes {
	data: fileTypes[];
	loading: boolean;
	total: number;
	param: fileParamTypes;
}
export interface fileStatsTypes {
	tableData: tableTypes;
}
