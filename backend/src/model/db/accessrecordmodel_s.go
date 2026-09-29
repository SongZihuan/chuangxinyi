package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"strings"
)

type (
	accessRecordModelSelf interface {
		GetCount(ctx context.Context, token string, startTime, endTime int64, timeType int64) (int64, error)
		GetList(ctx context.Context, token string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]AccessRecord, error)
		GetListByCond(ctx context.Context, cond string, page int64, pageSize int64) ([]AccessRecord, string, error)
		GetCountByCond(ctx context.Context, cond string) (int64, string, error)
	}
)

func (m *defaultAccessRecordModel) GetList(ctx context.Context, token string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]AccessRecord, error) {
	var resp []AccessRecord
	var err error

	cond := where.NewCond(m.table, accessRecordFieldNames).Add(token, "`user_token`='%s'", token).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", accessRecordRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []AccessRecord{}, nil
	default:
		return nil, err
	}
}

func (m *defaultAccessRecordModel) GetCount(ctx context.Context, token string, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, accessRecordFieldNames).Add(token, "`user_token`='%s'", token).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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

func (m *defaultAccessRecordModel) GetListByCond(ctx context.Context, cond string, page int64, pageSize int64) ([]AccessRecord, string, error) {
	var resp []AccessRecord
	var err error

	if len(strings.TrimSpace(cond)) != 0 {
		cond = fmt.Sprintf("WHERE (%s)", cond)
	}

	ctx = context.WithValue(ctx, "Allow-Table-Name", m.table)
	ctx = context.WithValue(ctx, "Allow-Col-Name", accessRecordFieldNames)
	ctx = context.WithValue(ctx, "Allow-Func-Name", where.SafeSqlFunc)

	query := fmt.Sprintf("select %s from %s %s order by id desc limit %d offset %d", accessRecordRows, m.table, cond, pageSize, (page-1)*pageSize)
	safe, r, err := where.CheckSQL(ctx, query)
	if err != nil {
		return nil, query, err
	} else if !safe {
		return nil, query, fmt.Errorf(r)
	}

	err = m.conn.QueryRowsCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp, query, nil
	case sqlc.ErrNotFound:
		return []AccessRecord{}, query, nil
	default:
		return nil, query, err
	}
}

func (m *defaultAccessRecordModel) GetCountByCond(ctx context.Context, cond string) (int64, string, error) {
	var err error
	var resp OneInt

	if len(strings.TrimSpace(cond)) != 0 {
		cond = fmt.Sprintf("WHERE (%s)", cond)
	}

	ctx = context.WithValue(ctx, "Allow-Table-Name", m.table)
	ctx = context.WithValue(ctx, "Allow-Col-Name", accessRecordFieldNames)
	ctx = context.WithValue(ctx, "Allow-Func-Name", where.SafeSqlFunc)

	query := fmt.Sprintf("select count(id) as res from %s %s", m.table, cond)
	safe, r, err := where.CheckSQL(ctx, query)
	if err != nil {
		return 0, query, err
	} else if !safe {
		return 0, query, fmt.Errorf(r)
	}

	err = m.conn.QueryRowCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp.Res, query, nil
	case sqlc.ErrNotFound:
		return 0, query, nil
	default:
		return 0, query, err
	}
}
