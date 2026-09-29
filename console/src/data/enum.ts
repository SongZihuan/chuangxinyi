/**
 * UserStatus_Register: "REGISTER",
 *    UserStatus_Normal:   "NORMAL",
 *    UserStatus_Banned:   "BANNED",
 *    UserStatus_Delete:   "DELETE",
 */
export interface userTypes {
	readonly [key: string]: string;
}
export enum UserStatusEnum {
	REGISTER = '注册',
	NORMAL = '正常',
	BANNED = '禁用',
	DELETE = '注销',
	FREEZE = '冻结',
}
export enum userColorTypes {
	REGISTER = 'parimary',
	NORMAL = 'success',
	BANNED = 'warning',
	DELETE = 'error',
	FREEZE = 'info',
}

/**
 * LogoutToken        = "LOGOUT"
 *    UpdateUserInfo     = "UPDATE_USER_INFO"
 *    UpdateWalletInfo   = "UPDATE_WALLET_INFO"
 *    UpdateRoleInfo     = "UPDATE_ROLE_INFO"
 *    NewAnnouncement    = "NEW_ANNOUNCEMENT"
 *    UpdateAnnouncement = "UPDATE_ANNOUNCEMENT"
 *    DeleteAnnouncement = "DELETE_ANNOUNCEMENT"
 *    UpdateMessage      = "NEW_MESSAGE"
 *    ReadMessage        = "READ_MESSAGE"
 *    NewOrderReply      = "NEW_ORDER_REPLY"
 *    UpdateOrder        = "UPDATE_ORDER"
 *    RoleChange         = "ROLE_CHANGE"
 *    Close              = "CLOSE"    // 关闭通道（前端用不到）
 *    BadCode            = "BAD_CODE" // 错误的Code
 *    BadData            = "BAD_DATA" // 错误的Data
 *    Pong               = "PONG"
 */
export interface messageType {
	readonly [key: string]: string;
}
// 订阅返回的消息类型
export enum MessageTypeEnum {
	LOGOUT = 'LOGOUT',
	UPDATE_USER_INFO = 'UPDATE_USER_INFO',
	UPDATE_WALLET_INFO = 'UPDATE_WALLET_INFO',
	UPDATE_ROLE_INFO = 'UPDATE_ROLE_INFO',
	NEW_ANNOUNCEMENT = 'NEW_ANNOUNCEMENT',
	UPDATE_ANNOUNCEMENT = 'UPDATE_ANNOUNCEMENT',
	DELETE_ANNOUNCEMENT = 'DELETE_ANNOUNCEMENT',
	NEW_MESSAGE = 'NEW_MESSAGE',
	READ_MESSAGE = 'READ_MESSAGE',
	NEW_ORDER_REPLY = 'NEW_ORDER_REPLY',
	UPDATE_ORDER = 'UPDATE_ORDER',
	ROLE_CHANGE = 'ROLE_CHANGE',
	CLOSE = 'CLOSE',
	BAD_CODE = 'BAD_CODE',
	BAD_DATA = 'BAD_DATA',
	PONG = 'PONG',
	TOKEN = 'TOKEN',
	BYE = "BYE",
	BAD_TOKEN = "BAD_TOKEN"
}

/**
 * 发送消息类型
 * GetUserInfo     = "GET_USER_INFO"    // 订阅 UpdateUserInfo RoleChange UpdateRoleInfo
 *    GetWalletInfo   = "GET_WALLET_INFO"  // 订阅 UpdateWalletInfo
 *    GetAnnouncement = "GET_ANNOUNCEMENT" // 订阅 NewAnnouncement, UpdateAnnouncement, DeleteAnnouncement
 *    GetMessage      = "GET_MESSAGE"      // 订阅 UpdateMessage
 *    GetOrder        = "GET_ORDER"        // 订阅 NewOrderReply UpdateOrder
 *    Ping            = "PING"
 */
export enum SendMessageTypeEnum {
	GET_USER_INFO = 'GET_USER_INFO',
	GET_WALLET_INFO = 'GET_WALLET_INFO',
	GET_ANNOUNCEMENT = 'GET_ANNOUNCEMENT',
	GET_MESSAGE = 'GET_MESSAGE',
	GET_ORDER = 'GET_ORDER',
	GET_TOKEN_INFO = 'GET_TOKEN_INFO',
	PING = 'PING',
}

/**
 * 发送和接收的对应关系
 */
export const MessageTypeMap = {
	[SendMessageTypeEnum.GET_TOKEN_INFO]: MessageTypeEnum.LOGOUT,
	[SendMessageTypeEnum.GET_USER_INFO]: [MessageTypeEnum.UPDATE_USER_INFO, MessageTypeEnum.ROLE_CHANGE, MessageTypeEnum.UPDATE_ROLE_INFO],
	[SendMessageTypeEnum.GET_ANNOUNCEMENT]: [
		MessageTypeEnum.NEW_ANNOUNCEMENT,
		MessageTypeEnum.UPDATE_ANNOUNCEMENT,
		MessageTypeEnum.DELETE_ANNOUNCEMENT,
	],
	[SendMessageTypeEnum.GET_MESSAGE]: [MessageTypeEnum.NEW_MESSAGE, MessageTypeEnum.READ_MESSAGE],
	[SendMessageTypeEnum.GET_ORDER]: [MessageTypeEnum.NEW_ORDER_REPLY, MessageTypeEnum.UPDATE_ORDER],
	[SendMessageTypeEnum.GET_WALLET_INFO]: MessageTypeEnum.UPDATE_WALLET_INFO,
	[SendMessageTypeEnum.PING]: MessageTypeEnum.PONG,
};
