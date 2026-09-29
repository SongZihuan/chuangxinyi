import request from '/@/utils/request';

export function useWalletRecordApi() {
    return {
        //新增站点域名
        walletList: (data: any) => {
            return request({
                url: '/admin/finance/wallet/record/list',
                method: 'get',
                data,
            });
        },
        userWalletList: (data: any) => {
            return request({
                url: '/user/center/wallet/record/list',
                method: 'get',
                data,
            });
        },
        // 更新备注
        updateRemark: (data: any) => {
            return request({
                url: '/admin/finance/wallet/record/remark/update',
                method: 'post',
                data,
            });
        },
        // 微信提现
        wechatWithdraw: (data: any) => {
            return request({
                url: '/user/center/withdraw/wechatpay',
                method: 'post',
                data,
            });
        },
        // 支付宝提现
        alipayWithdraw: (data: any) => {
            return request({
                url: '/user/center/withdraw/alipay',
                method: 'post',
                data,
            });
        },
        // 自定义提现
        selfpayWithdraw: (data: any) => {
            return request({
                url: '/user/center/withdraw/selfpay',
                method: 'post',
                data,
            });
        },
    };
}
