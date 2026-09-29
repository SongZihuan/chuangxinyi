package warp

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"math/big"
)

type WebsitePermission struct {
	types.WebsitePolicy
	Permission big.Int `json:"-"`
}

func (p WebsitePermission) GetWebsitePolicyType() types.WebsitePolicy {
	return p.WebsitePolicy
}

func (p WebsitePermission) GetWebsiteLittlePolicyType() types.WebsiteLittlePolicy {
	return types.WebsiteLittlePolicy{
		ID:       p.ID,
		Sign:     p.Sign,
		Name:     p.Name,
		Sort:     p.Sort,
		Describe: p.Describe,
		Status:   p.Status,
	}
}
