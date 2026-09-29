/**
 * OrderID string `json:"orderID"`
 *     Title string `json:"title"`
 *     From string `json:"from"`
 *     CreateAt int64 `json:"createAt"`
 *     LastReplyAt int64 `json:"lastReplyAt"`
 *     FinishAt int64 `json:"finishAt"`
 */
export interface IWorkOrderUserList {
    orderID: string;
    title: string;
    from: string;
    createAt: number;
    lastReplyAt: number;
    finishAt: number;
}
/**
 *  ID int64 `form:"id,optional"`
 *  UID string `form:"uid,optional"`
 *  page: 1,
 *            pagesize: 20,
 *            starttime: 0,
 *            endtime: 0,
 *            range: [],
 *            timetype: '',
 */
export interface IWorkOrderUserListParams {
    id?: number;
    uid?: string;
    page: number;
    pagesize: number;
    starttime: number;
    endtime: number;
    range: number[];
    timetype: string;
}
export interface IWorkOrderUserListRes {
    data: IWorkOrderUserList[];
    total: number;
    loading: boolean;
    param: IWorkOrderUserListParams;
}
export interface IWorkOrderUserListResData {
    tableData: IWorkOrderUserListRes;
}
