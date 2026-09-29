package cron

import (
	"context"
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

func WebsiteUrlPathCron(allP big.Int, website map[int64]warp.Website, permissionLst []warp.WebsitePermission) map[int64]warp.WebsiteUrlPath {
	defer utils.Recover(logger.Logger, nil, "")

	urlModel := db.NewWebsiteUrlPathModel(mysql.MySQLConn)
	urlList, err := urlModel.GetList(context.Background())
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return make(map[int64]warp.WebsiteUrlPath)
	}

	urlPath := make(map[int64]warp.WebsiteUrlPath, len(urlList))

	for _, u := range urlList {
		var p big.Int
		_, ok := p.SetString(u.Permission, 16)
		if !ok {
			continue
		}
		p = permission.ClearPermission(allP, p)

		websiteList := make([]types.LittleWebiste, 0, len(website))
		for _, w := range website {
			if w.Status == db.WebsiteStatusBanned {
				continue
			}
			if u.IsOrPolicy {
				if permission.HasOnePermission(w.PolicyPermissions, p) {
					websiteList = append(websiteList, w.GetLittleWebsiteType())
				}
			} else {
				if permission.HasAllPermission(w.PolicyPermissions, p) {
					websiteList = append(websiteList, w.GetLittleWebsiteType())
				}
			}
		}

		permissionList := make([]types.WebsiteLittlePolicy, 0, len(permissionLst))
		for _, wp := range permissionLst {
			if wp.Status == db.WebsitePolicyStatusBanned {
				continue
			}
			if permission.CheckPermission(p, wp.Permission) {
				permissionList = append(permissionList, wp.GetWebsiteLittlePolicyType())
			}
		}

		methodList := make([]string, 0, len(db.WebsitePathMethodStringMap))
		for method, rp := range db.WebsitePathMethodStringMap {
			if !permission.CheckPermissionInt64(u.Method, rp) {
				continue
			}
			methodList = append(methodList, method)
		}

		nu := warp.WebsiteUrlPath{
			WebsiteUrlPath: types.WebsiteUrlPath{
				ID:       u.Id,
				Describe: u.Describe,
				Mode:     u.Mode,
				Path:     u.Path,
				Status:   u.Status,
				IsOr:     u.IsOrPolicy,
				Websites: websiteList,
				Policy:   permissionList,
				Method:   methodList,
			},
			PolicyPermission: p,
			MethodPermission: u.Method,
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
