
/**
 *  WebId     int64        `json:"webID"`  // 要返回给用户
 *     WebName   string       `json:"webName"`
 *     WebDomain string       `json:"webDomain"`
 *     Ip        string       `json:"ip"`
 *     Geo       string       `json:"geo"`
 *     GeoCode   string       `json:"geoCode"`
 *     LoginTime int64        `json:"login_time"`
 */
export interface authorizationRecordDataType {
    webID: number;
    webName: string;
    webDomain: string;
    ip: string;
    geo: string;
    geoCode: string;
    login_time: number;
}

export interface authorizationRecordTypes {
    page: number;
    pagesize: number;
    starttime?: number;
    endtime?: number;
    timetype?: number | string;
    range: Date[];
    from: string;
    uid?: string;
}
export type tableType = {
    data: authorizationRecordDataType[];
    total: number;
    loading: boolean;
    param: authorizationRecordTypes;
};
export interface authorizationRecordStateType {
    tableData: tableType;
}
