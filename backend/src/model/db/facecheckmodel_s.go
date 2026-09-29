package db

import (
	"context"
	"fmt"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
)

type (
	faceCheckModelSelf interface {
		FindOneByCheckID(ctx context.Context, checkID string) (*FaceCheck, error)
		FindOneByCertifyID(ctx context.Context, certifyID string) (*FaceCheck, error)
	}
)

const (
	FaceCheckOK   = 1
	FaceCheckFail = 2
	FaceCheckWait = 3
)

func (m *defaultFaceCheckModel) FindOneByCheckID(ctx context.Context, checkID string) (*FaceCheck, error) {
	query := fmt.Sprintf("select %s from %s where `check_id` = ? and `delete_at` is null order by id desc limit 1", faceCheckRows, m.table)
	var resp FaceCheck
	err := m.conn.QueryRowCtx(ctx, &resp, query, checkID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFaceCheckModel) FindOneByCertifyID(ctx context.Context, certifyID string) (*FaceCheck, error) {
	query := fmt.Sprintf("select %s from %s where `certify_id` = ? and `delete_at` is null order by id desc limit 1", faceCheckRows, m.table)
	var resp FaceCheck
	err := m.conn.QueryRowCtx(ctx, &resp, query, certifyID)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
