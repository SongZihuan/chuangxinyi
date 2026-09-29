package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	roleModelSelf interface {
		HaveAny(ctx context.Context) (bool, error)
		FindBySignWithoutDelete(ctx context.Context, name string) (*Role, error)
		FindOneWithoutDelete(ctx context.Context, id int64) (*Role, error)
		GetList(ctx context.Context) ([]Role, error)
		GetCount(ctx context.Context) (int64, error)
	}
)

const (
	RoleStatusOK     = 1
	RoleStatusBanned = 2
)

func IsRoleStatus(roleStatus int64) bool {
	return roleStatus == RoleStatusOK || roleStatus == RoleStatusBanned
}

func (m *defaultRoleModel) HaveAny(ctx context.Context) (bool, error) {
	query := fmt.Sprintf("select %s from %s where delete_at is null order by id desc limit 1", roleRows, m.table)
	var resp Role
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return true, nil
	case sqlc.ErrNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (m *defaultRoleModel) FindBySignWithoutDelete(ctx context.Context, sign string) (*Role, error) {
	query := fmt.Sprintf("select %s from %s where `sign`=? and `delete_at` is null order by id desc limit 1", roleRows, m.table)
	var resp Role
	err := m.conn.QueryRowCtx(ctx, &resp, query, sign)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultRoleModel) FindOneWithoutDelete(ctx context.Context, id int64) (*Role, error) {
	var resp Role
	query := fmt.Sprintf("select %s from %s where `id` = ? order by id desc limit 1", roleRows, m.table)
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

func (m *defaultRoleModel) GetList(ctx context.Context) ([]Role, error) {
	var resp []Role
	cond := where.NewCond(m.table, roleFieldNames)
	query := fmt.Sprintf("select %s from %s where %s order by id %s", roleRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SystemResourceLimit*2))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Role{}, nil
	default:
		return nil, err
	}
}

func (m *defaultRoleModel) GetCount(ctx context.Context) (int64, error) {
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
