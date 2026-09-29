// 列表
export interface announcementTypes {
	page: number;
	pagesize: number;
}
// 新增
/**
 * Title string `json:"title"`
 * Content string `json:"content"`
 * StartAt int64 `json:"startAt"`
 * StopAt int64 `json:"stopAt"`
 */
export interface announcementCreateTypes {
	title: string;
	content: string;
	startAt: number;
	stopAt: number;
}
// 移动
export interface announcementMoveTypes {
	id: number;
	isUp: boolean;
}
