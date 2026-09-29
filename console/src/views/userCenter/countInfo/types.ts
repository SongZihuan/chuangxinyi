export interface checkPhoneTypes {
    phone: string;
    code: string;
    type: string;
}

/**
 * AllowPhone bool `json:"allowPhone"`
 *     AllowEmail bool `json:"allowEmail"`
 *     AllowPassword bool `json:"allowPassword"`
 *     AllowWeChat bool `json:"allowWeChat"`
 *     AllowSecondFA bool `json:"allowSecondFA"`
 */
export interface checkPhoneRes {
    uid?: string;
    id?: string;
    allowPhone: boolean;
    allowEmail: boolean;
    allowPassword: boolean;
    allowWeChat: boolean;
    allowSecondFA: boolean;
}