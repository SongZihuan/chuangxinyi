package where

import (
	"fmt"
	"strings"
	"time"
)

func JoinIntToString(lst []int64, seq string, removeDuplicates bool) string {
	if len(lst) == 0 {
		return ""
	}

	if removeDuplicates {
		lst = RemoveDuplicates(lst)
	}

	lstString := make([]string, 0, len(lst))
	for _, i := range lst {
		lstString = append(lstString, fmt.Sprintf("%d", i))
	}

	return strings.Join(lstString, seq)
}

func JoinStringWithQuotaToString(lst []string, seq string, removeDuplicates bool) string {
	if len(lst) == 0 {
		return ""
	}

	if removeDuplicates {
		lst = RemoveDuplicates(lst)
	}

	lstString := make([]string, 0, len(lst))
	for _, i := range lst {
		lstString = append(lstString, fmt.Sprintf(`'%s'`, i))
	}

	return strings.Join(lstString, seq)
}

func RemoveDuplicates[T string | int64](lst []T) []T {
	lstMap := make(map[T]bool, len(lst))
	res := make([]T, 0, len(lst))
	for _, i := range lst {
		exists, ok := lstMap[i]
		if ok && exists {
			continue
		}
		lstMap[i] = true
		res = append(res, i)
	}

	return res
}

func GetTime(t time.Time) string {
	return fmt.Sprintf(`'%s'`, t.Format("2006-01-02 15:04:05"))
}

func GetCol(col string) string {
	colLst := strings.Split(col, ".")
	if len(colLst) == 0 {
		return "``"
	}

	for i, c := range colLst {
		if strings.HasPrefix(c, "`") || strings.HasSuffix(c, "`") {
			continue
		}

		colLst[i] = fmt.Sprintf("`%s`", c)
	}

	return strings.Join(colLst, ".")
}

func InList[T string | int64](lst []T, element T) bool {
	for _, i := range lst {
		if i == element {
			return true
		}
	}

	return false
}
