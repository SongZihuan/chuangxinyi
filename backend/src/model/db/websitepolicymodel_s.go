package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	websitePolicyModelSelf interface {
		FindOneWithoutDelete(ctx context.Context, id int64) (*WebsitePolicy, error)
		GetList(ctx context.Context) ([]WebsitePolicy, error)
		GetCount(ctx context.Context) (int64, error)
		GetNewSortNumber(ctx context.Context) (res int64, err error)
		FindNear(ctx context.Context, sort int64, isUp bool) (res *WebsitePolicy, err error)
	}
)

const (
	WebsitePolicyStatusOK     = 1
	WebsitePolicyStatusBanned = 2
)

func IsWebsitePolicyStatus(policyStatus int64) bool {
	return policyStatus == WebsitePolicyStatusOK || policyStatus == WebsitePolicyStatusBanned
}

func (m *defaultWebsitePolicyModel) FindOneWithoutDelete(ctx context.Context, id int64) (*WebsitePolicy, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `delete_at` is null order by id desc limit 1", websitePolicyRows, m.table)
	var resp WebsitePolicy
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWebsitePolicyModel) GetList(ctx context.Context) ([]WebsitePolicy, error) {
	var resp []WebsitePolicy
	cond := where.NewCond(m.table, websitePolicyFieldNames)
	query := fmt.Sprintf("select %s from %s where %s order by sort %s", websitePolicyRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SystemResourceLimit*2))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WebsitePolicy{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWebsitePolicyModel) GetCount(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("select COUNT(id) as res from %s where delete_at is null", m.table)
	var resp OneInt
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp.Res, nil
	case sqlc.ErrNotFound:
		return 0, nil
	default:
		return 0, err
	}
}

func (m *defaultWebsitePolicyModel) GetNewSortNumber(ctx context.Context) (res int64, err error) {
	var resp OneIntOrNull
	query := fmt.Sprintf("select max(sort) as res from %s where delete_at is null", m.table)
	err = m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp.Res.Int64 + 1, nil
	case sqlc.ErrNotFound:
		return 1, nil
	default:
		return 0, err
	}
}

func (m *defaultWebsitePolicyModel) FindNear(ctx context.Context, sort int64, isUp bool) (res *WebsitePolicy, err error) {
	var resp WebsitePolicy
	if isUp {
		query := fmt.Sprintf("select %s from %s where `sort` < ? and delete_at is null order by `sort` desc limit 1", websitePolicyRows, m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, sort)
	} else {
		query := fmt.Sprintf("select %s from %s where `sort` > ? and delete_at is null order by `sort` asc limit 1", websitePolicyRows, m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, sort)
	}

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
