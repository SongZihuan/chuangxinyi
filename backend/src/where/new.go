package where

import (
	"fmt"
	"reflect"
	"strings"
)

func NewCond(table string, field []string) *Cond {
	return newCond(table, field).NotDeleteAt()
}

func NewCondWithStruct(table string, model any) *Cond {
	return newCond(table, rawFieldNames(model)).NotDeleteAt()
}

func NewCondWithoutDeleteAtWithStruct(table string, model any) *Cond {
	return newCond(table, rawFieldNames(model))
}

func NewCondWithoutDeleteAt(table string, field []string) *Cond {
	return newCond(table, field)
}

func newCond(table string, field []string) *Cond {
	return &Cond{
		table: table,
		cond:  "1=1",
		field: field,
	}
}

func NewPage(page int64, pageSize int64) *Page {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 20
	}

	return &Page{
		page: page,
		size: pageSize,
	}
}

func NewLimit(limit int64) *Limit {
	if limit <= 0 {
		limit = 0
	}

	return &Limit{
		limit: limit,
	}
}

const dbTag = "db"

// RawFieldNames converts golang struct field into slice string.
func rawFieldNames(in any) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		tagv := fi.Tag.Get(dbTag)
		switch tagv {
		case "-":
			continue
		case "":
			out = append(out, fi.Name)
		default:
			// get tag name with the tag opton, e.g.:
			// `db:"id"`
			// `db:"id,type=char,length=16"`
			// `db:",type=char,length=16"`
			// `db:"-,type=char,length=16"`
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if tagv == "-" {
				continue
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}
			out = append(out, tagv)
		}
	}

	return out
}
