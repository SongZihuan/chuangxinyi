import { authApi, AuthApiType } from '/@/utils/authFunction';
import { Session, Local } from '/@/utils/storage';

export default function useSubAuth(): AuthApiType {
	const subAuth = Session.get('userType') || Local.get('userType');
	if (subAuth && subAuth.subType) {
		return authApi(subAuth.subType);
	} else {
		return null;
	}
}
