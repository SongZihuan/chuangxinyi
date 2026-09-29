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
	oauth2BanedModelSelf interface {
		GetList(ctx context.Context, userID int64, limit int64) ([]Oauth2Baned, error)
		CheckAllow(ctx context.Context, userID int64, webID int64, t int64) (bool, error)
		InsertWithDelete(ctx context.Context, data *Oauth2Baned) (sql.Result, error)
	}
)

const (
	AllowLogin  = 1
	AllowDefray = 2
	AllowMsg    = 3
)

func (m *defaultOauth2BanedModel) GetList(ctx context.Context, userID int64, limit int64) ([]Oauth2Baned, error) {
	var resp []Oauth2Baned
	cond := where.NewCond(m.table, oauth2BanedFieldNames).UserID(userID)
	query := fmt.Sprintf("select %s from %s where %s order by id desc %s", oauth2BanedRows, m.table, cond, where.NewLimit(limit))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []Oauth2Baned{}, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultOauth2BanedModel) CheckAllow(ctx context.Context, userID int64, webID int64, t int64) (bool, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `web_id` = ? and delete_at is null order by id desc limit 1", oauth2BanedRows, m.table)
	var resp Oauth2Baned
	err := m.conn.QueryRowCtx(ctx, &resp, query, userID, webID)
	switch err {
	case nil:
		switch t {
		case AllowLogin:
			return resp.AllowLogin, nil
		case AllowDefray:
			return resp.AllowLogin && resp.AllowDefray, nil
		case AllowMsg:
			return resp.AllowLogin && resp.AllowMsg, nil
		default:
			return false, fmt.Errorf("bad type")
		}
	case sqlc.ErrNotFound:
		switch t {
		case AllowLogin:
			return false, nil
		case AllowDefray:
			return false, nil
		case AllowMsg:
			return false, nil
		default:
			return false, fmt.Errorf("bad type")
		}
	default:
		return false, err
	}
}

func (m *defaultOauth2BanedModel) InsertWithDelete(ctx context.Context, data *Oauth2Baned) (sql.Result, error) {
	key := fmt.Sprintf("db:insert:%s:%d:%d", m.table, data.UserId, data.WebId)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return nil, fmt.Errorf("delete fail")
	}
	defer redis.ReleaseLock(key)

	updateQuery := fmt.Sprintf("update %s set delete_at=? where user_id = ? and web_id = ? and delete_at is null", m.table)
	retUpdate, err := m.conn.ExecCtx(ctx, updateQuery, time.Now(), data.UserId, data.WebId)
	if err != nil {
		return retUpdate, err
	}

	return m.Insert(ctx, data)
}
