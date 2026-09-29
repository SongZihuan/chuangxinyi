import request from '/@/utils/request';
import type { phoneCodeTypes, sliderHeadersTypes, checkEmailTypes } from '../base/types';
import { checkFaceRecognitionTypes, faceRecognitionUrlTypes, startFaceRecognitionTypes } from '/@/api/login/types';

export function useLoginApi() {
	return {
		// 退出登录
		logout: () => {
			return request({
				url: '/user/center/token/delete',
				method: 'post',
			});
		},
		//手机登录
		phoneLogin: (data: phoneCodeTypes, headers: sliderHeadersTypes) => {
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
		//邮箱登录
		emailLogin: (data: checkEmailTypes, headers: sliderHeadersTypes) => {
			return request({
				url: '/public/checker/allow-website/email',
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
		//账号登录
		accountLogin: (data: any) => {
			delete data.password
			return request({
				url: '/public/checker/captcha/password',
				method: 'post',
				data,
				headers: {
					'X-CAPTCHA-Token': data.headers.token,
					'X-CAPTCHA-Sig': data.headers.sig,
					'X-CAPTCHA-SessionId': data.headers.sessionId,
					'X-CAPTCHA-Scene': 'nc_other',
				},
			});
		},
		//身份证号找回密码
		idcardLogin: (data: any) => {
			return request({
				url: '/public/checker/captcha/idcard',
				method: 'post',
				data,
			});
		},
		//企业身份证法人找回密码
		legalPersonIDCardLogin: (data: any) => {
			return request({
				url: '/public/checker/captcha/legalperson',
				method: 'post',
				data,
			});
		},
		// 开始人脸识别二维码
		startFaceRecognition: (data: startFaceRecognitionTypes) => {
			return request({
				url: '/public/checker-before/allow-website/alipay/face',
				method: 'post',
				data,
			});
		},
		// 获取人脸识别二维码url
		getFaceRecognitionUrl: (params: faceRecognitionUrlTypes) => {
			return request({
				url: '/public/checker-before/allow-website/alipay/face/url',
				method: 'get',
				params,
			});
		},
		// 检测人脸识别是否成功
		checkFaceRecognition: (data: checkFaceRecognitionTypes) => {
			return request({
				url: '/public/checker/allow-website/face',
				method: 'post',
				data,
			});
		},
	};
}
