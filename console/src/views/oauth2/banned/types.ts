/**
 *  WebId     int64        `json:"webID"`  // 要返回给用户
 *     WebName   string       `json:"webName"`
 *     WebDomain string       `json:"webDomain"`
 *     AllowLogin  bool         `json:"allowLogin"`
 *     AllowDefray bool         `json:"allowDefray"`
 *     AllowMsg    bool         `json:"allowMsg"`
 */
export interface bannedDataType {
	webID: number;
	webName: string;
	webDomain: string;
	allowLogin: boolean;
	allowDefray: boolean;
	allowMsg: boolean;
}
export interface bannedUpdateTypes {
	id?: number;
	webID: number;
	allowLogin: boolean;
	allowDefray: boolean;
	allowMsg: boolean;
}
export interface bannedTypes {
	limit: number;
	uid?: number;
}
export type tableType = {
	data: bannedDataType[];
	total: number;
	loading: boolean;
	param: bannedTypes;
};
export interface bannedStateType {
	tableData: tableType;
}
