import request from "/@/utils/request";
import {messageUserTypes} from "/@/views/message/user/types";

export function useMessageUserApi(){
    return {
        // 获取消息列表
        getMessageList: (params: messageUserTypes) => {
            return request({
                url: "/user/center/allow-website/message/list",
                method: "get",
                params
            })
        },
        // 已读消息
        readMessage: (data:{id: string}) => {
            return request({
                url: `/user/center/allow-website/message/read/`,
                method: "post",
                data
            })
        },
        // 全部已读
        readAllMessage: () => {
            return request({
                url: `/user/center/allow-website/message/read/all`,
                method: "post",
            })
        }
    }
}