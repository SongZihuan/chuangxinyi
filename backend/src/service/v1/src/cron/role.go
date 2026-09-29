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
	"math/big"
	"sort"
)

func RoleCron(allP big.Int, webAllP big.Int, website map[int64]warp.Website, menus map[int64]warp.Menu, urlpath map[int64]warp.UrlPath, permissionLst []warp.Permission) (map[int64]warp.Role, map[string]warp.Role) {
	roleModel := db.NewRoleModel(mysql.MySQLConn)
	res, err := roleModel.GetList(context.Background())
	if err != nil {
		logger.Logger.Error("mysql sql error: %s", err.Error())
		return make(map[int64]warp.Role), make(map[string]warp.Role, 0)
	}

	roles := make(map[int64]warp.Role, len(res))
	rolesSign := make(map[string]warp.Role, len(res))

	for _, r := range res {
		var p big.Int
		if r.Sign == config.BackendConfig.Admin.RootRole.RoleSign {
			p = allP
		} else {
			_, ok := p.SetString(r.Permissions, 16)
			if !ok {
				continue
			}
			p = permission.ClearPermission(allP, p)
		}

		menuList := make([]types.RoleMenu, 0, len(menus))
		for _, m := range menus {
			if m.Status == db.MenuStatusBanned {
				continue
			}
			if r.Sign == config.BackendConfig.Admin.RootRole.RoleSign { // 根用户拥有所有列表
				menuList = append(menuList, m.GetRoleMenuType())
			} else if m.IsOr {
				if permission.HasOnePermission(p, m.PolicyPermission) {
					menuList = append(menuList, m.GetRoleMenuType())
				}
			} else {
				if permission.HasAllPermission(p, m.PolicyPermission) {
					menuList = append(menuList, m.GetRoleMenuType())
				}
			}
		}

		urlPathList := make([]types.RoleUrlPath, 0, len(urlpath))
		for _, u := range urlpath {
			if u.Status == db.PathStatusDelete {
				continue
			}
			if r.Sign == config.BackendConfig.Admin.RootRole.RoleSign { // 根用户拥有所有列表
				urlPathList = append(urlPathList, u.GetRoleUrlPathType())
			} else if u.IsOr {
				if permission.HasOnePermission(p, u.PolicyPermission) {
					urlPathList = append(urlPathList, u.GetRoleUrlPathType())
				}
			} else {
				if permission.HasAllPermission(p, u.PolicyPermission) {
					urlPathList = append(urlPathList, u.GetRoleUrlPathType())
				}
			}
		}

		permissionList := make([]types.RolePolicy, 0, len(permissionLst))
		for _, rp := range permissionLst {
			if rp.Status == db.PolicyStatusBanned {
				continue
			}
			if r.Sign == config.BackendConfig.Admin.RootRole.RoleSign { // 根用户拥有所有权限
				permissionList = append(permissionList, rp.GetRolePermission())
			} else if permission.CheckPermission(p, rp.Permission) {
				permissionList = append(permissionList, rp.GetRolePermission())
			}
		}

		belong := getWebsite(r.Belong.Int64, website, webAllP)

		nr := warp.Role{
			Role: types.Role{
				ID:                   r.Id,
				Name:                 r.Name,
				Sign:                 r.Sign,
				Describe:             r.Describe,
				Status:               r.Status,
				NotDelete:            r.NotDelete,
				NotChangePermissions: r.NotChangePermissions,
				NotChangeSign:        r.NotChangeSign,
				CreateAt:             r.CreateAt.Unix(),
				Belong:               belong.ID,
				BelongName:           belong.Name,
				Menus:                menuList,
				UrlPaths:             urlPathList,
				Permission:           permissionList,
			},
			PolicyPermissions: p,
		}

		if nr.Sign == config.BackendConfig.Admin.RootRole.RoleSign {
			nr.Status = db.RoleStatusOK
		}

		sort.Slice(nr.Menus, func(i, j int) bool {
			return nr.Menus[i].Sort < nr.Menus[j].Sort
		})

		roles[r.Id] = nr
		rolesSign[r.Sign] = nr
	}

	return roles, rolesSign
}
