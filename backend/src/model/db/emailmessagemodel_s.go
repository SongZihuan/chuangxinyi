package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	emailMessageModelSelf interface {
		GetList(ctx context.Context, email string, page int64, pageSize int64, startTime, endTime int64, timeType int64, senderID int64) ([]EmailMessage, error)
		GetCount(ctx context.Context, email string, startTime, endTime int64, timeType int64, senderID int64) (int64, error)
	}
)

func (m *defaultEmailMessageModel) GetList(ctx context.Context, email string, page int64, pageSize int64, startTime, endTime int64, timeType int64, senderID int64) ([]EmailMessage, error) {
	var resp []EmailMessage
	var err error

	cond := where.NewCond(m.table, emailFieldNames).StringEQ("email", email).WebIDWithCenter(senderID, "sender_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", emailMessageRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []EmailMessage{}, nil
	default:
		return nil, err
	}
}

func (m *defaultEmailMessageModel) GetCount(ctx context.Context, email string, startTime, endTime int64, timeType int64, senderID int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, emailFieldNames).StringEQ("email", email).WebIDWithCenter(senderID, "sender_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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
