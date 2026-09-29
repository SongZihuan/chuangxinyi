export interface invoiceUserTypes {
    page: number,
    pagesize: number,
    starttime?:number,
    endtime?:number
    timetype?:number
}
/**
 * Amount int64 `json:"amount"`
 * Type int64 `json:"type"` // 1-个人普票 2-企业普票 3-企业专票
 */
export interface invoiceTypes{
    amount:number,
    type:number
}