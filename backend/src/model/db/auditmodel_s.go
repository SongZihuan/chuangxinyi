package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	auditModelSelf interface {
		GetList(ctx context.Context, userID int64, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64, fromID int64) ([]Audit, error)
		GetCount(ctx context.Context, userID int64, src string, startTime, endTime int64, timeType int64, fromID int64) (int64, error)
	}
)

func (m *customAuditModel) GetList(ctx context.Context, userID int64, src string, page int64, pageSize int64, startTime, endTime int64, timeType int64, fromID int64) ([]Audit, error) {
	var resp []Audit
	var err error

	cond := where.NewCond(m.table, auditFieldNames).UserID(userID).Like(src, true, "content", "from").WebIDWithCenter(fromID, "from_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", auditRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Audit{}, nil
	default:
		return nil, err
	}
}

func (m *customAuditModel) GetCount(ctx context.Context, userID int64, src string, startTime, endTime int64, timeType int64, fromID int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, auditFieldNames).UserID(userID).Like(src, true, "content").WebIDWithCenter(fromID, "from_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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
