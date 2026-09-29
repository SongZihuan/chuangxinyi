import request from '/@/utils/request';
import { announcementCreateTypes, announcementMoveTypes, announcementTypes } from '/@/api/announcement/types';

export function useAnnouncementApi() {
	return {
		//获取公告列表
		announcementList: (params: announcementTypes) => {
			return request({
				url: '/public/announcement/list',
				method: 'get',
				params,
			});
		},
		// 获取公告列表(管理员)
		getAnnouncementList: (params: announcementTypes) => {
			return request({
				url: '/admin/announcement/list',
				method: 'get',
				params,
			});
		},
		// 获取公告列表(管理员)
		userAnnouncementList: (params: announcementTypes) => {
			return request({
				url: '/public/announcement/list',
				method: 'get',
				params,
			});
		},
		//新增公告
		createAnnouncement: (data: announcementCreateTypes) => {
			return request({
				url: '/admin/announcement/create',
				method: 'post',
				data,
			});
		},
		//更新公告
		updateAnnouncement: (data: announcementCreateTypes) => {
			return request({
				url: '/admin/announcement/update',
				method: 'post',
				data,
			});
		},
		//删除公告
		delAnnouncement: (data: { id: string }) => {
			return request({
				url: '/admin/announcement/delete',
				method: 'post',
				data,
			});
		},
		// 移动
		moveAnnouncement: (data: announcementMoveTypes) => {
			return request({
				url: '/admin/announcement/move',
				method: 'post',
				data,
			});
		},
	};
}
