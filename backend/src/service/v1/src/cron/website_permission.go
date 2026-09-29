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
	"sort"
)

func WebsitePermissionCron(website map[int64]warp.Website) (map[string]warp.WebsitePermission, []warp.WebsitePermission, big.Int) {
	defer utils.Recover(logger.Logger, nil, "")

	permissionModel := db.NewWebsitePolicyModel(mysql.MySQLConn)
	permissionList, err := permissionModel.GetList(context.Background())
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return make(map[string]warp.WebsitePermission), make([]warp.WebsitePermission, 0), big.Int{}
	}

	permissionMap := make(map[string]warp.WebsitePermission, len(permissionList))
	permissionAllList := make([]warp.WebsitePermission, 0, len(permissionList))
	var allPermission big.Int

	for _, m := range permissionList {
		var p big.Int
		ok := p.Lsh(big.NewInt(1), uint(m.Id-10000))
		if ok == nil {
			continue
		}

		if m.Status == db.PolicyStatusOK {
			var res big.Int
			ok = res.Or(&allPermission, &p)
			if ok == nil {
				continue
			}
			allPermission = res

		}

		websiteList := make([]types.LittleWebiste, 0, len(website))
		for _, w := range website {
			if w.Status == db.WebsiteStatusBanned {
				continue
			}
			if permission.CheckPermission(w.PolicyPermissions, p) {
				websiteList = append(websiteList, w.GeLittleWebsiteType())
			}
		}

		np := warp.WebsitePermission{
			WebsitePolicy: types.WebsitePolicy{
				ID:       m.Id,
				Sign:     m.Sign,
				Name:     m.Name,
				Sort:     m.Sort,
				Describe: m.Describe,
				Status:   m.Status,
				Websites: websiteList,
			},
			Permission: p,
		}

		permissionMap[m.Sign] = np
		permissionAllList = append(permissionAllList, np)
	}

	sort.Slice(permissionAllList, func(i, j int) bool {
		return permissionAllList[i].Sort < permissionAllList[j].Sort
	})

	return permissionMap, permissionAllList, allPermission
}
