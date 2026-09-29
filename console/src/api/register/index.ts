import request from '/@/utils/request';
import type { checkPhoneCodeTypes, uploadUserInfoJsontTypes, startRegistrants, uploadCompanyInfoJsonTypes, updatePasswordTypes } from './types';
import type { phoneCodeTypes, sliderHeadersTypes, checkEmailTypes } from '../base/types';
import { useUserInfo } from '/@/stores/userInfo';
import { updateUserNameTypes } from './types';
const stores = useUserInfo();
export function useRegisterApi() {
	return {
		//发送手机验证码
		sendPhoneCode: (data: phoneCodeTypes, headers: sliderHeadersTypes) => {
			return request({
				url: '/public/checker-before/phone/send',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Token': headers.token,
					'X-CAPTCHA-Sig': headers.sig,
					'X-CAPTCHA-SessionId': headers.sessionId,
					'X-CAPTCHA-Scene': 'nc_other',
				},
			});
		},

		//验证手机验证码
		checkPhoneCode: (data: checkPhoneCodeTypes, headers: sliderHeadersTypes) => {
			return request({
				url: '/public/checker/allow-website/phone',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Token': headers.token,
					'X-CAPTCHA-Sig': headers.sig,
					'X-CAPTCHA-SessionId': headers.sessionId,
					'X-CAPTCHA-Scene': 'nc_other',
				},
			});
		},
		//验证邮箱验证码
		checkEmailCode: (data: checkEmailTypes, headers: sliderHeadersTypes) => {
			return request({
				url: '/public/checker/allow-website/email ',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Token': headers.token,
					'X-CAPTCHA-Sig': headers.sig,
					'X-CAPTCHA-SessionId': headers.sessionId,
					'X-CAPTCHA-Scene': 'nc_other',
				},
			});
		},
		//开始注册流程
		startRegistrants: (data: startRegistrants) => {
			return request({
				url: '/public/registrants',
				method: 'post',
				data,
			});
		},
		//上传用户身份信息(照片)
		userinfoUpload: (data: any) => {
			return request({
				url: '/user/center/root-only/userinfo/upload',
				method: 'post',
				data,
			});
		},
		//填写身份证姓名和身份证号
		uploadUserInfoJson: (data: uploadUserInfoJsontTypes) => {
			return request({
				url: '/user/center/root-only/userinfo/json/upload',
				method: 'post',
				data,
				headers: {
					'x-token': stores.userInfos.XToken,
				},
			});
		},
		// 个人信息人脸识别
		uploadUserInfoFace: (data: { faceToken: string }) => {
			return request({
				url: '/user/center/root-only/userinfo/face/upload',
				method: 'post',
				data,
			});
		},
		// 上传企业信息(照片)
		companyinfoUpload: (data: any) => {
			return request({
				url: '/user/center/root-only/companyinfo/upload',
				method: 'post',
				data,
			});
		},
		// 企业信息人脸识别
		uploadCompanyInfoFace: (data: { faceToken: string }) => {
			return request({
				url: '/user/center/root-only/companyinfo/face/upload',
				method: 'post',
				data,
			});
		},
		//上传企业法人身份证和企业营业执照
		uploadCompanyInfo: (data: any) => {
			return request({
				url: '/user/center/root-only/companyinfo/upload',
				method: 'post',
				data,
			});
		},
		// 上传身份证反正面
		uploadUserInfoBack: (data: any) => {
			return request({
				url: '/public/checker/allow-website/user',
				method: 'post',
				data,
				headers: {
					'Content-Type': 'multipart/form-data',
					'x-token': stores.userInfos.XToken,
				},
			});
		},
		// 上传企业身份证正面
		uploadCompanyInfoFront: (data: any) => {
			return request({
				url: '/public/checker/allow-website/company',
				method: 'post',
				data,
				headers: {
					'Content-Type': 'multipart/form-data',
					'x-token': stores.userInfos.XToken,
				},
			});
		},
		//填写企业法人身份证和身份证号
		uploadCompanyInfoJson: (data: uploadCompanyInfoJsonTypes) => {
			return request({
				url: '/user/center/root-only/companyinfo/json/upload',
				method: 'post',
				data,
			});
		},
		//设置密码
		updatePassword: (data: updatePasswordTypes) => {
			return request({
				url: '/user/center/password/update',
				method: 'post',
				data,
				headers: {
					'x-token': stores.userInfos.XToken,
				},
			});
		},
		//获取微信配置
		weixinConfig: () => {
			return request({
				url: '/public/checker-before/allow-website/wechat/get',
				method: 'get',
			});
		},
		// 更新用户名
		updateUserName: (data: updateUserNameTypes) => {
			return request({
				url: '/user/center/username/update',
				method: 'post',
				data,
			});
		},
		//更新昵称
		updateNickname: (data: updateUserNameTypes) => {
			return request({
				url: '/user/center/nickname/update',
				method: 'post',
				data,
			});
		},

		// 更新微信机器人
		updateWxRobot: (data: any) => {
			return request({
				url: '/user/center/wxrobot/update',
				method: 'post',
				data,
			});
		},
		//更新token时长
		tokenExpiration: (data: { tokenExpiration: string }) => {
			return request({
				url: '/user/center/token/expiration/update',
				method: 'post',
				data,
			});
		},
		// 更新登录控制
		updateLoginController: (data: any) => {
			return request({
				url: '/user/center/loginctrl/update',
				method: 'post',
				data,
			});
		}
	};
}
