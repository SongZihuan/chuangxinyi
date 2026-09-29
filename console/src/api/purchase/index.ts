import request from '/@/utils/request';
export function purchaseApi(){
    return {
        // 购买
        purchase: (data: any) => {
            return request({
                url: '/user/center/defray',
                method: 'post',
                data,
            });
        },
        // 获取订单信息
        getOrderInfo: (params: any) => {
            return request({
                url: '/public/defray/',
                method: 'get',
                params,
            });
        },
    }

}