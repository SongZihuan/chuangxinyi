package warp

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"math/big"
)

type Role struct {
	types.Role
	PolicyPermissions big.Int `json:"-"`
}

func (m Menu) GetRoleMenuType() types.RoleMenu {
	return types.RoleMenu{
		ID:             m.ID,
		Sort:           m.Sort,
		FatherID:       m.FatherID,
		Name:           m.Name,
		Path:           m.Path,
		Title:          m.Title,
		Icon:           m.Icon,
		Redirect:       m.Redirect,
		Superior:       m.Superior,
		Category:       m.Category,
		Component:      m.Component,
		ComponentAlias: m.ComponentAlias,
		MetaLink:       m.MetaLink,
		Type:           m.Type,
		IsLink:         m.IsLink,
		IsHide:         m.IsHide,
		IsKeepalive:    m.IsKeepalive,
		IsAffix:        m.IsAffix,
		IsIframe:       m.IsIframe,
		BtnPower:       m.BtnPower,
	}
}

func (r Role) GetMenuRoleTypes() types.MenuRole {
	return types.MenuRole{
		ID:                   r.ID,
		Describe:             r.Describe,
		Name:                 r.Name,
		Sign:                 r.Sign,
		NotDelete:            r.NotDelete,
		NotChangePermissions: r.NotChangePermissions,
		NotChangeSign:        r.NotChangeSign,
		Status:               r.Status,
		CreateAt:             r.CreateAt,
		Belong:               r.Belong,
	}
}

func (r Role) GetRole() types.Role {
	return r.Role
}

func (r *Role) IsBanned() bool {
	if r.Sign == config.BackendConfig.Admin.RootRole.RoleSign {
		return false
	}

	return r.Status == db.RoleStatusBanned
}
