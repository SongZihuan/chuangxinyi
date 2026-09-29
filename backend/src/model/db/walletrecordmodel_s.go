package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	walletRecordModelSelf interface {
		FindOneWithoutDelete(ctx context.Context, id int64) (*WalletRecord, error)
		GetList(ctx context.Context, walletID int64, t []int64, fundingID string, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]WalletRecord, error)
		GetCount(ctx context.Context, walletID int64, t []int64, fundingID string, src string, startTime, endTime int64, timeType int64) (int64, error)
	}
)

const (
	WalletPay      = 1
	WalletBack     = 2
	WalletDefray   = 3
	WalletInvoice  = 4
	WalletAdmin    = 5
	WalletWithdraw = 7
)

func IsWalletRecordType(t int64) bool {
	switch t {
	case WalletPay, WalletBack, WalletDefray, WalletInvoice, WalletAdmin, WalletWithdraw:
		return true
	default:
		return false
	}
}

func (m *defaultWalletRecordModel) FindOneWithoutDelete(ctx context.Context, id int64) (*WalletRecord, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and delete_at is null order by create_at desc limit 1", walletRecordRows, m.table)
	var resp WalletRecord
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

func (m *defaultWalletRecordModel) GetList(ctx context.Context, walletID int64, t []int64, fundingID string, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]WalletRecord, error) {
	var resp []WalletRecord
	var err error

	cond := where.NewCond(m.table, walletRecordFieldNames).Int64In("type", t).WalletID(walletID).StringEQ("funding_id", fundingID).Like(src, true, "reason").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", walletRecordRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))
	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WalletRecord{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWalletRecordModel) GetCount(ctx context.Context, walletID int64, t []int64, fundingID string, src string, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, walletRecordFieldNames).Int64In("type", t).WalletID(walletID).StringEQ("funding_id", fundingID).Like(src, true, "reason").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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
