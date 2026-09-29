export interface bannedTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number | string;
	from: string;
}

/**
 * WebID int64 `json:"webID"`
 *     AllowLogin bool `json:"allowLogin"`
 *     AllowDefray bool `json:"allowDefray"`
 *     AllowMsg bool `json:"allowMsg"`
 */
export interface bannedUpdateTypes {
	webID: number;
	allowLogin: boolean;
	allowDefray: boolean;
	allowMsg: boolean;
}
