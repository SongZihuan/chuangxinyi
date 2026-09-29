import request from '/@/utils/request';
import {invoiceTypes, invoiceUserTypes} from './types';

export function useInvoiceUserApi() {
    return {
        getInvoiceList(params: invoiceUserTypes) {
            return request({
                url: '/user/center/invoice/list',
                method: 'get',
                params,
            });
        },
        Invoice(data: invoiceTypes) {
            return request({
                url: '/user/center/invoice',
                method: 'post',
                data
            })
        },
        getInvoice(params: { invoiceID: string }) {
            return request({
                url: '/user/center/invoice/get',
                method: 'get',
                params
            })
        },
        // 下载发票
        downloadInvoice(params: { invoiceID: string }) {
            return request({
                url: '/user/center/invoice/download',
                method: 'get',
                params
            })
        },
    };
}