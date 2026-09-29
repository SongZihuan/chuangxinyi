/**
 * TokenType int64        `json:"tokenType"`
 *     Token     string       `json:"token"`
 *     Type      int64        `json:"type"`
 *     Data      string       `json:"data"`
 *     CreateAt int64 `json:"createAt"`
 */
export interface tokenDataType {
    tokenType: number;
    token: string;
    type: number;
    data: string;
    createAt: number;
}

export interface tokenTypes {
    page: number;
    pagesize: number;
    starttime?: number;
    endtime?: number;
    timetype?: number | string;
    range: Date[];
    token?: string;
}
export type tableType = {
    data: tokenDataType[];
    total: number;
    loading: boolean;
    param: tokenTypes;
};
export interface tokenStateType {
    tableData: tableType;
}