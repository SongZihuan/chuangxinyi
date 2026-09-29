package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	websiteFundingModelSelf interface {
		FindOneWithDelete(ctx context.Context, id int64) (*WebsiteFunding, error)
		GetList(ctx context.Context, websiteID int64, t []int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]WebsiteFunding, error)
		GetCount(ctx context.Context, websiteID int64, t []int64, startTime, endTime int64, timeType int64) (int64, error)
		GetYearSum(ctx context.Context, websiteID int64, year int64) ([]YearSumData, error)
	}
)

const (
	WebsiteFundingPay           = 1
	WebsiteFundingPayRefund     = 2
	WebsiteFundingPayRefundFail = 3
	WebsiteFundingDefray        = 4
	WebsiteFundingDefrayReturn  = 5
	WebsiteFundingBack          = 6
	WebsiteFundingWithdraw      = 7
	WebsiteFundingWithdrawFail  = 8
)

func IsWebsiteFundingType(t int64) bool {
	return t == WebsiteFundingPay || t == WebsiteFundingPayRefund || t == WebsiteFundingPayRefundFail || t == WebsiteFundingDefray || t == WebsiteFundingDefrayReturn || t == WebsiteFundingBack || t == WebsiteFundingWithdraw || t == WebsiteFundingWithdrawFail
}

func (m *defaultWebsiteFundingModel) FindOneWithDelete(ctx context.Context, id int64) (*WebsiteFunding, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `delete_at` is null limit 1", websiteFundingRows, m.table)
	var resp WebsiteFunding
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

func (m *defaultWebsiteFundingModel) GetList(ctx context.Context, websiteID int64, t []int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]WebsiteFunding, error) {
	var resp []WebsiteFunding
	var err error

	cond := where.NewCond(m.table, websiteFundingFieldNames).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).WebIDWithCenter(websiteID, "web_id").Int64In("type", t)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", websiteFundingRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WebsiteFunding{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWebsiteFundingModel) GetCount(ctx context.Context, websiteID int64, t []int64, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, websiteFundingFieldNames).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).WebIDWithCenter(websiteID, "web_id").Int64In("type", t)
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

type YearSumData struct {
	Profit      int64 `db:"profit"`
	Expenditure int64 `db:"expenditure"`
	Delta       int64 `db:"delta"`
	Month       int64 `db:"month"`
	Day         int64 `db:"day"`
}

func (m *defaultWebsiteFundingModel) GetYearSum(ctx context.Context, websiteID int64, year int64) ([]YearSumData, error) {
	var resp []YearSumData
	var err error

	query := fmt.Sprintf("select SUM(profit) as profit, SUM(expenditure) as expenditure, 0+res1-res2 as delta, `month` as month, `day` as day from %s where `web_id` = ? and delete_at is null and `year` = ? group by `day`,` month`", websiteFundingRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, websiteID, year)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []YearSumData{}, nil
	default:
		return []YearSumData{}, err
	}
}
