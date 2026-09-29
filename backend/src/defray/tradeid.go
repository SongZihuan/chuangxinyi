package defray

import "strings"

func TradeID(payID string) string {
	return strings.Replace(payID, "-", "", -1)
}
