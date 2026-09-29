import request from '/@/utils/request';
import {registerSubUserTypes} from "/@/api/subUser/types";

export function useSubUserCenterApi() {
    return {
        //获取子账号列表
        getSubUserList: () => {
            return request({
                url: '/user/center/son/tree',
                method: 'get',
                isLoginToken: false,
            } as any);
        },
        getSubUserListByLoginToken: () => {
            return request({
                url: '/user/center/son/tree',
                method: 'get',
                isLoginToken: true,
            } as any);
        },
        // 注册子账号
        registerSubUser: (data: registerSubUserTypes) => {
            return request({
                url: '/user/center/son/register',
                method: 'post',
                data,
            });
        }
    }
}