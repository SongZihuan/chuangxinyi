package utils

import (
	"fmt"
	errors "github.com/wuntsong-org/wterrors"
	"strconv"
	"strings"
)

func RemoveTrailingNil[T any](lst []T, isZero func(T) bool) []T {
	index := len(lst) - 1
	for ; index >= 0; index-- {
		if !isZero(lst[index]) {
			return lst[:index+1]
		}
	}
	return nil
}

func InList[T string | int64](lst []T, element T) bool {
	for _, i := range lst {
		if i == element {
			return true
		}
	}

	return false
}

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

func JoinStringToString(lst []string, seq string, removeDuplicates bool) string {
	if len(lst) == 0 {
		return ""
	}

	if removeDuplicates {
		lst = RemoveDuplicates(lst)
	}

	return strings.Join(lst, seq)
}

func ListConvert[T1, T2 string | int64](lst []T1, m map[T1]T2) []T2 {
	res := make([]T2, 0, len(lst))
	for _, i := range lst {
		n, ok := m[i]
		if !ok {
			continue
		}
		res = append(res, n)
	}

	return res
}

func ListConvertMust[T1, T2 string | int64](lst []T1, m map[T1]T2) ([]T2, errors.WTError) {
	res := make([]T2, 0, len(lst))
	for _, i := range lst {
		n, ok := m[i]
		if !ok {
			return nil, errors.Errorf("convert fail")
		}
		res = append(res, n)
	}

	return res, nil
}

func SplitStringToInt(lst string, seq string) []int64 {
	if len(lst) == 0 {
		return make([]int64, 0)
	}

	stringLst := strings.Split(lst, seq)
	res := make([]int64, 0, len(stringLst))
	for _, s := range stringLst {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		res = append(res, n)
	}

	return res
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
