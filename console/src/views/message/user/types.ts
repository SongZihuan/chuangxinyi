/**
 * type Message {
 *     ID int64 `json:"id"`
 *     Title string `json:"title"`
 *     Content string `json:"content"`
 *     Sender string `json:"sender"`
 *     SenderLink string `json:"senderLink"`
 *     CreateAt int64  `json:"createAt"`
 *     ReadAt int64 `json:"readAt"`
 * }
 */
export interface messageUserDataType {
	id: string;
	title: string;
	content: string;
	sender: string;
	senderLink: string;
	createAt: number;
	readAt: number;
}
export interface messageUserTypes {
	page: number;
	pagesize: number;
	starttime?: number;
	endtime?: number;
	range?: Date[];
	timetype?: number;
	senderID?: number,
}
export interface messageUserStateType {
	data: messageUserDataType[];
	total: number;
	loading: boolean;
	param: messageUserTypes;
}
export interface messageUserState {
	tableData: messageUserStateType;
}
