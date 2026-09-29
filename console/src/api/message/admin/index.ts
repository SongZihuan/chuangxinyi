import request from '/@/utils/request';
import {
	emailMessageAdminTypes,
	fuwuhaoMessageAdminTypes,
	messageAdminTypes,
	smsMessageAdminTypes,
	wxrobotMessageAdminTypes,
	sendMsgTypes,
	sendEmailTypes,
	sendSmsTypes,
	sendwxrobotTypes,
	sendFuwuhaoTypes,
} from '/@/api/message/admin/types';

export function useMessageAdminApi() {
	return {
		// 获取消息列表
		getMessageList: (params: messageAdminTypes) => {
			return request({
				url: '/admin/msg/allow-website/message/msg/list',
				method: 'get',
				params,
			});
		},
		//发送站内信
		sendMsg: (data: sendMsgTypes) => {
			return request({
				url: '/admin/msg/allow-website/send/msg',
				method: 'post',
				data,
			});
		},
		// 获取sms消息列表
		getSmsMessageList: (params: smsMessageAdminTypes) => {
			return request({
				url: '/admin/msg/allow-website/message/sms/list',
				method: 'get',
				params,
			});
		},
		//发送短信
		sendSms: (data: sendSmsTypes) => {
			return request({
				url: '/admin/msg/allow-website/send/sms',
				method: 'post',
				data,
			});
		},
		// 获取email消息列表
		getEmailMessageList: (params: emailMessageAdminTypes) => {
			return request({
				url: '/admin/msg/allow-website/message/email/list',
				method: 'get',
				params,
			});
		},
		//发送邮件
		sendEmail: (data: sendEmailTypes) => {
			return request({
				url: '/admin/msg/allow-website/send/email',
				method: 'post',
				data,
			});
		},
		// 获取fuwuhao消息列表
		getFuwuhaoMessageList: (params: fuwuhaoMessageAdminTypes) => {
			return request({
				url: '/admin/msg/allow-website/message/fuwuhao/list',
				method: 'get',
				params,
			});
		},
		//发送服务号消息
		sendFuwuhao: (data: sendFuwuhaoTypes) => {
			return request({
				url: '/admin/msg/allow-website/send/fuwuhao',
				method: 'post',
				data,
			});
		},
		// 获取wxrobot消息列表
		getWxrobotMessageList: (params: wxrobotMessageAdminTypes) => {
			return request({
				url: '/admin/msg/allow-website/message/wxrobot/list',
				method: 'get',
				params,
			});
		},
		//发送企业微信机器人
		sendWxrobot: (data: sendwxrobotTypes) => {
			return request({
				url: '/admin/msg/allow-website/send/wxrobot',
				method: 'post',
				data,
			});
		},
		//
	};
}
