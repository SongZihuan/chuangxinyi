package utils

import (
	"fmt"
	"reflect"
)

func GetStructName(i interface{}) string {
	return getStructName(i, 1000)
}

func getStructName(i interface{}, limit int64) string {
	if limit == 0 {
		return "unknown"
	}

	if i == nil {
		return "nil"
	}

	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Pointer {
		return getStructName(t.Elem(), limit-1)
	}

	if t.Kind() == reflect.Interface {
		return fmt.Sprintf("%s(%s)", getStructName(t.Elem(), limit-1), t.Name())
	}

	return t.Name()
}
