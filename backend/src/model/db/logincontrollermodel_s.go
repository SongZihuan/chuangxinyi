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
	loginControllerModelSelf interface {
		FindByUserID(ctx context.Context, userID int64) (*LoginController, error)
		InsertWithDelete(ctx context.Context, data *LoginController) (sql.Result, error)
	}
)

func makeDefault(userID int64) *LoginController {
	return &LoginController{
		UserId:        userID,
		AllowPhone:    true,
		AllowPassword: true,
		AllowWechat:   true,
		AllowEmail:    false,
		Allow2Fa:      true,
		CreateAt:      time.Now(),
	}
}

func (m *defaultLoginControllerModel) FindByUserID(ctx context.Context, userID int64) (*LoginController, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and delete_at is null order by id desc limit 1", loginControllerRows, m.table)
	var resp LoginController
	err := m.conn.QueryRowCtx(ctx, &resp, query, userID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return makeDefault(userID), nil
	default:
		return nil, err
	}
}

func (m *defaultLoginControllerModel) InsertWithDelete(ctx context.Context, data *LoginController) (sql.Result, error) {
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

	return m.Insert(ctx, data)
}
