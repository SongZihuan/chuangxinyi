import request from '/@/utils/request';
import {subUserInfoTypes} from "/@/api/user/subuser/types";
export function useSubUserApi(){
    return {
        //获取子账户信息
        subUserInfo: (params: subUserInfoTypes) => {
            return request({
                url: '/admin/user/son',
                method: 'get',
                params,
            });
        },
        // 获取token
        getToken: (data:{id: string}) => {
            return request({
                url: '/user/center/son/token',
                method: 'post',
                data,
                isLoginToken: true
            } as any);
        }
    }
}