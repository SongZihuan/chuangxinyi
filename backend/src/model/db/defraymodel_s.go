package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"time"
)

type (
	defrayModelSelf interface {
		FindByDefrayID(ctx context.Context, defrayID string) (*Defray, error)
		GetList(ctx context.Context, walletID int64, status []int64, src string, supplierID int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Defray, error)
		GetCount(ctx context.Context, walletID int64, status []int64, src string, supplierID int64, startTime, endTime int64, timeType int64) (int64, error)
		GetListWithOwnerID(ctx context.Context, ownerID int64, status []int64, src string, supplierID int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Defray, error)
		GetCountWithOwnerID(ctx context.Context, ownerID int64, status []int64, src string, supplierID int64, startTime, endTime int64, timeType int64) (int64, error)
		GetWaitDistribution(ctx context.Context, limit int64) ([]Defray, error)
	}
)

const (
	DefrayWait       = 1 // 订单等待支付
	DefraySuccess    = 2 // 订单支付完成
	DefrayWaitReturn = 3 // 订单等待退款
	DefrayReturn     = 4 // 订单退款
)

func IsDefrayStatus(status int64) bool {
	return status == DefrayWait || status == DefraySuccess || status == DefrayWaitReturn || status == DefrayReturn
}

func (m *defaultDefrayModel) FindByDefrayID(ctx context.Context, defrayID string) (*Defray, error) {
	query := fmt.Sprintf("select %s from %s where `defray_id` = ? and delete_at is null order by create_at desc limit 1", defrayRows, m.table)
	var resp Defray
	err := m.conn.QueryRowCtx(ctx, &resp, query, defrayID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultDefrayModel) GetListWithOwnerID(ctx context.Context, ownerID int64, status []int64, src string, supplierID int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Defray, error) {
	var resp []Defray
	var err error

	cond := where.NewCond(m.table, defrayFieldNames).Int64In("status", status).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).OwnerID(ownerID).WebIDWithCenter(supplierID, "supplier_id").Like(src, true, "subject")
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", defrayRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Defray{}, nil
	default:
		return nil, err
	}
}

func (m *defaultDefrayModel) GetCountWithOwnerID(ctx context.Context, ownerID int64, status []int64, src string, supplierID int64, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, defrayFieldNames).Int64In("status", status).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).OwnerID(ownerID).WebIDWithCenter(supplierID, "supplier_id").Like(src, true, "subject")
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

func (m *defaultDefrayModel) GetList(ctx context.Context, walletID int64, status []int64, src string, supplierID int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Defray, error) {
	var resp []Defray
	var err error

	cond := where.NewCond(m.table, defrayFieldNames).Int64In("status", status).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).WalletID(walletID).WebIDWithCenter(supplierID, "supplier_id").Like(src, true, "subject")
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", defrayRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Defray{}, nil
	default:
		return nil, err
	}
}

func (m *defaultDefrayModel) GetCount(ctx context.Context, walletID int64, status []int64, src string, supplierID int64, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, defrayFieldNames).Int64In("status", status).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).WalletID(walletID).WebIDWithCenter(supplierID, "supplier_id").Like(src, true, "subject")
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

func (m *defaultDefrayModel) GetWaitDistribution(ctx context.Context, limit int64) ([]Defray, error) {
	var resp []Defray
	var err error

	cond := where.NewCond(m.table, defrayFieldNames).Add("`last_return_at` < %s", where.GetTime(time.Now())).Add("`return_at` is null").Add("`has_distribution` = false").Add("`user_id` is not null")
	query := fmt.Sprintf("select %s from %s where %s order by id desc %s", defrayRows, m.table, cond, where.NewLimit(limit))
	err = m.conn.QueryRowsCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Defray{}, nil
	default:
		return nil, err
	}
}
