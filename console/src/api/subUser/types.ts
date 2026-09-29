/**
 * 注册子用户
 * PhoneToken string `json:"phoneToken"`
 *     RoleID int64 `json:"roleID"`
 */
export interface registerSubUserTypes {
	phoneToken: string;
	roleID?: number | undefined;
}
