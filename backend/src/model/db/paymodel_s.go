package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	payModelSelf interface {
		FindByPayID(ctx context.Context, payid string) (*Pay, error)
		GetList(ctx context.Context, walletID int64, src string, status []int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Pay, error)
		GetCount(ctx context.Context, walletID int64, src string, status []int64, startTime, endTime int64, timeType int64) (int64, error)
		InsertWithCreateAt(ctx context.Context, data *Pay) (sql.Result, error)
	}
)

const (
	PayWait                = 1 // 等待支付（不可退款，可查询订单）
	PaySuccess             = 2 // 支付成功（可退款）
	PayFinish              = 3 // 支付完成（可退款）
	PayClose               = 4 // 支付关闭（不可退款）
	PayWaitRefund          = 5 // 等待退款（可查询退款）
	PaySuccessRefund       = 6 // 支付退款（退款完成）
	PayCloseRefund         = 7 // 退款失败（可退款）
	PaySuccessRefundInside = 8 // 退款成功（仅限单边退款）
)

func IsPayStatus(status int64) bool {
	return status == PayWait || status == PaySuccess || status == PayFinish || status == PayClose || status == PayWaitRefund || status == PaySuccessRefund || status == PaySuccessRefundInside || status == PayCloseRefund
}

func (m *defaultPayModel) FindByPayID(ctx context.Context, payid string) (*Pay, error) {
	query := fmt.Sprintf("select %s from %s where `pay_id` = ? and delete_at is null order by create_at desc limit 1", payRows, m.table)
	var resp Pay
	err := m.conn.QueryRowCtx(ctx, &resp, query, payid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPayModel) GetList(ctx context.Context, walletID int64, src string, status []int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Pay, error) {
	var resp []Pay
	var err error

	cond := where.NewCond(m.table, payFieldNames).WalletID(walletID).Int64In("trade_status", status).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).Like(src, true, "subject", "pay_way")
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", payRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Pay{}, nil
	default:
		return nil, err
	}
}

func (m *defaultPayModel) GetCount(ctx context.Context, walletID int64, src string, status []int64, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, payFieldNames).WalletID(walletID).Int64In("trade_status", status).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).Like(src, true, "subject", "pay_way")
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

func (m *defaultPayModel) InsertWithCreateAt(ctx context.Context, data *Pay) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s, create_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, payRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Subject, data.PayWay, data.PayId, data.Cny, data.Get, data.CouponsId, data.TradeNo, data.BuyerId, data.TradeStatus, data.Balance, data.PayAt, data.RefundAt, data.DeleteAt, data.CreateAt)
	return ret, err
}
