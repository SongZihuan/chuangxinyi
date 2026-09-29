package db

import (
	"context"
	"fmt"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	footerModelSelf interface {
		FindTheNew(ctx context.Context) (*Footer, error)
	}
)

func (m *defaultFooterModel) FindTheNew(ctx context.Context) (*Footer, error) {
	query := fmt.Sprintf("select %s from %s order by create_at desc limit 1", footerRows, m.table)
	var resp Footer
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
