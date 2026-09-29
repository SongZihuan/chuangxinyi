package warp

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"math/big"
	"regexp"
)

type UrlPath struct {
	types.UrlPath
	PolicyPermission    big.Int        `json:"-"`
	SubPolicyPermission int64          `json:"-"`
	Regex               *regexp.Regexp `json:"-"`
	MethodPermission    int64          `json:"-"`
}

func (u UrlPath) GetUrlPathType() types.UrlPath {
	return u.UrlPath
}

func (u UrlPath) GetRoleUrlPathType() types.RoleUrlPath {
	return types.RoleUrlPath{
		ID:             u.ID,
		Path:           u.Path,
		Describe:       u.Describe,
		Mode:           u.Mode,
		Status:         u.Status,
		Authentication: u.Authentication,
		DoubleCheck:    u.DoubleCheck,
		CorsMode:       u.CorsMode,
		AdminMode:      u.AdminMode,
		BusyMode:       u.BusyMode,
		BusyCount:      u.BusyCount,
	}
}
