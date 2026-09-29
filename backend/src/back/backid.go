package back

import "strings"

func BackID(payID string) string {
	return strings.Replace(payID, "-", "", -1)
}
