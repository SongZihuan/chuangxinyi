package db

import (
	"context"
	"fmt"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	uncleModelSelf interface {
		FindByUserIDWithoutDelete(ctx context.Context, userID int64, uncleID int64) (*Uncle, error)
		GetUncleCount(ctx context.Context, userID int64) (int64, error)
		GetNephewCount(ctx context.Context, uncleID int64) (int64, error)
	}
)

const (
	UncleWaitOk = 1 // 等待确认
	UncleOK     = 2 // 确认
)

func (m *defaultUncleModel) FindByUserIDWithoutDelete(ctx context.Context, userID int64, uncleID int64) (*Uncle, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `uncle_id` = ? and `delete_at` is null order by id desc limit 1", uncleRows, m.table)
	var resp Uncle
	err := m.conn.QueryRowCtx(ctx, &resp, query, userID, uncleID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUncleModel) GetUncleCount(ctx context.Context, userID int64) (int64, error) {
	var err error
	var resp OneInt

	query := fmt.Sprintf("select count(id) as res from %s where delete_at is null and user_id = ?", m.table)
	err = m.conn.QueryRowCtx(ctx, &resp, query, userID)

	switch err {
	case nil:
		return resp.Res, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}

func (m *defaultUncleModel) GetNephewCount(ctx context.Context, uncleID int64) (int64, error) {
	var err error
	var resp OneInt

	query := fmt.Sprintf("select count(id) as res from %s where delete_at is null and uncle_id = ?", m.table)
	err = m.conn.QueryRowCtx(ctx, &resp, query, uncleID)

	switch err {
	case nil:
		return resp.Res, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}
