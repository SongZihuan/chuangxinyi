package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	invoiceModelSelf interface {
		FindByInvoiceID(ctx context.Context, invoiceID string) (*Invoice, error)
		GetList(ctx context.Context, walletID int64, t []int64, status []int64, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Invoice, error)
		GetCount(ctx context.Context, walletID int64, t []int64, status []int64, src string, startTime, endTime int64, timeType int64) (int64, error)
	}
)

const (
	PersonalInvoice           = 1 // 个人普票
	CompanyInvoice            = 2 // 企业普票
	CompanySpecializedInvoice = 3 // 企业专票
)

const (
	InvoiceWait       = 1 // 待开票
	InvoiceOK         = 2 // 已开票
	InvoiceReturn     = 3 // 已退票
	InvoiceBad        = 4 // 信息错误
	InvoiceRedFlush   = 5 // 已红冲
	InvoiceWaitReturn = 6 // 等待退票
)

func IsInvoiceType(t int64) bool {
	return t == PersonalInvoice || t == CompanyInvoice || t == CompanySpecializedInvoice
}

func IsInvoiceStatus(s int64) bool {
	return s == InvoiceWait || s == InvoiceOK || s == InvoiceReturn || s == InvoiceBad || s == InvoiceRedFlush || s == InvoiceWaitReturn
}

func (m *defaultInvoiceModel) FindByInvoiceID(ctx context.Context, invoiceID string) (*Invoice, error) {
	query := fmt.Sprintf("select %s from %s where `invoice_id` = ? and delete_at is null order by create_at desc limit 1", invoiceRows, m.table)
	var resp Invoice
	err := m.conn.QueryRowCtx(ctx, &resp, query, invoiceID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultInvoiceModel) GetList(ctx context.Context, walletID int64, t []int64, status []int64, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Invoice, error) {
	var resp []Invoice
	var err error

	cond := where.NewCond(m.table, invoiceFieldNames).Int64In("type", t).Int64In("status", status).WalletID(walletID).Like(src, true, "name", "tax_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", invoiceRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Invoice{}, nil
	default:
		return nil, err
	}
}

func (m *defaultInvoiceModel) GetCount(ctx context.Context, walletID int64, t []int64, status []int64, src string, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, invoiceFieldNames).Int64In("type", t).Int64In("status", status).WalletID(walletID).Like(src, true, "name", "tax_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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
