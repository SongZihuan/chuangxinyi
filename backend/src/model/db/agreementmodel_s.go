package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	agreementModelSelf interface {
		FindOneByAid(ctx context.Context, aid string) (*Agreement, error)
		GetList(ctx context.Context, page int64, pageSize int64) ([]Agreement, error)
		GetCount(ctx context.Context) (int64, error)
	}
)

func (m *defaultAgreementModel) FindOneByAid(ctx context.Context, aid string) (*Agreement, error) {
	query := fmt.Sprintf("select %s from %s where `aid` = ? order by id desc limit 1", agreementRows, m.table)
	var resp Agreement
	err := m.conn.QueryRowCtx(ctx, &resp, query, aid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAgreementModel) GetList(ctx context.Context, page int64, pageSize int64) ([]Agreement, error) {
	var resp []Agreement

	cond := where.NewCond(m.table, auditFieldNames)
	query := fmt.Sprintf("select %s from %s where %s order by `id` %s", agreementRows, m.table, cond, where.NewPage(page, pageSize))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Agreement{}, nil
	default:
		return nil, err
	}
}

func (m *defaultAgreementModel) GetCount(ctx context.Context) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, auditFieldNames)
	query := fmt.Sprintf("select count(id) as res from %s where %s", m.table, cond)
	err = m.conn.QueryRowCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp.Res, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}
