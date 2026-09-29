package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	workOrderModelSelf interface {
		FindOneWithoutDelete(ctx context.Context, id int64) (*WorkOrder, error)
		FindOneByUidWithoutDelete(ctx context.Context, uid string) (*WorkOrder, error)
		UpdateCh(ctx context.Context, data *WorkOrder) error
		GetList(ctx context.Context, userID int64, src string, status []int64, page int64, pageSize int64, startTime, endTime int64, timeType int64, fromID int64) ([]WorkOrder, error)
		GetCount(ctx context.Context, userID int64, src string, status []int64, startTime, endTime int64, timeType int64, fromID int64) (int64, error)
	}
)

const (
	WorkOrderStatusWaitUser  = 1
	WorkOrderStatusWaitReply = 2
	WorkOrderStatusFinish    = 3
)

func IsWorkOrderStatus(status int64) bool {
	return status == WorkOrderStatusWaitUser || status == WorkOrderStatusWaitReply || status == WorkOrderStatusFinish
}

func (m *defaultWorkOrderModel) UpdateCh(ctx context.Context, data *WorkOrder) error {
	err := m.Update(ctx, data)
	if err == nil && !data.DeleteAt.Valid {
		UpdateWorkOrder(data)
	}
	return err
}

func (m *defaultWorkOrderModel) FindOneWithoutDelete(ctx context.Context, id int64) (*WorkOrder, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `delete_at` is null order by id desc limit 1", workOrderRows, m.table)
	var resp WorkOrder
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

func (m *defaultWorkOrderModel) FindOneByUidWithoutDelete(ctx context.Context, uid string) (*WorkOrder, error) {
	query := fmt.Sprintf("select %s from %s where `uid` = ? and `delete_at` is null order by id desc limit 1", workOrderRows, m.table)
	var resp WorkOrder
	err := m.conn.QueryRowCtx(ctx, &resp, query, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWorkOrderModel) GetList(ctx context.Context, userID int64, src string, status []int64, page int64, pageSize int64, startTime, endTime int64, timeType int64, fromID int64) ([]WorkOrder, error) {
	var resp []WorkOrder
	var err error

	cond := where.NewCond(m.table, workOrderFieldNames).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).Int64In("status", status).UserID(userID).WebIDWithCenter(fromID, "from_id").Like(src, true, "title")
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", workOrderRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WorkOrder{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWorkOrderModel) GetCount(ctx context.Context, userID int64, src string, status []int64, startTime, endTime int64, timeType int64, fromID int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, workOrderFieldNames).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).Int64In("status", status).UserID(userID).WebIDWithCenter(fromID, "from_id").Like(src, true, "title")
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
