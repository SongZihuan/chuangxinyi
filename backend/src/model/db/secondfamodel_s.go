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
	secondfaModelSelf interface {
		FindByUserID(ctx context.Context, userID int64) (*Secondfa, error)
		InsertWithDelete(ctx context.Context, data *Secondfa) (sql.Result, error)
		InsertCh(ctx context.Context, data *Secondfa) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Secondfa) error
	}
)

func (m *defaultSecondfaModel) InsertCh(ctx context.Context, data *Secondfa) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultSecondfaModel) UpdateCh(ctx context.Context, data *Secondfa) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultSecondfaModel) FindByUserID(ctx context.Context, userID int64) (*Secondfa, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", secondfaRows, m.table)
	var resp Secondfa
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

func (m *defaultSecondfaModel) InsertWithDelete(ctx context.Context, data *Secondfa) (sql.Result, error) {
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
