package alipay

import "strings"

func AlipayID(payID string) string {
	return strings.Replace(payID, "-", "", -1)
}

func AlipayFaceID(faceID string) string {
	return strings.Replace(faceID, "-", "", -1)
}
