package oss

import (
	"bytes"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	errors "github.com/wuntsong-org/wterrors"
)

func UploadInvoice(invoiceID string, file []byte, fast bool, isRed bool) (string, errors.WTError) {
	var key string
	if !utils.IsUID(invoiceID) {
		return "", errors.Errorf("bad invoice id")
	}

	if isRed {
		key = fmt.Sprintf("红字发票/%s.pdf", invoiceID)
	} else {
		key = fmt.Sprintf("蓝字发票/%s.pdf", invoiceID)
	}

	if fast {
		go func() {
			err := InvoiceBucket.PutObject(key, bytes.NewReader(file))
			if err != nil {
				logger.Logger.Error("upload invoice file oss error: %s", err.Error())
			}
		}()
	} else {
		err := InvoiceBucket.PutObject(key, bytes.NewReader(file))
		if err != nil {
			return "", errors.WarpQuick(err)
		}
	}
	return key, nil
}

func GetInvoice(invoiceID string, isRed bool) (string, errors.WTError) {
	var key string
	if !utils.IsUID(invoiceID) {
		return "", errors.Errorf("bad invoice id")
	}

	if isRed {
		key = fmt.Sprintf("红字发票/%s.pdf", invoiceID)
	} else {
		key = fmt.Sprintf("蓝字发票/%s.pdf", invoiceID)
	}

	url, err := InvoiceSignBucket.SignURL(key, oss.HTTPGet, 30)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	return url, nil
}
