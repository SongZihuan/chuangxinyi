/**
 * IsOr bool `json:"isOr"`
 *     Policy []string `json:"policy"`
 *     MenuID int64 `json:"menuID"`
 */
export interface MenuRauthTypes {
	isOr: boolean;
	policy: string[];
	menuID: number;
}

/**
 * policy: [],
 *   roles: [],
 *   关系展示
 */
export interface MenuRauthShowTypes {
	policy: string[];
	roles: any[];
}
