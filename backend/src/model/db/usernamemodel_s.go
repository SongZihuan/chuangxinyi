package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"time"
)

var UserNameRepeat = fmt.Errorf("UserName Repeat")

type (
	usernameModelSelf interface {
		FindByUsername(ctx context.Context, username string) (*Username, error)
		FindByUsernameWithBase64(ctx context.Context, username string) (*Username, error)
		FindByUserID(ctx context.Context, userID int64) (*Username, error)
		InsertSafe(ctx context.Context, data *Username) (sql.Result, error)
		DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error)
		InsertCh(ctx context.Context, data *Username) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Username) error
	}
)

func (m *defaultUsernameModel) InsertCh(ctx context.Context, data *Username) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultUsernameModel) UpdateCh(ctx context.Context, data *Username) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultUsernameModel) FindByUsername(ctx context.Context, username string) (*Username, error) {
	query := fmt.Sprintf("select %s from %s where `username` = ? and delete_at is null and is_delete = false order by create_at desc limit 1", usernameRows, m.table)
	var resp Username
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUsernameModel) FindByUsernameWithBase64(ctx context.Context, username string) (*Username, error) {
	return m.FindByUsername(ctx, base64.StdEncoding.EncodeToString([]byte(username)))
}

func (m *defaultUsernameModel) FindByUserID(ctx context.Context, userID int64) (*Username, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", usernameRows, m.table)
	var resp Username
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

func (m *defaultUsernameModel) InsertSafe(ctx context.Context, data *Username) (sql.Result, error) {
	if !data.IsDelete {
		key1 := fmt.Sprintf("db:insert:%s:%s", m.table, data.Username)
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

	if !data.IsDelete {
		var resp1 OneIntOrNull
		query1 := fmt.Sprintf("select COUNT(id) as res from %s where `username` = ? and delete_at is null and is_delete = false", m.table)
		err := m.conn.QueryRowCtx(ctx, &resp1, query1, data.Username)
		if !errors.Is(err, ErrNotFound) && err != nil {
			return nil, err
		} else if resp1.Res.Int64 != 0 {
			return nil, UserNameRepeat
		}
	}

	updateQuery2 := fmt.Sprintf("update %s set delete_at=? where user_id = ? and delete_at is null", m.table)
	retUpdate2, err := m.conn.ExecCtx(ctx, updateQuery2, time.Now(), data.UserId)
	if err != nil {
		return retUpdate2, err
	}

	return m.InsertCh(ctx, data)
}

func (m *defaultUsernameModel) DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error) {
	updateQuery2 := fmt.Sprintf("update %s set is_delete=? where user_id = ?", m.table)
	return m.conn.ExecCtx(ctx, updateQuery2, true, userID)
}
