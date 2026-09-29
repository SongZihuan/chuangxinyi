package warp

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"math/big"
)

type Permission struct {
	types.Policy
	Permission big.Int `json:"-"`
}

func (p Permission) GetPolicyType() types.Policy {
	return p.Policy
}

func (p Permission) GetRolePermission() types.RolePolicy {
	return types.RolePolicy{
		ID:          p.ID,
		Sign:        p.Sign,
		Name:        p.Name,
		Sort:        p.Sort,
		IsUser:      p.IsUser,
		IsAnonymous: p.IsAnonymous,
		Describe:    p.Describe,
		Status:      p.Status,
	}
}
