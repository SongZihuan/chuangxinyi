package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	menuModelSelf interface {
		HaveAny(ctx context.Context) (bool, error)
		FindOneWithoutDelete(ctx context.Context, id int64) (*Menu, error)
		GetList(ctx context.Context) ([]Menu, error)
		GetNewSortNumber(ctx context.Context, fatherID int64) (res int64, err error)
		FindNear(ctx context.Context, fatherID int64, sort int64, isUp bool) (res *Menu, err error)
		GetSonList(ctx context.Context, fatherID int64, limit int64) ([]Menu, error)
		GetCount(ctx context.Context) (int64, error)
	}
)

const (
	MenuStatusOK     = 1
	MenuStatusBanned = 2
)

func IsMenuStatus(menuStatus int64) bool {
	return menuStatus == MenuStatusOK || menuStatus == MenuStatusBanned
}

func (m *defaultMenuModel) HaveAny(ctx context.Context) (bool, error) {
	query := fmt.Sprintf("select %s from %s where delete_at is null order by id desc limit 1", menuRows, m.table)
	var resp Menu
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

func (m *defaultMenuModel) FindOneWithoutDelete(ctx context.Context, id int64) (*Menu, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `delete_at` is null order by id desc limit 1", menuRows, m.table)
	var resp Menu
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

func (m *defaultMenuModel) GetSonList(ctx context.Context, fatherID int64, limit int64) ([]Menu, error) {
	var resp []Menu
	if limit == 0 || limit > config.BackendConfig.MySQL.SystemResourceLimit {
		limit = config.BackendConfig.MySQL.SystemResourceLimit * 2
	}

	cond := where.NewCond(m.table, applicationFieldNames).HasFatherID(fatherID)
	query := fmt.Sprintf("select %s from %s where %s order by id desc %s", menuRows, m.table, cond, where.NewLimit(limit))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Menu{}, nil
	default:
		return nil, err
	}
}

func (m *defaultMenuModel) GetList(ctx context.Context) ([]Menu, error) {
	var resp []Menu
	cond := where.NewCond(m.table, applicationFieldNames)
	query := fmt.Sprintf("select %s from %s where %s order by sort %s", menuRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SystemResourceLimit*2))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Menu{}, nil
	default:
		return nil, err
	}
}

func (m *defaultMenuModel) GetCount(ctx context.Context) (int64, error) {
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

func (m *defaultMenuModel) GetNewSortNumber(ctx context.Context, fatherID int64) (res int64, err error) {
	var resp OneIntOrNull
	if fatherID == 0 {
		query := fmt.Sprintf("select max(sort) as res from %s where father_id is null and delete_at is null", m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query)
	} else {
		query := fmt.Sprintf("select max(sort) as res from %s where father_id = ? and delete_at is null", m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, fatherID)
	}
	switch err {
	case nil:
		return resp.Res.Int64 + 1, nil
	case sqlc.ErrNotFound:
		return 1, nil
	default:
		return 0, err
	}
}

func (m *defaultMenuModel) FindNear(ctx context.Context, fatherID int64, sort int64, isUp bool) (res *Menu, err error) {
	var resp Menu
	if isUp {
		if fatherID == 0 {
			query := fmt.Sprintf("select %s from %s where father_id is null and `sort` < ? and delete_at is null order by `sort` desc limit 1", menuRows, m.table)
			err = m.conn.QueryRowCtx(ctx, &resp, query, sort)
		} else {
			query := fmt.Sprintf("select %s from %s where father_id = ? and `sort` < ? and delete_at is null order by `sort` desc limit 1", menuRows, m.table)
			err = m.conn.QueryRowCtx(ctx, &resp, query, fatherID, sort)
		}
	} else {
		if fatherID == 0 {
			query := fmt.Sprintf("select %s from %s where father_id is null and `sort` > ? and delete_at is null order by `sort` asc limit 1", menuRows, m.table)
			err = m.conn.QueryRowCtx(ctx, &resp, query, sort)
		} else {
			query := fmt.Sprintf("select %s from %s where father_id = ? and `sort` > ? and delete_at is null order by `sort` asc limit 1", menuRows, m.table)
			err = m.conn.QueryRowCtx(ctx, &resp, query, fatherID, sort)
		}
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
