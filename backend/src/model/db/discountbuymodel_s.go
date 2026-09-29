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
	discountBuyModelSelf interface {
		FindOneByUserID(ctx context.Context, userID int64, discountID int64) (*DiscountBuy, error)
		InsertWithDelete(ctx context.Context, data *DiscountBuy) (sql.Result, error)
	}
)

func (m *defaultDiscountBuyModel) FindOneByUserID(ctx context.Context, userID int64, discountID int64) (*DiscountBuy, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `discount_id` = ? and delete_at is null order by create_at desc limit 1", discountBuyRows, m.table)
	var resp DiscountBuy
	err := m.conn.QueryRowCtx(ctx, &resp, query, userID, discountID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultDiscountBuyModel) InsertWithDelete(ctx context.Context, data *DiscountBuy) (sql.Result, error) {
	key := fmt.Sprintf("db:insert:%s:%d:%d", m.table, data.UserId, data.DiscountId)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key)

	update := fmt.Sprintf("update %s set `delete_at` = ? where `user_id` = ? and `discount_id` = ? and delete_at is null", m.table)
	_, err := m.conn.ExecCtx(ctx, update, time.Now(), data.UserId, data.DiscountId)
	if err != nil {
		return nil, err
	}

	return m.Insert(ctx, data)
}
