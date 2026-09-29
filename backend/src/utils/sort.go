/*
使用go排序的例子
此处只是举个例子
具体使用需要复制代码并修改类型

因为golang的泛型等问题导致只能这么做
*/
package utils

import "sort"

type LessFunction func(any, any) bool

type Sorter struct {
	Data   []any
	LessFn LessFunction
}

func (s Sorter) Len() int {
	return len(s.Data)
}

func (s Sorter) Less(i, j int) bool {
	return s.LessFn(s.Data[i], s.Data[j])
}

func (s Sorter) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

func SortList(lst []any, less LessFunction) {
	sort.Sort(Sorter{
		Data:   lst,
		LessFn: less,
	})
}
