import request from '/@/utils/request';
import {invoiceTitleTypes} from "/@/api/invoice/title/types";

export function useInvoiceTitleApi() {
    return {
        updateInvoice(data: invoiceTitleTypes) {
            return request({
                url: '/user/center/title/update',
                method: 'post',
                data
            })
        }
    };
}