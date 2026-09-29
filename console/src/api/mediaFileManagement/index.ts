import request from '/@/utils/request';
export function useMediaFileApi() {
	return {
		// 获取文件列表
		getMediaFileList(params) {
			return request({
				url: 'admin/ui/file/list',
				method: 'get',
				params,
			});
		},
		// 删除文件
		deleteMediaFile(data) {
			return request({
				url: 'admin/ui/file/delete',
				method: 'post',
				data,
			});
		},
		// 新增文件
		addMediaFile(data) {
			return request({
				url: 'admin/ui/file/update',
				method: 'post',
				headers: {
					'Content-Type': 'multipart/form-data',
				},
				data,
			});
		},
	};
}
