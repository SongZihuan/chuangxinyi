import request from '/@/utils/request';

export function useNephewUserApi() {
    return {
        //获取协作人信息
        getNephewUserInfo: () => {
            return request({
                url: '/user/center/nephew/list',
                method: 'get',
            });
        },
        // 删除协作人
        delNephewUser: (data: { nephewID: string }) => {
            return request({
                url: '/user/center/nephew/del',
                method: 'post',
                data: data,
            });
        },
        // 添加协作人
        addNephewUser: (data: { nephewID: string }) => {
            return request({
                url: '/user/center/nephew/add',
                method: 'post',
                data: data,
            });
        },

    };
}