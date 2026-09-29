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
	headerModelSelf interface {
		FindByHeader(ctx context.Context, header string) (*Header, error)
		FindByHeaderWithoutDelete(ctx context.Context, header string) (*Header, error)
		FindByUserID(ctx context.Context, userID int64) (*Header, error)
		InsertWithDelete(ctx context.Context, data *Header) (sql.Result, error)
	}
)

func (m *defaultHeaderModel) InsertCh(ctx context.Context, data *Header) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultHeaderModel) UpdateCh(ctx context.Context, data *Header) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultHeaderModel) FindByHeader(ctx context.Context, header string) (*Header, error) {
	// 不需要delete_at is null
	query := fmt.Sprintf("select %s from %s where `header` = ? order by create_at desc limit 1", headerRows, m.table)
	var resp Header
	err := m.conn.QueryRowCtx(ctx, &resp, query, header)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultHeaderModel) FindByHeaderWithoutDelete(ctx context.Context, header string) (*Header, error) {
	// 不需要delete_at is null
	query := fmt.Sprintf("select %s from %s where `header` = ? and `delete_at` is null order by create_at desc limit 1", headerRows, m.table)
	var resp Header
	err := m.conn.QueryRowCtx(ctx, &resp, query, header)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultHeaderModel) FindByUserID(ctx context.Context, userID int64) (*Header, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", headerRows, m.table)
	var resp Header
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

func (m *defaultHeaderModel) InsertWithDelete(ctx context.Context, data *Header) (sql.Result, error) {
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
