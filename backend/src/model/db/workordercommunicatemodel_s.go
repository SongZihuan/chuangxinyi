package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	workOrderCommunicateModelSelf interface {
		FindOneWithoutDelete(ctx context.Context, id int64) (*WorkOrderCommunicate, error)
		GetList(ctx context.Context, orderID int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]WorkOrderCommunicate, error)
		GetCount(ctx context.Context, orderID int64, startTime, endTime int64, timeType int64) (int64, error)
	}
)

const (
	CommunicateFromUser  = 1
	CommunicateFromAdmin = 2
)

func (m *defaultWorkOrderCommunicateModel) FindOneWithoutDelete(ctx context.Context, id int64) (*WorkOrderCommunicate, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `delete_at` is null order by id desc limit 1", workOrderCommunicateRows, m.table)
	var resp WorkOrderCommunicate
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

func (m *defaultWorkOrderCommunicateModel) GetList(ctx context.Context, orderID int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]WorkOrderCommunicate, error) {
	var resp []WorkOrderCommunicate
	var err error

	cond := where.NewCond(m.table, workOrderFieldNames).Add("`order_id`=%d", orderID).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", workOrderCommunicateRows, m.table, cond, cond.OrderByAsc(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WorkOrderCommunicate{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWorkOrderCommunicateModel) GetCount(ctx context.Context, orderID int64, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, workOrderFieldNames).Add("`order_id`=%d", orderID).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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
