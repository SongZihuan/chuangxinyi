package warp

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"math/big"
	"regexp"
)

type WebsiteUrlPath struct {
	types.WebsiteUrlPath
	PolicyPermission big.Int        `json:"-"`
	Regex            *regexp.Regexp `json:"-"`
	MethodPermission int64          `json:"-"`
}

func (u WebsiteUrlPath) GetWebsiteUrlPathType() types.WebsiteUrlPath {
	return u.WebsiteUrlPath
}

func (u WebsiteUrlPath) GetWebsiteLittleUrlPathType() types.WebsiteLittleUrlPath {
	return types.WebsiteLittleUrlPath{
		ID:       u.ID,
		Path:     u.Path,
		Describe: u.Describe,
		Mode:     u.Mode,
		Status:   u.Status,
	}
}
