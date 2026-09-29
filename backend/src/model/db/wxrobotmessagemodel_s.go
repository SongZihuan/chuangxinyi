package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	wxrobotMessageModelSelf interface {
		GetList(ctx context.Context, webhook string, page int64, pageSize int64, startTime, endTime int64, timeType int64, senderID int64) ([]WxrobotMessage, error)
		GetCount(ctx context.Context, webhook string, startTime, endTime int64, timeType int64, senderID int64) (int64, error)
	}
)

func (m *defaultWxrobotMessageModel) GetList(ctx context.Context, webhook string, page int64, pageSize int64, startTime, endTime int64, timeType int64, senderID int64) ([]WxrobotMessage, error) {
	var resp []WxrobotMessage
	var err error

	cond := where.NewCond(m.table, wxrobotFieldNames).StringEQ("webhook", webhook).WebIDWithCenter(senderID, "sender_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", wxrobotMessageRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WxrobotMessage{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWxrobotMessageModel) GetCount(ctx context.Context, webhook string, startTime, endTime int64, timeType int64, senderID int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, wxrobotFieldNames).StringEQ("webhook", webhook).WebIDWithCenter(senderID, "sender_id").TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime)
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
