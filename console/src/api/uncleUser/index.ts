import request from '/@/utils/request';

export function useUncleUserApi() {
    return {
        //获取协作人信息
        getUncleUserInfo: () => {
            return request({
                url: '/user/center/uncle/list',
                method: 'get',
            });
        },
        // 删除协作人
        delUncleUser: (data: { uncleID: string }) => {
            return request({
                url: '/user/center/uncle/del',
                method: 'post',
                data: data,
            });
        },
        // 添加协作人
        addUncleUser: (data: { uncleID: string }) => {
            return request({
                url: '/user/center/uncle/add',
                method: 'post',
                data: data,
            });
        }
    };
}