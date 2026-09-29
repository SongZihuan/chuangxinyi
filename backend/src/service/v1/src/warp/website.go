package warp

import (
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"math/big"
)

type Website struct {
	types.Website
	PolicyPermissions big.Int `json:"-"`
}

const UserCenterWebsite = 0 // 必须是0
const UnknownWebsite = -1

func (w Website) GetWebsiteType() types.Website {
	return w.Website
}

func (w Website) GetLittleWebsiteType() types.LittleWebiste {
	return types.LittleWebiste{
		ID:        w.ID,
		Name:      w.Name,
		Describe:  w.Describe,
		KeyMap:    w.KeyMap,
		PubKey:    w.PubKey,
		Agreement: w.Agreement,
		Status:    w.Status,
		CreateAt:  w.CreateAt,
	}
}

func (w Website) GetWebsiteEasyType() types.WebsiteEasy {
	return types.WebsiteEasy{
		ID:   w.ID,
		Name: w.Name,
	}
}

func (w Website) GetGetDomainDataType() types.GetDomainData {
	return types.GetDomainData{
		Name:      w.Name,
		Describe:  w.Describe,
		KeyMap:    w.KeyMap,
		Agreement: w.Agreement,
	}
}

func (w Website) GeLittleWebsiteType() types.LittleWebiste {
	return types.LittleWebiste{
		ID:        w.ID,
		Name:      w.Name,
		Describe:  w.Describe,
		KeyMap:    w.KeyMap,
		PubKey:    w.PubKey,
		Agreement: w.Agreement,
		CreateAt:  w.CreateAt,
	}
}

func (w Website) GetWebsiteIPListType() []types.WebsiteIP {
	return w.IP
}

func (w Website) GetWebsiteIPStringListType() []string {
	res := make([]string, 0, len(w.IP))
	for _, i := range w.IP {
		res = append(res, i.IP)
	}

	return res
}

func (w Website) GetWebsiteDomainListType() []types.WebsiteDomain {
	return w.Domain
}
