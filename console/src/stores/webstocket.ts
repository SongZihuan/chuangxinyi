import { ref } from 'vue';
import { defineStore } from 'pinia';
import { MessageTypeEnum, MessageTypeMap, SendMessageTypeEnum } from '/@/data/enum';
import { useUserInfo } from '/@/stores/userInfo';
import { Local, Session } from '/@/utils/storage';
import {ElLoading} from 'element-plus';
import { useSocketListInfo } from '/@/stores/socketListInfo';
import { useSocket } from '/@/hooks/useWebstocket';
import {useLoginSignIn, useLoginSignSub} from "/@/hooks/useLoginSignIn";

export const useSocketStore = () => {
	const userInfoStore = useUserInfo();
	const socketListInfoStore = useSocketListInfo();

	return defineStore("socket", () => {
		let isConnected = ref(false)
		let message = ref("")
		let reconnectError = ref(false)
		let heartBeatInterval = ref(50000)
		let heartBeatTimer = ref(0)
		let socketInstance = ref(null)
		let orderData = ref([])
		let tokenStatus = ref(false)
		let sendTokenStatus = ref(false)
		let orderDataFinishStatus = ref(null)


		const resetState = () => {
			isConnected.value = false;
			message.value = '';
			reconnectError.value = false;
			heartBeatInterval.value = 50000;
			heartBeatTimer.value = 0;
			socketInstance.value = null;
			orderData.value = [];
			tokenStatus.value = false;
			localStorage.setItem('socketStatus', 'false');
		}

		const SOCKET_ONOPEN = (event: any) => {
			socketInstance.value = event.currentTarget;
			isConnected.value = true;
			if (Session.get('token')) {
				useSocketStore().sendMessage(MessageTypeEnum.TOKEN, {
					data: JSON.stringify({
						Token: Session.get('token'),
					}),
				});
			}
		}

		const SOCKET_ONCLOSE = () => {
			isConnected.value = false;
			// 连接关闭时停掉心跳消息
			window.clearInterval(heartBeatTimer.value);
			heartBeatTimer.value = 0;
		}

		const SOCKET_ONERROR = () => {
			resetState();
			useSocket();
		}

		// 收到服务端发送的消息
		const SOCKET_ONMESSAGE = async (msg: any) => {
			// 获取用户信息
			const messageTypes = MessageTypeMap[SendMessageTypeEnum.GET_USER_INFO];
			const announcementMessageTypes = MessageTypeMap[SendMessageTypeEnum.GET_ANNOUNCEMENT];
			const messageMessageTypes = MessageTypeMap[SendMessageTypeEnum.GET_MESSAGE];
			const orderMessageTypes = MessageTypeMap[SendMessageTypeEnum.GET_ORDER];

			if (msg.code === MessageTypeEnum.PONG && tokenStatus.value == false) {
				await bindTokenCallback();
			} else if (msg.code === MessageTypeEnum.BAD_TOKEN || msg.code === MessageTypeEnum.LOGOUT) {
				if (Session.get('token') && Session.get('login-token')) {  // 处于登入状态
					const { regetToken } = useLoginSignSub()
					const { isLoginUser, resetSignIn } = useLoginSignIn()

					let res = await regetToken()
					if (!res) {
						if (isLoginUser()) {
							const loadingInstance = ElLoading.service({ fullscreen: true, text: '您已经被登出 ' });
							setTimeout(() => {
								Session.clear(); // 清除浏览器全部临时缓存
								Local.clear();
								// 清除子账号登录信息
								window.location.reload()  // 使用reload，保存路由进度
								loadingInstance.close();
							}, 3000);
						} else {
							let res2 = await resetSignIn()
							if (!res2) {
								const loadingInstance = ElLoading.service({ fullscreen: true, text: '您已经被登出 ' });
								setTimeout(() => {
									Session.clear(); // 清除浏览器全部临时缓存
									Local.clear();
									// 清除子账号登录信息
									window.location.reload()  // 使用reload，保存路由进度
									loadingInstance.close();
								}, 3000);
							}
						}
					}

					// 重新连接
					sendMessageNow(MessageTypeEnum.BYE);
					(socketInstance.value as any).close()
					resetState();
					useSocket();
				}
			} else if (msg.code === MessageTypeEnum.BAD_DATA || msg.code === MessageTypeEnum.BAD_CODE) {
				// 忽略
			} else if (messageTypes.includes(msg.code)) {
				if (msg.code === MessageTypeEnum.UPDATE_USER_INFO) {
					if (Session.get('userInfo')?.role && Session.get('userInfo').role.menus) {
						Session.set('userInfo', { ...msg.data, role: Session.get('userInfo').role });
						Session.set('userData', { ...msg.data, role: Session.get('userInfo').role });
						Local.set('userInfo', { ...msg.data, role: Session.get('userInfo').role });
					} else {
						Session.set('userInfo', msg.data);
						Session.set('userData', msg.data);
						Local.set('userInfo', msg.data);
					}
					await userInfoStore.setUserInfos();
				}
				if (msg.code === MessageTypeEnum.ROLE_CHANGE) {
					// console.log("获取用户角色: " ,message.data)
					await userInfoStore.setRoles(msg.data);
				}
				if (msg.code === MessageTypeEnum.UPDATE_ROLE_INFO) {
					// console.log("获取用户角色: " ,message.data)
					await userInfoStore.setRoles(msg.data);
				}
			} else if (msg.code === MessageTypeEnum.UPDATE_WALLET_INFO) {
				// console.log('获取钱包信息: ', message.data);
				await userInfoStore.setWalletInfos(msg.data);
			} else if (announcementMessageTypes.includes(msg.code)) {
				await socketListInfoStore.setAnnouncementList();
			} else if (messageMessageTypes.includes(msg.code)) {
				await socketListInfoStore.setMessageList();
			} else if (orderMessageTypes.includes(msg.code)) {
				if (msg.code === MessageTypeEnum.NEW_ORDER_REPLY) {
					// console.log('获取工单信息: ', message.data);
					orderData.value = msg.data;
				}
				if (msg.code === MessageTypeEnum.UPDATE_ORDER) {
					// console.log('获取工单信息: ', message.data);
					orderDataFinishStatus.value = msg.data;
				}
			}

			message.value = msg;
		}
		// 自动重连
		const SOCKET_RECONNECT = () => {
			// console.info("消息系统重连中...", count);
		}
		// 重连错误
		const SOCKET_RECONNECT_ERROR = () => {
			reconnectError.value = true;
		}
		// 发送消息
		const sendMessage = (code: any, options?: any) => {
			if (!socketInstance.value) {
				setTimeout(() => {
					sendMessage(code, options);
				}, 1500);
				return;
			}
			if (code === MessageTypeEnum.TOKEN) {
				sendTokenStatus.value = true;
			}

			(socketInstance.value as any).sendObj({
				code: code,
				...options,
			});
		}
		const sendMessageNow = (code: any, options?: any) => {
			if (code === MessageTypeEnum.TOKEN) {
				sendTokenStatus.value = true;
			}

			(socketInstance.value as any).sendObj({
				code: code,
				...options,
			});
		}
		// 绑定token的回调
		const bindTokenCallback = async () => {
			isConnected.value = true;
			tokenStatus.value = true;
			// 连接成功时启动定时发送心跳消息，避免被服务器断开连接
			heartBeatTimer.value = window.setInterval(() => {
				const message = 'PING';
				isConnected &&
				(socketInstance.value as any).sendObj({
					code: message,
				});
			}, heartBeatInterval.value);
			// 连接成功时发送消息
			sendMessage(SendMessageTypeEnum.GET_USER_INFO);
			sendMessage(SendMessageTypeEnum.GET_WALLET_INFO);
			sendMessage(SendMessageTypeEnum.GET_ANNOUNCEMENT);
			sendMessage(SendMessageTypeEnum.GET_MESSAGE);
			await socketListInfoStore.setAnnouncementList();
			await socketListInfoStore.setMessageList();
		}

		return {
			isConnected,
			message,
			reconnectError,
			heartBeatInterval,
			heartBeatTimer,
			orderData,
			tokenStatus,
			socketInstance,
			sendTokenStatus,
			orderDataFinishStatus,

			resetState,
			SOCKET_ONOPEN,
			SOCKET_ONCLOSE,
			SOCKET_ONERROR,
			SOCKET_ONMESSAGE,
			SOCKET_RECONNECT,
			SOCKET_RECONNECT_ERROR,
			sendMessage,
			sendMessageNow,
			bindTokenCallback,
		}
	})()
}
