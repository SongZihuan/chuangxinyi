package action

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"math/big"
	"time"
)

func GetWebsiteByUID(webID string) warp.Website {
	if webID == config.BackendConfig.User.WebsiteUID {
		return warp.Website{
			Website: types.Website{
				ID:       warp.UserCenterWebsite,
				UID:      config.BackendConfig.User.WebsiteUID,
				Name:     config.BackendConfig.User.ReadableName,
				Describe: "create by system",
				CreateAt: time.Now().Unix(),
				Status:   db.WebsiteStatusOK,
			},
			PolicyPermissions: *model.WebsiteAllPermission(),
		}
	}

	web, ok := (model.WebsitesUID())[webID]
	if ok {
		return web
	}

	return warp.Website{
		Website: types.Website{
			ID:       warp.UnknownWebsite,
			UID:      "",
			Name:     "未知",
			Describe: "create by system",
			CreateAt: time.Now().Unix(),
			Status:   db.WebsiteStatusBanned,
		},
		PolicyPermissions: *big.NewInt(0),
	}
}

func GetWebsite(webID int64) warp.Website {
	if webID == warp.UserCenterWebsite { // webID == 0
		return warp.Website{
			Website: types.Website{
				ID:       warp.UserCenterWebsite,
				UID:      config.BackendConfig.User.WebsiteUID,
				Name:     config.BackendConfig.User.ReadableName,
				Describe: "create by system",
				CreateAt: time.Now().Unix(),
				Status:   db.WebsiteStatusOK,
			},
			PolicyPermissions: *model.WebsiteAllPermission(),
		}
	}

	web, ok := (model.Websites())[webID]
	if ok {
		return web
	}

	return warp.Website{
		Website: types.Website{
			ID:       warp.UnknownWebsite,
			UID:      "",
			Name:     "未知",
			Describe: "create by system",
			CreateAt: time.Now().Unix(),
			Status:   db.WebsiteStatusBanned,
		},
		PolicyPermissions: *big.NewInt(0),
	}
}
