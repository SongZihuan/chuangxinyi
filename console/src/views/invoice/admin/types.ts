/**
 * {
 *     "name": "某某科技有限公司",
 *     "taxID": "91xxxxxxxxxx00000X",
 *     "bandID": "00000",
 *     "band": "某某银行"
 * }
 */
interface invoiceUserDataType {
	name: string;
	taxID: string;
	bandID: string;
	band: string;
}

interface invoiceAddressDataType {
	name: string;
	phone: string;
	email: string;
	province: string;
	city: string;
	district: string;
	address: string;
}
/**
 * InvoiceID string `json:"invoiceID"`
 *     Type int64 `json:"type"`  // 1-个人普票 2-企业普票 3-企业专票
 *     Title Title `json:"title"`
 *     Address InvoiceAddress `json:"address"`
 *     Amount int64 `json:"amount"`
 *     Status int64 `json:"status"`  // 1-待开票 2-已开票 3-已退票 4-信息错误
 *     CreateAt int64 `json:"createAt"`
 *     BillingAt int64 `json:"billingAt,omitempty"`
 *     ReturnAt int64 `json:"returnAt,omitempty"`
 */
export interface invoiceAdminDataType {
	invoiceID: string;
	type: number;
	title: invoiceUserDataType;
	address: invoiceAddressDataType;
	amount: number;
	status: number;
	createAt: number;
	billingAt: number;
	returnAt: number;
}
export interface invoiceAdminTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
	range: Date[];
}
export interface invoiceAdminStateType {
	data: invoiceAdminDataType[];
	total: number;
	loading: boolean;
	param: invoiceAdminTypes;
}
export interface invoiceAdminState {
	tableData: invoiceAdminStateType;
}
