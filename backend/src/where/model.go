package where

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"strings"
	"time"
)

type Cond struct {
	cond    string
	field   []string
	table   string
	orderBy string
}

type Page struct {
	page int64
	size int64
}

type Limit struct {
	limit int64
}

func (c *Cond) TimeBetweenWithTimeMap(colNum int64, m map[int64]string, start, end int64) *Cond {
	if colNum == 0 || m == nil {
		return c
	}

	col, ok := m[colNum]
	if !ok || len(col) == 0 {
		return c
	}

	if !c.IsColumn(col) {
		return c
	}

	return c.TimeBetween(col, start, end)
}

func (c *Cond) TimeBetween(col string, start, end int64) *Cond {
	if len(col) == 0 {
		return c
	}

	if !c.IsColumn(col) {
		return c
	}

	if start > end {
		start, end = end, start
	}

	if start == 0 && end == 0 {
		return c
	}

	if start == 0 {
		c.cond += fmt.Sprintf(" and %s < %s", GetCol(col), GetTime(time.Unix(end, 0)))
		return c
	}

	if end == 0 {
		c.cond += fmt.Sprintf(" and %s > %s", GetCol(col), GetTime(time.Unix(start, 0)))
		return c
	}

	c.cond += fmt.Sprintf(" and %s between %s and %s", GetCol(col), GetTime(time.Unix(start, 0)), GetTime(time.Unix(end, 0)))
	c.orderBy = col
	return c
}

func (c *Cond) Int64In(col string, num []int64) *Cond {
	if len(col) == 0 {
		return c
	}

	if !c.IsColumn(col) {
		return c
	}

	if num == nil || len(num) == 0 {
		return c
	}

	c.cond += fmt.Sprintf(" and %s in (%s)", GetCol(col), JoinIntToString(num, ",", true))
	return c
}

func (c *Cond) StringIn(col string, str []string) *Cond {
	if len(col) == 0 {
		return c
	}

	if !c.IsColumn(col) {
		return c
	}

	if str == nil || len(str) == 0 {
		return c
	}

	c.cond += fmt.Sprintf(" and %s in (%s)", GetCol(col), JoinStringWithQuotaToString(str, ",", true))
	return c
}

func (c *Cond) Like(src string, doublePre bool, col1 string, col ...string) *Cond {
	if len(src) == 0 {
		return c
	}

	if len(col) == 0 {
		return c
	}

	var like string
	if doublePre {
		like = "%%" + src + "%%" // 使用%%，满足转义
	} else {
		like = "%" + src + "%" // 使用%%，满足转义
	}

	base := "0=1"
	if c.IsColumn(col1) {
		base = fmt.Sprintf("%s like '%s'", GetCol(col1), like)
	}

	for _, v := range col {
		if len(v) == 0 {
			continue
		}

		if !c.IsColumn(v) {
			return c
		}

		base += fmt.Sprintf(" or %s like '%s'", GetCol(v), like)
	}

	c.cond += fmt.Sprintf(" and (%s)", base)
	return c
}

func (c *Cond) Add(cond string, args ...any) *Cond {
	if len(cond) == 0 {
		return c
	}

	c.cond += fmt.Sprintf(" and (%s)", fmt.Sprintf(cond, args...))
	return c
}

func (c *Cond) AddIfExists(prob string, cond string, args ...any) *Cond {
	if len(cond) == 0 || len(prob) == 0 {
		return c
	}

	return c.Add(cond, args...)
}

func (c *Cond) NotDeleteAt() *Cond {
	if !c.IsColumn("delete_at") {
		return c
	}

	c.cond += " and `delete_at` is null"
	return c
}

func (c *Cond) UserID(userID int64) *Cond {
	if userID == 0 {
		return c
	}

	if !c.IsColumn("user_id") {
		return c
	}

	c.cond += fmt.Sprintf(" and `user_id` = %d", userID)
	return c
}

func (c *Cond) StringEQ(col, str string) *Cond {
	if len(str) == 0 || len(col) == 0 {
		return c
	}

	if !c.IsColumn(col) {
		return c
	}

	c.cond += fmt.Sprintf(" and %s = '%s'", GetCol(col), str)
	return c
}

func (c *Cond) HasFatherID(fatherID int64) *Cond {
	if !c.IsColumn("father_id") {
		return c
	}

	if fatherID == 0 {
		c.cond += fmt.Sprintf(" and `father_id` is null")
		return c
	}

	c.cond += fmt.Sprintf(" and `father_id` = %d", fatherID)
	return c
}

func (c *Cond) OwnerID(ownerID int64) *Cond {
	if !c.IsColumn("owner_id") {
		return c
	}

	if ownerID == 0 {
		return c
	}

	c.cond += fmt.Sprintf(" and `owner_id` = %d", ownerID)
	return c
}

func (c *Cond) WalletID(walletID int64) *Cond {
	if !c.IsColumn("wallet_id") {
		return c
	}

	if walletID == 0 {
		return c
	}

	c.cond += fmt.Sprintf(" and `wallet_id` = %d", walletID)
	return c
}

func (c *Cond) WebIDWithoutCenter(webID int64, col string) *Cond {
	if webID == 0 {
		return c
	}

	if !c.IsColumn(col) {
		return c
	}

	c.cond += fmt.Sprintf(" and %s = %d", GetCol(col), webID)
	return c
}

func (c *Cond) WebIDWithCenter(webID int64, col string) *Cond {
	if webID < 0 {
		return c
	}

	if !c.IsColumn(col) {
		return c
	}

	c.cond += fmt.Sprintf(" and %s = %d", GetCol(col), webID)
	return c
}

func (c *Cond) NotUserID() *Cond {
	if !c.IsColumn("user_id") {
		return c
	}

	c.cond += fmt.Sprintf(" and `user_id` is null")
	return c
}

func (c *Cond) HasUserID(userID int64) *Cond {
	if !c.IsColumn("user_id") {
		return c
	}

	if userID == 0 {
		c.cond += fmt.Sprintf(" and `user_id` is not null")
	}

	return c.UserID(userID)
}

func (c *Cond) LinkUUIDAndType(uuid string, uuidCol string, t int64, typeCol string) *Cond {
	if len(uuid) == 0 || t == 0 || len(uuidCol) == 0 || len(typeCol) == 0 {
		return c
	}

	if !c.IsColumn(uuidCol) {
		return c
	}

	if !c.IsColumn(typeCol) {
		return c
	}

	c.cond += fmt.Sprintf(" and %s = '%s' and %s = %d", GetCol(uuidCol), uuid, GetCol(typeCol), t)
	return c
}

func (c *Cond) LinkID(id int64, idCol string) *Cond {
	if id == 0 || len(idCol) == 0 {
		return c
	}

	if !c.IsColumn(idCol) {
		return c
	}

	c.cond += fmt.Sprintf(" and %s = %d", GetCol(idCol), id)
	return c
}

func (c *Cond) IsColumn(col string) bool {
	colLst := strings.Split(col, ".")
	if len(colLst) == 0 {
		return false
	} else if len(colLst) == 1 {
		return InList(c.field, col) || InList(c.field, fmt.Sprintf("`%s`", col))
	}

	col = colLst[len(colLst)-1]
	return InList(c.field, col) || InList(c.field, fmt.Sprintf("`%s`", col))
}

func (c *Cond) OrderBy() string {
	if c.orderBy == "" {
		if c.IsColumn("create_at") {
			return "`create_at` desc"
		} else {
			return "`id` desc"
		}
	}

	return fmt.Sprintf("%s desc", GetCol(c.orderBy))
}

func (c *Cond) OrderByAsc() string {
	if c.orderBy == "" {
		if c.IsColumn("create_at") {
			return "`create_at` asc"
		} else {
			return "`id` asc"
		}
	}

	return fmt.Sprintf("%s asc", GetCol(c.orderBy))
}

func (c *Cond) String() string {
	ctx := context.WithValue(context.Background(), "Allow-Func-Name", SafeSqlFunc)
	ctx = context.WithValue(ctx, "Allow-Table-Name", c.table)
	ctx = context.WithValue(ctx, "Allow-Col-Name", c.field)

	s := fmt.Sprintf("select 1 from %s where %s", c.table, c.cond)
	safe, msg, err := CheckSQL(ctx, s)
	if err != nil {
		logger.Logger.Error("SQL Error: %s, %s", err.Error(), s)
		return "0=1"
	} else if !safe {
		logger.Logger.Error("SQL NotSafe: %s, %s", msg, s)
		return "0=1"
	}

	return c.cond
}

func (p *Page) String() string {
	return fmt.Sprintf("limit %d offset %d", p.size, (p.page-1)*p.size)
}

func (l *Limit) String() string {
	return fmt.Sprintf("limit %d offset %d", l.limit, 0)
}
