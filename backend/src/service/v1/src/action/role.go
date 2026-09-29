package action

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"time"
)

func GetRootRole() warp.Role {
	belong := GetWebsite(warp.UserCenterWebsite)

	menuList := make([]types.RoleMenu, 0, len(model.Menus()))
	for _, m := range model.Menus() {
		menuList = append(menuList, m.GetRoleMenuType())
	}

	urlPathList := make([]types.RoleUrlPath, 0, len(model.UrlPathMap()))
	for _, u := range model.UrlPathMap() {
		urlPathList = append(urlPathList, u.GetRoleUrlPathType())
	}

	return warp.Role{
		Role: types.Role{
			ID:                   0,
			Name:                 config.BackendConfig.Admin.RootRole.RoleName,
			Describe:             config.BackendConfig.Admin.RootRole.RoleDescribe,
			Sign:                 config.BackendConfig.Admin.RootRole.RoleSign,
			NotDelete:            config.BackendConfig.Admin.RootRole.NotDelete,
			NotChangeSign:        config.BackendConfig.Admin.RootRole.NotChangeSign,
			NotChangePermissions: config.BackendConfig.Admin.RootRole.NotChangePermissions,
			Status:               db.RoleStatusOK,
			Belong:               belong.ID,
			BelongName:           belong.Name,
			CreateAt:             time.Now().Unix(),
			Menus:                menuList,
			UrlPaths:             urlPathList,
		},
		PolicyPermissions: *model.AllPermission(),
	}
}

func GetUserRole() warp.Role {
	belong := GetWebsite(warp.UserCenterWebsite)

	menuList := make([]types.RoleMenu, 0, len(model.Menus()))
	for _, m := range model.Menus() {
		if m.IsOr {
			if permission.HasOnePermission(*model.UserPermission(), m.PolicyPermission) {
				menuList = append(menuList, m.GetRoleMenuType())
			}
		} else {
			if permission.HasAllPermission(*model.UserPermission(), m.PolicyPermission) {
				menuList = append(menuList, m.GetRoleMenuType())
			}
		}
	}

	urlPathList := make([]types.RoleUrlPath, 0, len(model.UrlPathMap()))
	for _, u := range model.UrlPathMap() {
		if u.IsOr {
			if permission.HasOnePermission(*model.UserPermission(), u.PolicyPermission) {
				urlPathList = append(urlPathList, u.GetRoleUrlPathType())
			}
		} else {
			if permission.HasAllPermission(*model.UserPermission(), u.PolicyPermission) {
				urlPathList = append(urlPathList, u.GetRoleUrlPathType())
			}
		}
	}

	return warp.Role{
		Role: types.Role{
			ID:                   0,
			Name:                 config.BackendConfig.Admin.UserRole.RoleName,
			Describe:             config.BackendConfig.Admin.UserRole.RoleDescribe,
			Sign:                 config.BackendConfig.Admin.UserRole.RoleSign,
			NotDelete:            config.BackendConfig.Admin.UserRole.NotDelete,
			NotChangeSign:        config.BackendConfig.Admin.UserRole.NotChangeSign,
			NotChangePermissions: config.BackendConfig.Admin.UserRole.NotChangePermissions,
			Status:               db.RoleStatusOK,
			Belong:               belong.ID,
			BelongName:           belong.Name,
			CreateAt:             time.Now().Unix(),
			Menus:                menuList,
			UrlPaths:             urlPathList,
		},
		PolicyPermissions: *model.UserPermission(),
	}
}

func GetAnonymousRole() warp.Role {
	belong := GetWebsite(warp.UserCenterWebsite)

	menuList := make([]types.RoleMenu, 0, len(model.Menus()))
	for _, m := range model.Menus() {
		if m.IsOr {
			if permission.HasOnePermission(*model.AnonymousPermission(), m.PolicyPermission) {
				menuList = append(menuList, m.GetRoleMenuType())
			}
		} else {
			if permission.HasAllPermission(*model.AnonymousPermission(), m.PolicyPermission) {
				menuList = append(menuList, m.GetRoleMenuType())
			}
		}
	}

	urlPathList := make([]types.RoleUrlPath, 0, len(model.UrlPathMap()))
	for _, u := range model.UrlPathMap() {
		if u.IsOr {
			if permission.HasOnePermission(*model.AnonymousPermission(), u.PolicyPermission) {
				urlPathList = append(urlPathList, u.GetRoleUrlPathType())
			}
		} else {
			if permission.HasAllPermission(*model.AnonymousPermission(), u.PolicyPermission) {
				urlPathList = append(urlPathList, u.GetRoleUrlPathType())
			}
		}
	}

	return warp.Role{
		Role: types.Role{
			ID:                   0,
			Name:                 config.BackendConfig.Admin.AnonymousRole.RoleName,
			Describe:             config.BackendConfig.Admin.AnonymousRole.RoleDescribe,
			Sign:                 config.BackendConfig.Admin.AnonymousRole.RoleSign,
			NotDelete:            config.BackendConfig.Admin.AnonymousRole.NotDelete,
			NotChangeSign:        config.BackendConfig.Admin.AnonymousRole.NotChangeSign,
			NotChangePermissions: config.BackendConfig.Admin.AnonymousRole.NotChangePermissions,
			Status:               db.RoleStatusOK,
			Belong:               belong.ID,
			BelongName:           belong.Name,
			CreateAt:             time.Now().Unix(),
			Menus:                menuList,
			UrlPaths:             urlPathList,
		},
		PolicyPermissions: *model.AnonymousPermission(),
	}
}

func GetRole(id int64, isAdmin bool) warp.Role {
	var res *warp.Role = nil
	if isAdmin {
		tmp, ok := (model.RolesSign())[config.BackendConfig.Admin.RootRole.RoleSign]
		if ok {
			res = &tmp
			res.Status = db.RoleStatusOK
		} else {
			tmp = GetRootRole()
			res = &tmp
		}
	}

	if res == nil {
		tmp, ok := (model.Roles())[id]
		if ok {
			res = &tmp
		}
	}

	if res == nil {
		tmp, ok := (model.RolesSign())[config.BackendConfig.Admin.AnonymousRole.RoleSign]
		if ok {
			res = &tmp
		}
	}

	if res == nil {
		return GetAnonymousRole()
	} else if res.IsBanned() {
		return GetAnonymousRole()
	} else {
		return *res
	}
}

func GetRoleWithoutBanned(id int64, isAdmin bool) warp.Role {
	var res *warp.Role = nil
	if isAdmin {
		tmp, ok := (model.RolesSign())[config.BackendConfig.Admin.RootRole.RoleSign]
		if ok {
			res = &tmp
			res.Status = db.RoleStatusOK
		} else {
			tmp = GetRootRole()
			res = &tmp
		}
	}

	if res == nil {
		tmp, ok := (model.Roles())[id]
		if ok {
			res = &tmp
		}
	}

	if res == nil {
		tmp, ok := (model.RolesSign())[config.BackendConfig.Admin.AnonymousRole.RoleSign]
		if ok {
			res = &tmp
		}
	}

	if res == nil {
		return GetAnonymousRole()
	} else if res.IsBanned() {
		return GetAnonymousRole()
	} else {
		return *res
	}
}

func GetRoleBySign(sign string, isAdmin bool) warp.Role {
	var res *warp.Role = nil
	if isAdmin {
		tmp, ok := (model.RolesSign())[config.BackendConfig.Admin.RootRole.RoleSign]
		if ok {
			res = &tmp
			res.Status = db.RoleStatusOK
		} else {
			tmp = GetRootRole()
			res = &tmp
		}
	}

	if res == nil {
		tmp, ok := (model.RolesSign())[sign]
		if ok {
			res = &tmp
		} else {
			res = nil
		}
	}

	if res == nil {
		tmp, ok := (model.RolesSign())[config.BackendConfig.Admin.AnonymousRole.RoleSign]
		if ok {
			res = &tmp
		}
	}

	if res == nil {
		return GetAnonymousRole()
	} else if res.IsBanned() {
		return GetAnonymousRole()
	} else {
		return *res
	}
}

func ClearRoleMenu(menu []types.RoleMenu, role warp.Role, subType int) []types.RoleMenu {
	res := make([]types.RoleMenu, 0, len(menu))
	if role.Sign == config.BackendConfig.Admin.RootRole.RoleSign && (subType == jwt.UserRootToken || subType == jwt.UserHighAuthorityRootToken) {
		return menu
	}

	subTypePermission, ok := jwt.UserSubTokenPermissionMap[subType]
	if !ok {
		return []types.RoleMenu{}
	}

	for _, m := range menu {
		srcMenu, ok := (model.Menus())[m.ID]
		if !ok || !permission.CheckPermissionInt64(srcMenu.SubPolicyPermission, subTypePermission) {
			continue
		}
		res = append(res, m)
	}

	return res
}

func ClearRoleUrlPath(path []types.RoleUrlPath, role warp.Role, subType int) []types.RoleUrlPath {
	res := make([]types.RoleUrlPath, 0, len(path))
	if role.Sign == config.BackendConfig.Admin.RootRole.RoleSign && (subType == jwt.UserRootToken || subType == jwt.UserHighAuthorityRootToken) {
		return path
	}

	subTypePermission, ok := jwt.UserSubTokenPermissionMap[subType]
	if !ok {
		return []types.RoleUrlPath{}
	}

	for _, p := range path {
		srcPathUrl, ok := (model.UrlPathMap())[p.ID]
		if !ok || !permission.CheckPermissionInt64(srcPathUrl.SubPolicyPermission, subTypePermission) {
			continue
		}
		res = append(res, p)
	}

	return res
}
