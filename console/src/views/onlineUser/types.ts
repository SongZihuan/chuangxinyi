/**
 * id: '',
 *   roleID: 0,
 *   roleName: '',
 *   roleSign: '',
 *   phone: '',
 *   userName: '',
 *   nickname: '',
 *   header: '',
 *   email: '',
 *   userRealName: '',
 *   companyName: '',
 *   wechatNickName: '',
 *   wechatHeader: '',
 *   unionID: '',
 *   signin: false,
 *   status: '',
 *   inviteCount: 0,
 *   tokenExpire: 0,
 *   createAt: 0,
 */
export interface fatherUserType {
    id: string;
    roleID: number;
    roleName: string;
    roleSign: string;
    phone: string;
    userName: string;
    nickname: string;
    header: string;
    email: string;
    userRealName: string;
    companyName: string;
    wechatNickName: string;
    wechatHeader: string;
    unionID: string;
    signin: boolean;
    status: string;
    inviteCount: number;
    tokenExpire: number;
    createAt: number;
}