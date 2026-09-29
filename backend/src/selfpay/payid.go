package selfpay

import "strings"

func SelfpayID(payID string) string {
	return strings.Replace(payID, "-", "", -1)
}
