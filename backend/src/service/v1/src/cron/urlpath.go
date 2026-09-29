package cron

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"math/big"
	"regexp"
)

func UrlPathCron(allP big.Int, roles map[int64]warp.Role, permissionLst []warp.Permission) map[int64]warp.UrlPath {
	defer utils.Recover(logger.Logger, nil, "")

	urlModel := db.NewUrlPathModel(mysql.MySQLConn)
	urlList, err := urlModel.GetList(context.Background())
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return make(map[int64]warp.UrlPath)
	}

	urlPath := make(map[int64]warp.UrlPath, len(urlList))

	for _, u := range urlList {
		var p big.Int
		_, ok := p.SetString(u.Permission, 16)
		if !ok {
			continue
		}
		p = permission.ClearPermission(allP, p)

		roleList := make([]types.MenuRole, 0, len(roles))
		for _, r := range roles {
			if r.IsBanned() {
				continue
			}
			if r.Sign == config.BackendConfig.Admin.RootRole.RoleSign {
				roleList = append(roleList, r.GetMenuRoleTypes())
			} else if u.IsOrPolicy {
				if permission.HasOnePermission(r.PolicyPermissions, p) {
					roleList = append(roleList, r.GetMenuRoleTypes())
				}
			} else {
				if permission.HasAllPermission(r.PolicyPermissions, p) {
					roleList = append(roleList, r.GetMenuRoleTypes())
				}
			}
		}

		permissionList := make([]types.RolePolicy, 0, len(permissionLst))
		for _, rp := range permissionLst {
			if rp.Status == db.PolicyStatusBanned {
				continue
			}
			if permission.CheckPermission(p, rp.Permission) {
				permissionList = append(permissionList, rp.GetRolePermission())
			}
		}

		subPolicyList := make([]string, 0, len(jwt.UserSubTokenStringList))
		for _, rp := range jwt.UserSubTokenStringList {
			if !permission.CheckPermissionInt64(u.SubPolicy, jwt.UserSubTokenStringMap[rp]) {
				continue
			}
			subPolicyList = append(subPolicyList, rp)
		}

		methodList := make([]string, 0, len(db.PathMethodStringMap))
		for method, rp := range db.PathMethodStringMap {
			if !permission.CheckPermissionInt64(u.Method, rp) {
				continue
			}
			methodList = append(methodList, method)
		}

		nu := warp.UrlPath{
			UrlPath: types.UrlPath{
				ID:             u.Id,
				Describe:       u.Describe,
				Mode:           u.Mode,
				Path:           u.Path,
				Status:         u.Status,
				IsOr:           u.IsOrPolicy,
				Roles:          roleList,
				Policy:         permissionList,
				SubPolicy:      subPolicyList,
				Authentication: u.Authentication,
				DoubleCheck:    u.DoubleCheck,
				CorsMode:       u.CorsMode,
				AdminMode:      u.AdminMode,
				BusyMode:       u.BusyMode,
				BusyCount:      u.BusyCount,
				Method:         methodList,
				CaptchaMode:    u.CaptchaMode,
			},
			PolicyPermission:    p,
			SubPolicyPermission: u.SubPolicy,
			MethodPermission:    u.Method,
		}

		if nu.Mode == db.PathModeRegex {
			nu.Regex, err = regexp.Compile(nu.Path)
			if err != nil {
				continue
			}
		}

		urlPath[u.Id] = nu
	}

	return urlPath
}
