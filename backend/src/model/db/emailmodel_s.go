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
	emailModelSelf interface {
		FindByEmail(ctx context.Context, email string) (*Email, error)
		FindByUserID(ctx context.Context, userID int64) (*Email, error)
		InsertWithDelete(ctx context.Context, data *Email) (sql.Result, error)
		DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error)
		InsertCh(ctx context.Context, data *Email) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Email) error
	}
)

func (m *defaultEmailModel) InsertCh(ctx context.Context, data *Email) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultEmailModel) UpdateCh(ctx context.Context, data *Email) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultEmailModel) FindByEmail(ctx context.Context, email string) (*Email, error) {
	query := fmt.Sprintf("select %s from %s where `email` = ? and delete_at is null and is_delete = false order by create_at desc limit 1", emailRows, m.table)
	var resp Email
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultEmailModel) FindByUserID(ctx context.Context, userID int64) (*Email, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", emailRows, m.table)
	var resp Email
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

func (m *defaultEmailModel) InsertWithDelete(ctx context.Context, data *Email) (sql.Result, error) {
	if !data.IsDelete && data.Email.Valid {
		key1 := fmt.Sprintf("db:insert:%s:%s", m.table, data.Email.String)
		if !redis.AcquireLockMore(ctx, key1, time.Minute*2) {
			return nil, fmt.Errorf("delete fail")
		}
		defer redis.ReleaseLock(key1)
	}

	key2 := fmt.Sprintf("db:insert:%s:%d", m.table, data.UserId)
	if !redis.AcquireLockMore(ctx, key2, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key2)

	if !data.IsDelete && data.Email.Valid {
		updateQuery1 := fmt.Sprintf("update %s set delete_at=? where email = ? and is_delete = false and delete_at is null", m.table)
		retUpdate1, err := m.conn.ExecCtx(ctx, updateQuery1, time.Now(), data.Email)
		if err != nil {
			return retUpdate1, err
		}
	}

	updateQuery2 := fmt.Sprintf("update %s set delete_at=? where user_id = ? and delete_at is null", m.table)
	retUpdate2, err := m.conn.ExecCtx(ctx, updateQuery2, time.Now(), data.UserId)
	if err != nil {
		return retUpdate2, err
	}

	return m.InsertCh(ctx, data)
}

func (m *defaultEmailModel) DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error) {
	updateQuery1 := fmt.Sprintf("update %s set is_delete=? where user_id = ?", m.table)
	return m.conn.ExecCtx(ctx, updateQuery1, true, userID)
}
