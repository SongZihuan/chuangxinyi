package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	withdrawModelSelf interface {
		FindByWithdrawID(ctx context.Context, withdrawID string) (*Withdraw, error)
		GetList(ctx context.Context, walletID int64, status []int64, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Withdraw, error)
		GetCount(ctx context.Context, walletID int64, status []int64, src string, startTime, endTime int64, timeType int64) (int64, error)
	}
)

const (
	WithdrawWait = 1
	WithdrawOK   = 2
	WithdrawFail = 3
)

func IsWithdrawStatus(status int64) bool {
	return status == WithdrawWait || status == WithdrawOK || status == WithdrawFail
}

func (m *defaultWithdrawModel) FindByWithdrawID(ctx context.Context, withdrawID string) (*Withdraw, error) {
	query := fmt.Sprintf("select %s from %s where `withdraw_id` = ? and delete_at is null order by create_at desc limit 1", withdrawRows, m.table)
	var resp Withdraw
	err := m.conn.QueryRowCtx(ctx, &resp, query, withdrawID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWithdrawModel) GetList(ctx context.Context, walletID int64, status []int64, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Withdraw, error) {
	var resp []Withdraw
	var err error

	cond := where.NewCond(m.table, withdrawFieldNames).Int64In("status", status).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).WalletID(walletID).Like(src, true, "withdraw_way", "name")
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", withdrawRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Withdraw{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWithdrawModel) GetCount(ctx context.Context, walletID int64, status []int64, src string, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, withdrawFieldNames).Int64In("status", status).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).WalletID(walletID).Like(src, true, "withdraw_way", "name")
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
