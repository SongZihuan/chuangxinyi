package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"time"
)

type (
	nicknameModelSelf interface {
		FindByUserID(ctx context.Context, userID int64) (*Nickname, error)
		InsertWithDelete(ctx context.Context, data *Nickname) (sql.Result, error)
		InsertCh(ctx context.Context, data *Nickname) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Nickname) error
	}
)

func (m *defaultNicknameModel) InsertCh(ctx context.Context, data *Nickname) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultNicknameModel) UpdateCh(ctx context.Context, data *Nickname) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultNicknameModel) FindByUserID(ctx context.Context, userID int64) (*Nickname, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", nicknameRows, m.table)
	var resp Nickname
	err := m.conn.QueryRowCtx(ctx, &resp, query, userID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNicknameModel) InsertWithDelete(ctx context.Context, data *Nickname) (sql.Result, error) {
	key := fmt.Sprintf("db:insert:%s:%d", m.table, data.UserId)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key)

	updateQuery := fmt.Sprintf("update %s set delete_at=? where user_id = ? and delete_at is null", m.table)
	retUpdate, err := m.conn.ExecCtx(ctx, updateQuery, time.Now(), data.UserId)
	if err != nil {
		return retUpdate, err
	}

	return m.InsertCh(ctx, data)
}
