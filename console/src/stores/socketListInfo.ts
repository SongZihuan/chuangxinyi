import { defineStore } from 'pinia';
import { useAnnouncementApi } from '/@/api/announcement';
import { useMessageUserApi } from '/@/api/message/user';
import dayjs from "dayjs";

/**
 * socket数据存取
 */
export const useSocketListInfo = defineStore({
	id: 'socketListInfo',
	state: () => ({
		// 公告列表
		announcementList: [],
		// 消息列表
		messageList: [],
		// 订单列表
		orderList: [],
		// 用户信息
		userInfo: {},
		// 钱包信息
		walletInfo: {},
		// params
		params: {},
        // 用户order
        userOrderData: []
	}),
	actions: {
		// 设置公告列表
		async setAnnouncementList() {
			const params = {
				page: 1,
				pagesize: 3,
			};
			await useAnnouncementApi()
				.userAnnouncementList(params)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						this.announcementList = res.data.announcement;
					}
				});
		},
		// 设置消息列表
		async setMessageList() {
			const param = {
				page: 1,
				pagesize: 3,
				senderID: -1,
			};
			useMessageUserApi()
				.getMessageList(param)
				.then((res: any) => {
					if (res.code === "SUCCESS") {
						this.messageList = res.data.message;
					}
				});
		},
        // 用户order
        async setUserOrderData() {

        }
	},
});
