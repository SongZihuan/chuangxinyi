package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"time"
)

type (
	ossFileModelSelf interface {
		FindByFidWithoutDelete(ctx context.Context, fid string) (*OssFile, error)
		GetList(ctx context.Context, fid string, mediaType []string, page int64, pageSize int64) ([]OssFile, error)
		GetCount(ctx context.Context, fid string, mediaType []string) (int64, error)
		InsertWithDelete(ctx context.Context, data *OssFile) (sql.Result, error)
		DeleteAll(ctx context.Context, fid string) (sql.Result, error)
	}
)

func (m *defaultOssFileModel) FindByFidWithoutDelete(ctx context.Context, fid string) (*OssFile, error) {
	query := fmt.Sprintf("select %s from %s where fid = ? and delete_at is null order by create_at desc limit 1", ossFileRows, m.table)
	var resp OssFile
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

func (m *defaultOssFileModel) GetList(ctx context.Context, fid string, mediaType []string, page int64, pageSize int64) ([]OssFile, error) {
	var resp []OssFile
	var err error

	cond := where.NewCond(m.table, ossFileFieldNames).Like(fid, true, "fid").StringIn("media_type", mediaType)
	query := fmt.Sprintf("select %s from %s where %s order by `id` %s", ossFileRows, m.table, cond, where.NewPage(page, pageSize))

	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []OssFile{}, nil
	default:
		return nil, err
	}
}

func (m *defaultOssFileModel) GetCount(ctx context.Context, fid string, mediaType []string) (int64, error) {
	var err error
	var resp OneInt

	cond := where.NewCond(m.table, ossFileFieldNames).Like(fid, true, "fid").StringIn("media_type", mediaType)
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

func (m *defaultOssFileModel) InsertWithDelete(ctx context.Context, data *OssFile) (sql.Result, error) {
	key := fmt.Sprintf("db:insert:%s:%s", m.table, data.Fid)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key)

	updateQuery := fmt.Sprintf("update %s set delete_at=? where `fid` = ? and delete_at is null", m.table)
	retUpdate, err := m.conn.ExecCtx(ctx, updateQuery, time.Now(), data.Fid)
	if err != nil {
		return retUpdate, err
	}

	return m.Insert(ctx, data)
}

func (m *defaultOssFileModel) DeleteAll(ctx context.Context, fid string) (sql.Result, error) {
	updateQuery := fmt.Sprintf("update %s set delete_at=? where `fid` = ? and delete_at is null", m.table)
	return m.conn.ExecCtx(ctx, updateQuery, time.Now(), fid)
}
