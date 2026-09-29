package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

func StringCalculateHMAC(key, data string) string {
	// 创建一个新的HMAC-SHA1哈希器，传入密钥
	h := hmac.New(sha1.New, []byte(key))

	// 写入要计算HMAC的数据
	h.Write([]byte(data))

	// 计算HMAC并返回结果，使用Base64编码
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
