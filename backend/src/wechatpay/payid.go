package wechatpay

import "strings"

func WeChatPayID(payID string) string {
	return strings.Replace(payID, "-", "", -1)
}
