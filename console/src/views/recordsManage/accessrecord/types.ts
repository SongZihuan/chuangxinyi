/**
 * ID                int64          `json:"id"`
 *     RequestsID string `json:"requestsID"`
 *     UserID            int64  `json:"userID"`
 *     UserUid           string `json:"userUid"`
 *     RoleID            int64  `json:"roleID"`
 *     RoleName          string `json:"roleName"`
 *     RoleSign          string `json:"roleSign"`
 *     WebId             int64  `json:"webID"`
 *     WebName           string `json:"webName"`
 *     WebDomain         string `json:"webDomain"`
 *     RequestsWebId     int64  `json:"requestsWebID"`
 *     RequestsWebName   string `json:"requestsWebName"`
 *     RequestsWebDomain string `json:"requestsWebDomain"`
 *     Ip1               string         `json:"ip1"`
 *     Ip2               string         `json:"ip2"`
 *     Port1             string         `json:"port1"`
 *     Port2             string         `json:"port2"`
 *     GeoCode           string         `json:"geo_code"`
 *     Geo               string         `json:"geo"`
 *     Scheme            string         `json:"scheme"`
 *     Method            string         `json:"method"`
 *     Host              string         `json:"host"`
 *     Path              string         `json:"path"`
 *     Query             string         `json:"query"`
 *     RequestsBody      string         `json:"requestsBody"`
 *     ResponseBody      string `json:"responseBody"`
 *     ResponseBodyError string `json:"responseBodyError"`
 *     RequestsHeader    string         `json:"requestsJeader"`
 *     ResponseHeader    string `json:"responseHeader"`
 *     StatusCode        int64  `json:"statusCode"`
 *     PanicError        string `json:"panicError"`
 *     Message           string `json:"message"`
 *     UseTimeMs           int64  `json:"useTimeMs"`
 *     CreateAt          int64      `json:"createAt"`
 *     StartAt           int64   `json:"startAt"`
 *     EndAt             int64   `json:"endAt"`
 */
export interface accessrecordDataType {
    id: number;
    requestsID: string;
    userID: number;
    userUid: string;
    roleID: number;
    roleName: string;
    roleSign: string;
    webId: number;
    webName: string;
    webDomain: string;
    requestsWebId: number;
    requestsWebName: string;
    requestsWebDomain: string;
    ip1: string;
    ip2: string;
    port1: string;
    port2: string;
    geoCode: string;
    geo: string;
    scheme: string;
    method: string;
    host: string;
    path: string;
    query: string;
    requestsBody: string;
    responseBody: string;
    responseBodyError: string;
    requestsHeader: string;
    responseHeader: string;
    statusCode: number;
    panicError: string;
    message: string;
    useTimeMs: number;
    createAt: number;
    startAt: number;
    endAt: number;
}

export interface accessrecordTypes {
    page: number;
    pagesize: number;
    starttime?: number;
    endtime?: number;
    timetype?: number | string;
    range: Date[];
}
export type tableType = {
    data: accessrecordDataType[];
    total: number;
    loading: boolean;
    param: accessrecordTypes;
};
export interface accessrecordStateType {
    tableData: tableType;
}