/**
 * Name   string `json:"name"`       // 姓名，公司名
 *     TaxID  string `json:"tax_id"`   // 税号，身份证号
 *     BandID string `json:"band_id"` // 银行卡号
 *     Band   string `json:"band"`       // 开户行
 */
export interface invoiceTitleTypes {
	name: string;
	taxID: string;
	bankID: string;
	bank: string;
}
