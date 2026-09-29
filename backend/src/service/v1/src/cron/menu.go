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
)

func MenuCron(allP big.Int, roles map[int64]warp.Role, permissionLst []warp.Permission) map[int64]warp.Menu {
	defer utils.Recover(logger.Logger, nil, "")

	menuModel := db.NewMenuModel(mysql.MySQLConn)
	menuList, err := menuModel.GetList(context.Background())
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return make(map[int64]warp.Menu)
	}

	menus := make(map[int64]warp.Menu, len(menuList))

	for _, m := range menuList {
		var p big.Int
		_, ok := p.SetString(m.Policy, 16)
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
			} else if m.IsOrPolicy {
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
			if !permission.CheckPermissionInt64(m.SubPolicy, jwt.UserSubTokenStringMap[rp]) {
				continue
			}
			subPolicyList = append(subPolicyList, rp)
		}

		nm := warp.Menu{
			Menu: types.Menu{
				ID:             m.Id,
				Sort:           m.Sort,
				FatherID:       m.FatherId.Int64,
				Name:           m.Name,
				Path:           m.Path,
				Title:          m.Title,
				Icon:           m.Icon,
				Redirect:       m.Redirect.String,
				Superior:       m.Superior,
				Category:       m.Category,
				Component:      m.Component,
				ComponentAlias: m.ComponentAlias,
				MetaLink:       m.MetaLink.String,
				Type:           m.Type,
				IsLink:         m.IsLink,
				IsHide:         m.IsHide,
				IsKeepalive:    m.IsKeepalive,
				IsAffix:        m.IsAffix,
				IsIframe:       m.IsIframe,
				BtnPower:       m.BtnPower,
				Roles:          roleList,
				Policy:         permissionList,
				SubPolicy:      subPolicyList,
				Status:         m.Status,
				Describe:       m.Describe,
				IsOr:           m.IsOrPolicy,
			},
			PolicyPermission:    p,
			SubPolicyPermission: m.SubPolicy,
		}

		menus[m.Id] = nm
	}

	return menus
}
