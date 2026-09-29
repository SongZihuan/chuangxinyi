export interface invoiceAdminTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
}

/**
 * type AdminProcessInvoiceReq {
 *     ID string `json:"id"`
 *     Status int64 `json:"status"`
 * }
 */
export interface adminProcessInvoiceTypes {
	orderID: string;
	id?: string;
	status?: number;
	phoneToken?: string;
}
