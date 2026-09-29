package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	backModelSelf interface {
		FindByBackID(ctx context.Context, backid string) (*Back, error)
		GetList(ctx context.Context, walletID int64, supplierId int64, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Back, error)
		GetCount(ctx context.Context, walletID int64, supplierId int64, src string, startTime, endTime int64, timeType int64) (int64, error)
	}
)

func (m *defaultBackModel) FindByBackID(ctx context.Context, backid string) (*Back, error) {
	query := fmt.Sprintf("select %s from %s where `back_id` = ? and delete_at is null order by create_at desc limit 1", backRows, m.table)
	var resp Back
	err := m.conn.QueryRowCtx(ctx, &resp, query, backid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultBackModel) GetList(ctx context.Context, walletID int64, supplierId int64, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Back, error) {
	var resp []Back
	var err error

	cond := where.NewCond(m.table, backFieldNames).WalletID(walletID).Like(src, true, "subject").WebIDWithCenter(supplierId, "supplier_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", backRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Back{}, nil
	default:
		return nil, err
	}
}

func (m *defaultBackModel) GetCount(ctx context.Context, walletID int64, supplierId int64, src string, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, backFieldNames).WalletID(walletID).Like(src, true, "subject").WebIDWithCenter(supplierId, "supplier_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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
