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
	idcardModelSelf interface {
		FindByIDCardWithoutCompany(ctx context.Context, idcard string) (*Idcard, error)
		FindByUserID(ctx context.Context, userID int64) (*Idcard, error)
		InsertWithDelete(ctx context.Context, data *Idcard) (sql.Result, error)
		DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error)
		InsertCh(ctx context.Context, data *Idcard) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Idcard) error
	}
)

func (m *defaultIdcardModel) InsertCh(ctx context.Context, data *Idcard) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultIdcardModel) UpdateCh(ctx context.Context, data *Idcard) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultIdcardModel) FindByIDCardWithoutCompany(ctx context.Context, idcard string) (*Idcard, error) {
	query := fmt.Sprintf("select %s from %s where `user_id_card` = ? and is_company = false and is_delete = false and delete_at is null order by create_at desc limit 1", idcardRows, m.table)
	var resp Idcard
	err := m.conn.QueryRowCtx(ctx, &resp, query, idcard)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultIdcardModel) FindByUserID(ctx context.Context, userID int64) (*Idcard, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", idcardRows, m.table)
	var resp Idcard
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

func (m *defaultIdcardModel) InsertWithDelete(ctx context.Context, data *Idcard) (sql.Result, error) {
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

func (m *defaultIdcardModel) DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error) {
	updateQuery := fmt.Sprintf("update %s set is_delete=? where user_id = ?", m.table)
	return m.conn.ExecCtx(ctx, updateQuery, true, userID)
}
