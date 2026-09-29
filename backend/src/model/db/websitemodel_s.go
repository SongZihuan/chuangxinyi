package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	websiteModelSelf interface {
		FindOneByUIDWithoutDelete(ctx context.Context, uid string) (*Website, error)
		FindOneWithoutDelete(ctx context.Context, id int64) (*Website, error)
		GetList(ctx context.Context) ([]Website, error)
		UpdateByPermission(ctx context.Context, newPermission string) error
		GetCount(ctx context.Context) (int64, error)
	}
)

const (
	WebsiteStatusOK     = 1
	WebsiteStatusBanned = 2
)

func IsWebsiteStatus(websiteStatus int64) bool {
	return websiteStatus == WebsitePolicyStatusOK || websiteStatus == WebsiteStatusBanned
}

func (m *defaultWebsiteModel) FindOneByUIDWithoutDelete(ctx context.Context, uid string) (*Website, error) {
	var resp Website
	query := fmt.Sprintf("select %s from %s where `uid` = ? and delete_at is null order by id desc limit 1", websiteRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, uid)
	switch err {
	case nil:
		if resp.DeleteAt.Valid {
			return nil, ErrNotFound
		}
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWebsiteModel) FindOneWithoutDelete(ctx context.Context, id int64) (*Website, error) {
	var resp Website
	query := fmt.Sprintf("select %s from %s where `id` = ? and delete_at is null order by id desc limit 1", websiteRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		if resp.DeleteAt.Valid {
			return nil, ErrNotFound
		}
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWebsiteModel) GetList(ctx context.Context) ([]Website, error) {
	var resp []Website
	cond := where.NewCond(m.table, websiteFieldNames)
	query := fmt.Sprintf("select %s from %s where %s order by create_at %s", websiteRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SystemResourceLimit*2))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Website{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWebsiteModel) GetCount(ctx context.Context) (int64, error) {
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

func (m *defaultWebsiteModel) UpdateByPermission(ctx context.Context, newPermission string) error {
	query := fmt.Sprintf("update %s set `permission`=? where `permission` != '0' and `permission` != ''", m.table)
	_, err := m.conn.ExecCtx(ctx, query, newPermission)
	return err
}
