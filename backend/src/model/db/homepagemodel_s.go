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
	homepageModelSelf interface {
		FindByUserID(ctx context.Context, userID int64) (*Homepage, error)
		InsertWithDelete(ctx context.Context, data *Homepage) (sql.Result, error)
		InsertCh(ctx context.Context, data *Homepage) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Homepage) error
	}
)

func (m *defaultHomepageModel) InsertCh(ctx context.Context, data *Homepage) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultHomepageModel) UpdateCh(ctx context.Context, data *Homepage) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultHomepageModel) FindByUserID(ctx context.Context, userID int64) (*Homepage, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", homepageRows, m.table)
	var resp Homepage
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

func (m *defaultHomepageModel) InsertWithDelete(ctx context.Context, data *Homepage) (sql.Result, error) {
	key := fmt.Sprintf("db:insert:%s:%d", m.table, data.UserId)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key)

	updateQuery1 := fmt.Sprintf("update %s set delete_at=? where user_id = ? and delete_at is null", m.table)
	retUpdate1, err := m.conn.ExecCtx(ctx, updateQuery1, time.Now(), data.UserId)
	if err != nil {
		return retUpdate1, err
	}

	return m.InsertCh(ctx, data)
}
