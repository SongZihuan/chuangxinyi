package warp

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"math/big"
)

type Menu struct {
	types.Menu
	PolicyPermission    big.Int `json:"-"`
	SubPolicyPermission int64   `json:"-"`
}

func (m Menu) GetMenuType() types.Menu {
	return m.Menu
}
