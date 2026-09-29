package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"time"
)

var PhoneRepeat = fmt.Errorf("Phone Repeat")

type (
	phoneModelSelf interface {
		FindByPhone(ctx context.Context, phone string) (*Phone, error)
		FindByUserID(ctx context.Context, userID int64) (*Phone, error)
		InsertSafe(ctx context.Context, data *Phone) (sql.Result, error)
		InsertWithDelete(ctx context.Context, data *Phone) (sql.Result, error)
		DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error)
		InsertCh(ctx context.Context, data *Phone) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Phone) error
	}
)

func (m *defaultPhoneModel) InsertCh(ctx context.Context, data *Phone) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultPhoneModel) UpdateCh(ctx context.Context, data *Phone) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultPhoneModel) FindByPhone(ctx context.Context, phone string) (*Phone, error) {
	query := fmt.Sprintf("select %s from %s where `phone` = ? and delete_at is null and is_delete = false order by create_at desc limit 1", phoneRows, m.table)
	var resp Phone
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPhoneModel) FindByUserID(ctx context.Context, userID int64) (*Phone, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by id desc limit 1", phoneRows, m.table)
	var resp Phone
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

func (m *defaultPhoneModel) InsertSafe(ctx context.Context, data *Phone) (sql.Result, error) {
	if !data.IsDelete {
		key1 := fmt.Sprintf("db:insert:%s:%s", m.table, data.Phone)
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
		query1 := fmt.Sprintf("select COUNT(id) as res from %s where `phone` = ? and is_delete = false and delete_at is null", m.table)
		err := m.conn.QueryRowCtx(ctx, &resp1, query1, data.Phone)
		if !errors.Is(err, ErrNotFound) && err != nil {
			return nil, err
		} else if resp1.Res.Int64 != 0 {
			return nil, PhoneRepeat
		}
	}

	updateQuery2 := fmt.Sprintf("update %s set delete_at=? where user_id = ? and delete_at is null", m.table)
	retUpdate2, err := m.conn.ExecCtx(ctx, updateQuery2, time.Now(), data.UserId)
	if err != nil {
		return retUpdate2, err
	}

	return m.InsertCh(ctx, data)
}

func (m *defaultPhoneModel) InsertWithDelete(ctx context.Context, data *Phone) (sql.Result, error) {
	key1 := fmt.Sprintf("db:insert:%s:%s", m.table, data.Phone)
	if !redis.AcquireLockMore(ctx, key1, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key1)

	key2 := fmt.Sprintf("db:insert:%s:%d", m.table, data.UserId)
	if !redis.AcquireLockMore(ctx, key2, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key2)

	updateQuery1 := fmt.Sprintf("update %s set delete_at=? where user_id = ? and delete_at is null", m.table)
	retUpdate1, err := m.conn.ExecCtx(ctx, updateQuery1, time.Now(), data.UserId)
	if err != nil {
		return retUpdate1, err
	}

	updateQuery2 := fmt.Sprintf("update %s set delete_at=? where phone = ? and delete_at is null", m.table)
	retUpdate2, err := m.conn.ExecCtx(ctx, updateQuery2, time.Now(), data.Phone)
	if err != nil {
		return retUpdate2, err
	}

	return m.InsertCh(ctx, data)
}

func (m *defaultPhoneModel) DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error) {
	updateQuery1 := fmt.Sprintf("update %s set is_delete=? where user_id = ?", m.table)
	return m.conn.ExecCtx(ctx, updateQuery1, true, userID)
}
