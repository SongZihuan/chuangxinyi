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
	wechatModelSelf interface {
		FindByUserID(ctx context.Context, userID int64) (*Wechat, error)
		FindByOpenID(ctx context.Context, openID string) (*Wechat, error)
		InsertWithDelete(ctx context.Context, data *Wechat) (sql.Result, error)
		DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error)
		FindByUnionID(ctx context.Context, unionID string) (*Wechat, error)
		InsertCh(ctx context.Context, data *Wechat) (sql.Result, error)
		UpdateCh(ctx context.Context, data *Wechat) error
	}
)

func (m *defaultWechatModel) InsertCh(ctx context.Context, data *Wechat) (sql.Result, error) {
	ret, err := m.Insert(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return ret, err
}

func (m *defaultWechatModel) UpdateCh(ctx context.Context, data *Wechat) error {
	err := m.Update(ctx, data)
	UpdateUser(data.UserId, m.conn, nil)
	return err
}

func (m *defaultWechatModel) FindByUserID(ctx context.Context, userID int64) (*Wechat, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by create_at desc limit 1", wechatRows, m.table)
	var resp Wechat
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

func (m *defaultWechatModel) FindByOpenID(ctx context.Context, openID string) (*Wechat, error) {
	query := fmt.Sprintf("select %s from %s where `open_id` = ? and is_delete = false and delete_at is null order by create_at desc limit 1", wechatRows, m.table)
	var resp Wechat
	err := m.conn.QueryRowCtx(ctx, &resp, query, openID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWechatModel) FindByUnionID(ctx context.Context, unionID string) (*Wechat, error) {
	query := fmt.Sprintf("select %s from %s where `union_id` = ? and is_delete = false and delete_at is null order by create_at desc limit 1", wechatRows, m.table)
	var resp Wechat
	err := m.conn.QueryRowCtx(ctx, &resp, query, unionID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWechatModel) InsertWithDelete(ctx context.Context, data *Wechat) (sql.Result, error) {
	if data.OpenId.Valid {
		key1 := fmt.Sprintf("db:insert:%s:%s", m.table, data.OpenId.String)
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

	updateQuery1 := fmt.Sprintf("update %s set delete_at=? where user_id = ? and delete_at is null", m.table)
	retUpdate1, err := m.conn.ExecCtx(ctx, updateQuery1, time.Now(), data.UserId)
	if err != nil {
		return retUpdate1, err
	}

	updateQuery2 := fmt.Sprintf("update %s set delete_at=? where open_id = ? and is_delete = false and delete_at is null", m.table)
	retUpdate2, err := m.conn.ExecCtx(ctx, updateQuery2, time.Now(), data.OpenId)
	if err != nil {
		return retUpdate2, err
	}

	return m.InsertCh(ctx, data)
}

func (m *defaultWechatModel) DeleteByUserID(ctx context.Context, userID int64) (sql.Result, error) {
	updateQuery := fmt.Sprintf("update %s set is_delete=? where user_id = ?", m.table)
	return m.conn.ExecCtx(ctx, updateQuery, true, userID)
}
