package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	oauth2RecordModelSelf interface {
		GetList(ctx context.Context, userID int64, webID int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Oauth2Record, error)
		GetCount(ctx context.Context, userID int64, webID int64, startTime, endTime int64, timeType int64) (int64, error)
	}
)

func (m *defaultOauth2RecordModel) GetList(ctx context.Context, userID int64, webID int64, page int64, pageSize int64, startTime, endTime int64, timeType int64) ([]Oauth2Record, error) {
	var resp []Oauth2Record
	var err error

	cond := where.NewCond(m.table, oauth2RecordFieldNames).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).UserID(userID).WebIDWithoutCenter(webID, "web_id")
	query := fmt.Sprintf("select %s from %s where %s order by `login_time` desc %s", oauth2RecordRows, m.table, cond, where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Oauth2Record{}, nil
	default:
		return nil, err
	}
}

func (m *defaultOauth2RecordModel) GetCount(ctx context.Context, userID int64, webID int64, startTime, endTime int64, timeType int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, oauth2RecordFieldNames).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).UserID(userID).WebIDWithoutCenter(webID, "web_id")
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
