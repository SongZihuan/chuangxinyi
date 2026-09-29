package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	tokenRecordModelSelf interface {
		GetList(ctx context.Context, token string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]TokenRecord, error)
		GetCount(ctx context.Context, token string, startTime, endTime int64, timeType int64) (int64, error)
		InsertWithCreate(ctx context.Context, data *TokenRecord) (sql.Result, error)
	}
)

const (
	LoginToken = 1
	UserToken  = 2
)

const (
	TokenCreate      = 1
	TokenGeoIPUpdate = 2
	TokenDelete      = 3
)

func (m *defaultTokenRecordModel) GetList(ctx context.Context, token string, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]TokenRecord, error) {
	var resp []TokenRecord
	var err error

	cond := where.NewCond(m.table, tokenRecordFieldNames).Add(token, "`token`='%s'", token).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", tokenRecordRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []TokenRecord{}, nil
	default:
		return nil, err
	}
}

func (m *defaultTokenRecordModel) GetCount(ctx context.Context, token string, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, tokenRecordFieldNames).Add(token, "`token`='%s'", token).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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

func (m *defaultTokenRecordModel) InsertWithCreate(ctx context.Context, data *TokenRecord) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s, create_at) values (?, ?, ?, ?, ?, ?)", m.table, tokenRecordRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.TokenType, data.Token, data.Type, data.Data, data.DeleteAt, data.CreateAt)
	return ret, err
}
