package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	smsMessageModelSelf interface {
		GetList(ctx context.Context, phone string, page int64, pageSize int64, startTime, endTime int64, timeType int64, senderID int64) ([]SmsMessage, error)
		GetCount(ctx context.Context, phone string, startTime, endTime int64, timeType int64, senderID int64) (int64, error)
	}
)

func (m *defaultSmsMessageModel) GetList(ctx context.Context, phone string, page int64, pageSize int64, startTime, endTime int64, timeType int64, senderID int64) ([]SmsMessage, error) {
	var resp []SmsMessage
	var err error

	cond := where.NewCond(m.table, smsMessageFieldNames).StringEQ("phone", phone).WebIDWithCenter(senderID, "sender_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", smsMessageRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []SmsMessage{}, nil
	default:
		return nil, err
	}
}

func (m *defaultSmsMessageModel) GetCount(ctx context.Context, phone string, startTime, endTime int64, timeType int64, senderID int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, smsMessageFieldNames).StringEQ("phone", phone).WebIDWithCenter(senderID, "sender_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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
