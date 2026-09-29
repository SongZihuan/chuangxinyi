import { useUserInfo } from '/@/stores/userInfo';
import { judementSameArr } from '/@/utils/arrayOperation';

/**
 * 单个权限验证
 * @param value 权限值
 * @returns 有权限，返回 `true`，反之则反
 */
export function auth(value: string): boolean {
	const stores = useUserInfo();
	return stores.userInfos.authBtnList.some((v: string) => v === value);
}

/**
 * 多个权限验证，满足一个则为 true
 * @param value 权限值
 * @returns 有权限，返回 `true`，反之则反
 */
export function auths(value: Array<string>): boolean {
	let flag = false;
	const stores = useUserInfo();
	stores.userInfos.authBtnList.map((val: string) => {
		value.map((v: string) => {
			if (val === v) flag = true;
		});
	});
	return flag;
}

/**
 * 多个权限验证，全部满足则为 true
 * @param value 权限值
 * @returns 有权限，返回 `true`，反之则反
 */
export function authAll(value: Array<string>): boolean {
	const stores = useUserInfo();
	return judementSameArr(value, stores.userInfos.authBtnList);
}

/**
 * 接口权限验证
 * UserRootToken 非子用户（非高权限）
 * UserSonToken 子用户（直接登录）
 * UserFatherToken 父用户登录的子用户（非根父亲）
 * UserRootFatherToken 根父用户登录的子用户
 * UserUncleToken 协作人登录的子用户
 * UserHighAuthorityRootToken 非子用户（高权限）
 *
 * 返回：addUncle: uncle不允许
 * aliPay（PC和WAP）: 非子用户或子用户根父亲
 * bind2FA: uncle不允许
 * delete2FA: uncle不允许，只有非子用户（非高权限）或子用户直接登录的时候需要验证密钥
 * deleteAllToken: uncle不允许
 * deleteUser: 仅限非子用户，或者子用户的父用户
 * delUncle: uncle不允许
 * getAllToken: uncle不允许
 * getinfo: 仅限非是子用户，根父用户可以获得法人身份证号码和企业社会统一信用代码
 * getSonToken: uncle不允许
 * registerSon: uncle不允许
 * selfPay: 同alipay
 * updateAdress: uncle不允许
 * updateEmail: uncle不允许
 * updateHeader: uncle不允许
 * updateHomepage: uncle不允许
 * updatePassword: 非子用户或子用户的父亲
 * updatePhone: 非子用户或子用户的父亲
 * updateSignOne:  非子用户或子用户的父亲
 * updateTitle: uncle不允许
 * updateTokenExpiration: 非子用户或子用户的父亲
 * updateUserName: 非子用户或子用户父亲
 * updateWechat: 非子用户或子用户父亲
 * updateWxrobot: 非子用户或子用户父亲
 * updateRealName：非子用户或子用户根父亲
 * wechatPay: 同alipay
 * UpdateLoginController: 非子用户或子用户父亲
 * is2FA: 根父账户和非子用户（含高权限）无需验证
 */
export interface AuthApiType {
	addUncle: boolean;
	aliPay: boolean;
	bind2FA: boolean;
	delete2FA: boolean;
	delete2FAOther: boolean;
	deleteAllToken: boolean;
	deleteUser: boolean;
	delUncle: boolean;
	getAllToken: boolean;
	getinfo: boolean;
	getSonToken: boolean;
	registerSon: boolean;
	selfPay: boolean;
	updateAdress: boolean;
	updateEmail: boolean;
	updateHeader: boolean;
	updateHomepage: boolean;
	updatePassword: boolean;
	updatePhone: boolean;
	updateSignOne: boolean;
	updateTitle: boolean;
	updateWechat: boolean;
	updateRealName: boolean;
	wechatPay: boolean;
	updateWxrobot: boolean;
	updateTokenExpiration: boolean;
	updateUserName: boolean;
	UpdateLoginController: boolean;
	is2FA: boolean;
}
export function authApi(SubType: string): AuthApiType {
	return {
		addUncle: SubType !== 'UserUncleToken',
		aliPay: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken',
		bind2FA: SubType !== 'UserUncleToken',
		delete2FA: SubType !== 'UserUncleToken',
		delete2FAOther: SubType === 'UserHighAuthorityRootToken' || SubType === 'UserRootToken' || SubType === 'UserRootFatherToken',
		deleteAllToken: SubType !== 'UserUncleToken',
		deleteUser: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		delUncle: SubType !== 'UserUncleToken',
		getAllToken: SubType !== 'UserUncleToken',
		getinfo: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken',
		getSonToken: SubType !== 'UserUncleToken',
		registerSon: SubType !== 'UserUncleToken',
		selfPay: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken',
		updateAdress: SubType !== 'UserUncleToken',
		updateEmail: SubType !== 'UserUncleToken',
		updateHeader: SubType !== 'UserUncleToken',
		updateHomepage: SubType !== 'UserUncleToken',
		updatePassword: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		updatePhone: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		updateSignOne: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' ||	SubType === 'UserFatherToken',
		updateTitle: SubType !== 'UserUncleToken',
		updateTokenExpiration: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		updateUserName: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		updateWechat: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		updateWxrobot: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		updateRealName: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken'|| SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		wechatPay: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken' || SubType === 'UserFatherToken',
		UpdateLoginController: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken',
		is2FA: SubType === 'UserRootToken' || SubType === 'UserRootFatherToken' || SubType === 'UserHighAuthorityRootToken',
	};
}
