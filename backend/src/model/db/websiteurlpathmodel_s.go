package db

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/where"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlc"
	"net/http"
)

type (
	websiteUrlPathModelSelf interface {
		FindBySignWithoutDelete(ctx context.Context, name string) (*WebsiteUrlPath, error)
		FindOneWithoutDelete(ctx context.Context, id int64) (*WebsiteUrlPath, error)
		GetList(ctx context.Context) ([]WebsiteUrlPath, error)
		GetCount(ctx context.Context) (int64, error)
	}
)

const (
	WebsitePathStatusOk     = 1
	WebsitePathStatusDelete = 2
	WebsitePathStatusBanned = 3
)

const (
	WebsitePathModePrefix   = 1 // 前缀匹配
	WebsitePathModeComplete = 2 // 完整匹配
	WebsitePathModeRegex    = 3 // 正则匹配
)

const (
	WebsitePathGet = 1 << iota
	WebsitePathPost
)

var WebsitePathMethodStringMap = map[string]int64{
	http.MethodGet:  WebsitePathGet,
	http.MethodPost: WebsitePathPost,
}

func IsWebsitePathStatus(pathStatus int64) bool {
	return pathStatus == WebsitePathStatusBanned || pathStatus == WebsitePathStatusOk || pathStatus == WebsitePathStatusDelete
}

func IsWebsitePathMode(pathMode int64) bool {
	return pathMode == WebsitePathModePrefix || pathMode == WebsitePathModeComplete || pathMode == WebsitePathModeRegex
}

func (m *defaultWebsiteUrlPathModel) FindBySignWithoutDelete(ctx context.Context, sign string) (*WebsiteUrlPath, error) {
	query := fmt.Sprintf("select %s from %s where `sign`=? and `delete_at` is null order by id desc limit 1", websiteUrlPathRows, m.table)
	var resp WebsiteUrlPath
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

func (m *defaultWebsiteUrlPathModel) FindOneWithoutDelete(ctx context.Context, id int64) (*WebsiteUrlPath, error) {
	var resp WebsiteUrlPath
	query := fmt.Sprintf("select %s from %s where `id` = ? order by id desc limit 1", websiteUrlPathRows, m.table)
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

func (m *defaultWebsiteUrlPathModel) GetList(ctx context.Context) ([]WebsiteUrlPath, error) {
	var resp []WebsiteUrlPath
	cond := where.NewCond(m.table, websiteUrlPathFieldNames)
	query := fmt.Sprintf("select %s from %s where %s order by id %s", websiteUrlPathRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SystemResourceLimit*2))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []WebsiteUrlPath{}, nil
	default:
		return nil, err
	}
}

func (m *defaultWebsiteUrlPathModel) GetCount(ctx context.Context) (int64, error) {
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
