package db

import (
	"context"
	"fmt"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	walletModelSelf interface {
		FindByWalletID(ctx context.Context, userID int64) (*Wallet, error)
	}
)

func (m *defaultWalletModel) FindByWalletID(ctx context.Context, walletID int64) (*Wallet, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and delete_at is null order by create_at desc limit 1", walletRows, m.table)
	var resp Wallet
	err := m.conn.QueryRowCtx(ctx, &resp, query, walletID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
