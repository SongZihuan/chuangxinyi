package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	discountModelSelf interface {
		FindOneWithoutDelete(ctx context.Context, id int64) (*Discount, error)
		GetList(ctx context.Context, src string, page int64, pageSize int64, show bool) ([]Discount, error)
		GetCount(ctx context.Context, src string, show bool) (int64, error)
	}
)

func (m *defaultDiscountModel) FindOneWithoutDelete(ctx context.Context, id int64) (*Discount, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and delete_at is null order by id desc limit 1", discountRows, m.table)
	var resp Discount
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

func (m *defaultDiscountModel) GetList(ctx context.Context, src string, page int64, pageSize int64, show bool) ([]Discount, error) {
	var resp []Discount
	var err error

	cond := where.NewCond(m.table, discountBuyFieldNames).Like(src, true, "name", "short_describe")
	if show {
		cond = cond.Add("`show` = true")
	}
	query := fmt.Sprintf("select %s from %s where %s order by id desc %s", discountRows, m.table, cond, where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Discount{}, nil
	default:
		return nil, err
	}
}

func (m *defaultDiscountModel) GetCount(ctx context.Context, src string, show bool) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, discountBuyFieldNames).Like(src, true, "name", "short_describe")
	if show {
		cond = cond.Add("`show` = true")
	}
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
