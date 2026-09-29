package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	websiteIpModelSelf interface {
		GetListByWebsiteIDWithoutDelete(ctx context.Context, websiteID int64) ([]WebsiteIp, error)
		GetList(ctx context.Context, websiteID int64) ([]WebsiteIp, error)
		GetCount(ctx context.Context) (int64, error)
		FindOneWithoutDelete(ctx context.Context, id int64) (*WebsiteIp, error)
	}
)

func (m *defaultWebsiteIpModel) GetListByWebsiteIDWithoutDelete(ctx context.Context, websiteID int64) ([]WebsiteIp, error) {
	var resp []WebsiteIp
	query := fmt.Sprintf("select %s from %s where `website_id` = ? and delete_at is null order by id desc", websiteIpRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, websiteID)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WebsiteIp{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWebsiteIpModel) GetList(ctx context.Context, websiteID int64) ([]WebsiteIp, error) {
	var resp []WebsiteIp
	cond := where.NewCond(m.table, websiteIpFieldNames).LinkID(websiteID, "website_id")
	query := fmt.Sprintf("select %s from %s where %s order by id %s", websiteIpRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SystemResourceLimit*2))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WebsiteIp{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWebsiteIpModel) GetCount(ctx context.Context) (int64, error) {
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

func (m *defaultWebsiteIpModel) FindOneWithoutDelete(ctx context.Context, id int64) (*WebsiteIp, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and delete_at is null order by id desc limit 1", websiteIpRows, m.table)
	var resp WebsiteIp
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
