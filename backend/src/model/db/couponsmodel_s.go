package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	couponsModelSelf interface {
		FindOneWithoutDelete(ctx context.Context, id int64) (*Coupons, error)
		GetList(ctx context.Context, userID int64, t []int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Coupons, error)
		GetCount(ctx context.Context, userID int64, t []int64, startTime, endTime int64, timeType int64) (int64, error)
	}
)

func (m *defaultCouponsModel) FindOneWithoutDelete(ctx context.Context, id int64) (*Coupons, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and delete_at is null order by id desc limit 1", couponsRows, m.table)
	var resp Coupons
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCouponsModel) GetList(ctx context.Context, userID int64, t []int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Coupons, error) {
	var resp []Coupons
	var err error

	cond := where.NewCond(m.table, couponsFieldNames).UserID(userID).Int64In("type", t).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", couponsRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Coupons{}, nil
	default:
		return nil, err
	}
}

func (m *defaultCouponsModel) GetCount(ctx context.Context, userID int64, t []int64, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, couponsFieldNames).UserID(userID).Int64In("type", t).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select count(id) as res from %s where %s", m.table, cond)

	err = m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp.Res, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}
