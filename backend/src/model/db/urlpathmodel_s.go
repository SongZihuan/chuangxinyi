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
	urlPathModelSelf interface {
		FindBySignWithoutDelete(ctx context.Context, name string) (*UrlPath, error)
		FindOneWithoutDelete(ctx context.Context, id int64) (*UrlPath, error)
		GetList(ctx context.Context) ([]UrlPath, error)
		GetCount(ctx context.Context) (int64, error)
	}
)

const (
	PathStatusOk     = 1
	PathStatusDelete = 2
	PathStatusBanned = 3
)

const (
	PathModePrefix   = 1 // 前缀匹配
	PathModeComplete = 2 // 完全匹配
	PathModeRegex    = 3 // 正则匹配
)

const (
	PathCorsAll     = 1 // 允许跨域
	PathCorsWebsite = 2 // 允许外站跨域
	PathCorsCenter  = 3 // 限制跨域
)

const (
	PathNotAdmin     = 1 // 非管理员接口
	PathWebsiteAdmin = 2 // 管理员接口，外站允许
	PathCenterAdmin  = 3 // 管理员接口
)

const (
	PathBusyModeIP   = 1 // IP限制模式
	PathBusyModeUser = 2 // 用户限制模式
)

const (
	PathGet = 1 << iota
	PathPost
)

var PathMethodStringMap = map[string]int64{
	http.MethodGet:  PathGet,
	http.MethodPost: PathPost,
}

const (
	CaptchaModeNone        = 1 // 不需要验证码
	CaptchaModeOn          = 2 // 需要验证码
	CaptchaModeSilenceOnly = 3 // 需要验证码，仅限静默
	CaptchaModeSliderOnly  = 4 // 需要验证码，仅限滑块
)

func IsCaptchaMode(captchaMode int64) bool {
	return captchaMode == CaptchaModeNone || captchaMode == CaptchaModeOn || captchaMode == CaptchaModeSilenceOnly || captchaMode == CaptchaModeSliderOnly
}

func IsPathBusyMode(busyMode int64) bool {
	return busyMode == PathBusyModeIP || busyMode == PathBusyModeUser
}

func IsPathAdminMode(adminMode int64) bool {
	return adminMode == PathNotAdmin || adminMode == PathWebsiteAdmin || adminMode == PathCenterAdmin
}

func IsPathCorsModel(coreModel int64) bool {
	return coreModel == PathCorsAll || coreModel == PathCorsWebsite || coreModel == PathCorsCenter
}

func IsPathStatus(pathStatus int64) bool {
	return pathStatus == PathStatusBanned || pathStatus == PathStatusOk || pathStatus == PathStatusDelete
}

func IsPathMode(pathMode int64) bool {
	return pathMode == PathModePrefix || pathMode == PathModeComplete || pathMode == PathModeRegex
}

func (m *defaultUrlPathModel) FindBySignWithoutDelete(ctx context.Context, sign string) (*UrlPath, error) {
	query := fmt.Sprintf("select %s from %s where `sign`=? and `delete_at` is null order by id desc limit 1", urlPathRows, m.table)
	var resp UrlPath
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

func (m *defaultUrlPathModel) FindOneWithoutDelete(ctx context.Context, id int64) (*UrlPath, error) {
	var resp UrlPath
	query := fmt.Sprintf("select %s from %s where `id` = ? order by id desc limit 1", urlPathRows, m.table)
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

func (m *defaultUrlPathModel) GetList(ctx context.Context) ([]UrlPath, error) {
	var resp []UrlPath
	cond := where.NewCond(m.table, urlPathFieldNames)
	query := fmt.Sprintf("select %s from %s where %s order by id %s", urlPathRows, m.table, cond, where.NewLimit(config.BackendConfig.MySQL.SystemResourceLimit*2))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []UrlPath{}, nil
	default:
		return nil, err
	}
}

func (m *defaultUrlPathModel) GetCount(ctx context.Context) (int64, error) {
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
