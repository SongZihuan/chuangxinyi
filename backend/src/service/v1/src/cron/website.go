package cron

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"math/big"
)

func WebsiteCron(allP big.Int, urlpath map[int64]warp.WebsiteUrlPath, permissionLst []warp.WebsitePermission) (map[int64]warp.Website, map[string]warp.Website, []warp.Website) {
	websiteModel := db.NewWebsiteModel(mysql.MySQLConn)
	websiteIPModel := db.NewWebsiteIpModel(mysql.MySQLConn)
	websiteDomainModel := db.NewWebsiteDomainModel(mysql.MySQLConn)
	res, err := websiteModel.GetList(context.Background())
	if err != nil {
		logger.Logger.Error("mysql sql error: %s", err.Error())
		return make(map[int64]warp.Website), make(map[string]warp.Website), make([]warp.Website, 0)
	}

	websitesList := make([]warp.Website, 0, len(res))
	websitesByID := make(map[int64]warp.Website, len(res))
	websitesByUID := make(map[string]warp.Website, len(res))

	for _, r := range res {
		var p big.Int
		_, ok := p.SetString(r.Permission, 16)
		if !ok {
			continue
		}
		p = permission.ClearPermission(allP, p)

		ipList, err := websiteIPModel.GetList(context.Background(), r.Id)
		if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			continue
		}

		webIP := make([]types.WebsiteIP, 0, len(ipList))
		for _, m := range ipList {
			webIP = append(webIP, types.WebsiteIP{
				ID: m.Id,
				IP: m.Ip,
			})
		}

		domainList, err := websiteDomainModel.GetList(context.Background(), r.Id)
		if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			continue
		}

		webDomain := make([]types.WebsiteDomain, 0, len(domainList))
		for _, m := range domainList {
			webDomain = append(webDomain, types.WebsiteDomain{
				ID:     m.Id,
				Domain: m.Domain,
			})
		}

		urlPathList := make([]types.WebsiteLittleUrlPath, 0, len(urlpath))
		for _, u := range urlpath {
			if u.Status == db.WebsitePathStatusDelete {
				continue
			}
			if u.IsOr {
				if permission.HasOnePermission(p, u.PolicyPermission) {
					urlPathList = append(urlPathList, u.GetWebsiteLittleUrlPathType())
				}
			} else {
				if permission.HasAllPermission(p, u.PolicyPermission) {
					urlPathList = append(urlPathList, u.GetWebsiteLittleUrlPathType())
				}
			}
		}

		permissionList := make([]types.WebsiteLittlePolicy, 0, len(permissionLst))
		for _, wp := range permissionLst {
			if wp.Status == db.WebsitePolicyStatusBanned {
				continue
			}
			if permission.CheckPermission(p, wp.Permission) {
				permissionList = append(permissionList, wp.GetWebsiteLittlePolicyType())
			}
		}

		var keyMap map[string]string
		err = utils.JsonUnmarshal([]byte(r.Keymap), &keyMap)
		if err != nil {
			logger.Logger.Error("utils.JsonUnmarshal: %s", err.Error())
			continue
		}

		km := make([]types.LabelValueRecord, 0, len(keyMap))
		for l, v := range keyMap {
			km = append(km, types.LabelValueRecord{
				Label: l,
				Value: v,
			})
		}

		w := warp.Website{
			Website: types.Website{
				ID:        r.Id,
				UID:       r.Uid,
				Name:      r.Name,
				Describe:  r.Describe,
				KeyMap:    km,
				PubKey:    r.Pubkey,
				Agreement: r.Agreement,
				Status:    r.Status,
				CreateAt:  r.CreateAt.Unix(),
				UrlPath:   urlPathList,
				Policy:    permissionList,
				IP:        webIP,
				Domain:    webDomain,
			},
			PolicyPermissions: p,
		}

		websitesList = append(websitesList, w)
		websitesByID[r.Id] = w
		websitesByUID[r.Uid] = w
	}

	return websitesByID, websitesByUID, websitesList
}
