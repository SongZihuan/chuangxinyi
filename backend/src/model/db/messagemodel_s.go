package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"time"
)

type (
	messageModelSelf interface {
		FindOneWithoutDelete(ctx context.Context, id int64) (*Message, error)
		ReadAll(ctx context.Context, userID int64) error
		UpdateReadCh(ctx context.Context, data *Message) error
		InsertCh(ctx context.Context, data *Message) (sql.Result, error)
		GetList(ctx context.Context, userID int64, src string, justNotRead bool, page int64, pageSize int64, startTime, endTime int64, timeType int64, senderID int64) ([]Message, error)
		GetCount(ctx context.Context, userID int64, src string, justNotRead bool, startTime, endTime int64, timeType int64, senderID int64) (int64, error)
		ReadAllWebsite(ctx context.Context, userID int64, senderID int64) error
	}
)

func (m *defaultMessageModel) InsertCh(ctx context.Context, data *Message) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	if err == nil {
		id, err := ret.LastInsertId()
		if err == nil {
			data.Id = id
			UpdateMessage(data)
		}
	}
	return ret, err
}

func (m *defaultMessageModel) UpdateReadCh(ctx context.Context, data *Message) error {
	err := m.Update(ctx, data)
	if err == nil && !data.DeleteAt.Valid {
		ReadMessage(data)
	}
	return err
}

func (m *defaultMessageModel) FindOneWithoutDelete(ctx context.Context, id int64) (*Message, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `delete_at` is null order by id desc limit 1", messageRows, m.table)
	var resp Message
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

func (m *defaultMessageModel) ReadAll(ctx context.Context, userID int64) error {
	query := fmt.Sprintf("update %s set read_at = ? where `user_id` = ? and delete_at is null", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), userID)
	return err
}

func (m *defaultMessageModel) ReadAllWebsite(ctx context.Context, userID int64, senderID int64) error {
	query := fmt.Sprintf("update %s set read_at = ? where `user_id` = ? and `sender_id` = ? and delete_at is null", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), userID, senderID)
	return err
}

func (m *defaultMessageModel) GetList(ctx context.Context, userID int64, src string, justNotRead bool, page int64, pageSize int64, startTime, endTime int64, timeType int64, senderID int64) ([]Message, error) {
	var resp []Message
	var err error

	cond := where.NewCond(m.table, messageFieldNames).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).UserID(userID).WebIDWithCenter(senderID, "sender_id").Like(src, true, "title")
	if justNotRead {
		cond = cond.Add("read_at is null")
	}
	query := fmt.Sprintf("select %s from %s where %s order by %s %s", messageRows, m.table, cond, cond.OrderBy(), where.NewPage(page, pageSize))
	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Message{}, nil
	default:
		return nil, err
	}
}

func (m *defaultMessageModel) GetCount(ctx context.Context, userID int64, src string, justNotRead bool, startTime, endTime int64, timeType int64, senderID int64) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, messageFieldNames).TimeBetweenWithTimeMap(timeType, timeMap, startTime, endTime).UserID(userID).WebIDWithCenter(senderID, "sender_id").Like(src, true, "title")
	if justNotRead {
		cond = cond.Add("read_at is null")
	}
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
