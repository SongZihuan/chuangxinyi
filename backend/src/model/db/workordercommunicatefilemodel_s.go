package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	workOrderCommunicateFileModelSelf interface {
		FindByFidWithoutDelete(ctx context.Context, fid string) (*WorkOrderCommunicateFile, error)
		FindByKeyWithoutDelete(ctx context.Context, key string) (*WorkOrderCommunicateFile, error)
		GetList(ctx context.Context, communicateID int64) ([]WorkOrderCommunicateFile, error)
	}
)

func (m *defaultWorkOrderCommunicateFileModel) GetList(ctx context.Context, communicateID int64) ([]WorkOrderCommunicateFile, error) {
	var resp []WorkOrderCommunicateFile
	var err error

	cond := where.NewCond(m.table, workOrderCommunicateFileFieldNames).LinkID(communicateID, "communicate_id")
	query := fmt.Sprintf("select %s from %s where %s order by `id` %s", workOrderCommunicateFileRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.WorkOrderFileLimit*2))
	err = m.conn.QueryRowsCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WorkOrderCommunicateFile{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWorkOrderCommunicateFileModel) FindByFidWithoutDelete(ctx context.Context, fid string) (*WorkOrderCommunicateFile, error) {
	query := fmt.Sprintf("select %s from %s where fid = ? and delete_at is null order by create_at desc limit 1", workOrderCommunicateFileRows, m.table)
	var resp WorkOrderCommunicateFile
	err := m.conn.QueryRowCtx(ctx, &resp, query, fid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWorkOrderCommunicateFileModel) FindByKeyWithoutDelete(ctx context.Context, key string) (*WorkOrderCommunicateFile, error) {
	query := fmt.Sprintf("select %s from %s where `key` = ? and delete_at is null order by create_at desc limit 1", workOrderCommunicateFileRows, m.table)
	var resp WorkOrderCommunicateFile
	err := m.conn.QueryRowCtx(ctx, &resp, query, key)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
