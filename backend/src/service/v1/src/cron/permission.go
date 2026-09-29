package cron

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/config"
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

func PermissionCron(roles map[int64]warp.Role) (map[string]warp.Permission, []warp.Permission, big.Int, big.Int, big.Int) {
	defer utils.Recover(logger.Logger, nil, "")

	permissionModel := db.NewPolicyModel(mysql.MySQLConn)
	permissionList, err := permissionModel.GetList(context.Background())
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return make(map[string]warp.Permission, 0), make([]warp.Permission, 0), big.Int{}, big.Int{}, big.Int{}
	}

	permissionMap := make(map[string]warp.Permission, len(permissionList))
	permissionAllList := make([]warp.Permission, 0, len(permissionList))
	var allPermission, userPermission, anonymousPermission big.Int

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

		if m.Status == db.PolicyStatusOK && m.IsUser {
			var res big.Int
			ok = res.Or(&userPermission, &p)
			if ok == nil {
				continue
			}
			userPermission = res
		}

		if m.Status == db.PolicyStatusOK && m.IsAnonymous {
			var res big.Int
			ok = res.Or(&anonymousPermission, &p)
			if ok == nil {
				continue
			}
			anonymousPermission = res
		}

		roleList := make([]types.MenuRole, 0, len(roles))
		for _, r := range roles {
			if r.IsBanned() {
				continue
			}
			if r.Sign == config.BackendConfig.Admin.RootRole.RoleSign { // 根用户拥有所有权限
				roleList = append(roleList, r.GetMenuRoleTypes())
			} else {
				if permission.CheckPermission(r.PolicyPermissions, p) {
					roleList = append(roleList, r.GetMenuRoleTypes())
				}
			}
		}

		np := warp.Permission{
			Policy: types.Policy{
				ID:          m.Id,
				Sign:        m.Sign,
				Name:        m.Name,
				Sort:        m.Sort,
				Describe:    m.Describe,
				Status:      m.Status,
				Roles:       roleList,
				IsAnonymous: m.IsAnonymous,
				IsUser:      m.IsUser,
			},
			Permission: p,
		}

		permissionMap[m.Sign] = np
		permissionAllList = append(permissionAllList, np)
	}

	sort.Slice(permissionAllList, func(i, j int) bool {
		return permissionAllList[i].Sort < permissionAllList[j].Sort
	})

	return permissionMap, permissionAllList, allPermission, userPermission, anonymousPermission
}
