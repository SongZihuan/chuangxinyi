/**
 * UserID int64 `json:"userID"`
 *     Title string `json:"title"`
 *     Content string `json:"content"`
 *     Sender string `json:"sender"`
 *     SenderID int64 `json:"senderID"`
 *     SenderLink string `json:"senderLink"`
 *     CreateAt int64  `json:"createAt"`
 *     ReadAt int64 `json:"readAt"`
 */
export interface messageAdminDataType {
	id: string;
	userID: string;
	title: string;
	content: string;
	sender: string;
	senderID: string;
	senderLink: string;
	createAt: number;
	readAt: number;
}

/**
 * Phone         string         `json:"phone"`
 *     Sig           string         `json:"sig"`
 *     Template      string         `json:"template"`
 *     TemplateParam []LabelInterfaceValueRecord         `json:"templateParam"`
 *     SenderId      int64  `json:"senderID"`
 *     Success       bool           `json:"success"`
 *     ErrorMsg      string `json:"errorMsg"`
 *     CreateAt      int64      `json:"createAt"`
 */
export interface smsMessageAdminDataType {
	phone: string;
	sig: string;
	template: string;
	templateParam: any[];
	senderId: string;
	success: boolean;
	errorMsg: string;
	createAt: number;
}

/**
 * Email    string         `json:"email"`
 *     Subject  string         `json:"subject"`
 *     Content  string         `json:"content"`
 *     Sender   string         `json:"sender"`
 *     SenderId int64  `json:"senderIDd"`
 *     Success  bool           `json:"success"`
 *     ErrorMsg string `json:"errorMsg"`
 *     CreateAt int64      `json:"createAt"`
 */
export interface emailMessageAdminDataType {
	email: string;
	subject: string;
	content: string;
	sender: string;
	senderIDd: string;
	success: boolean;
	errorMsg: string;
	createAt: number;
}

/**
 * OpenID   string         `json:"openID"`
 *     Template string         `json:"template"`
 *     Url      string         `json:"url"`
 *     Val      []LabelInterfaceValueRecord        `json:"val"`
 *     SenderId int64  `json:"senderID"`
 *     Success  bool           `json:"success"`
 *     ErrorMsg string `json:"errorMsg"`
 *     CreateAt int64      `json:"createAt"`
 */
export interface fuwuhaoMessageAdminDataType {
	openID: string;
	template: string;
	url: string;
	val: any[];
	senderID: string;
	success: boolean;
	errorMsg: string;
	createAt: number;
}

/**
 * Webhook  string         `json:"webhook"`
 *     Text     string         `json:"text"`
 *     AtAll    bool           `json:"atAll"`
 *     SenderId int64  `json:"senderId"`
 *     Success  bool           `json:"success"`
 *     ErrorMsg string `json:"errorMsg"`
 *     CreateAt int64     `json:"createAt"`
 */
export interface wxrobotMessageAdminDataType {
	webhook: string;
	text: string;
	atAll: boolean;
	senderId: string;
	success: boolean;
	errorMsg: string;
	createAt: number;
}
export interface messageAdminTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	timetype?: number;
	range: Date[];
	id?: string;
	uid?: string;
	phone?: string;
	email?: string;
	openID?: string;
	webhook?: string;
	senderID?: number;
}
export interface messageAdminStateType {
	data: messageAdminDataType[] | smsMessageAdminDataType[] | emailMessageAdminDataType[] | fuwuhaoMessageAdminDataType[] | wxrobotMessageAdminDataType[];
	total: number;
	loading: boolean;
	param: messageAdminTypes;
}
export interface messageAdminState {
	tableData: messageAdminStateType;
}
