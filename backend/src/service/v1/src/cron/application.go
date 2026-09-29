package cron

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"math/big"
	"sort"
)

func ApplicationCron(webAllP big.Int, website map[int64]warp.Website) (map[string]warp.Application, []warp.Application) {
	defer utils.Recover(logger.Logger, nil, "")

	applicationModel := db.NewApplicationModel(mysql.MySQLConn)
	applicationList, err := applicationModel.GetList(context.Background())
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return make(map[string]warp.Application, 0), make([]warp.Application, 0)
	}

	applicationMap := make(map[string]warp.Application, len(applicationList))
	applicationAllList := make([]warp.Application, 0, len(applicationList))

	for _, a := range applicationList {
		web := getWebsite(a.WebId, website, webAllP)
		if web.ID == warp.UnknownWebsite {
			continue
		}

		np := warp.Application{
			AdminApplication: types.AdminApplication{
				ID:       a.Id,
				Name:     a.Name,
				Sort:     a.Sort,
				Describe: a.Describe,
				WebID:    web.ID,
				WebName:  web.Name,
				WebUID:   web.UID,
				Url:      a.Url,
				Icon:     a.Icon,
				Status:   a.Status,
			},
		}

		applicationMap[a.Name] = np
		applicationAllList = append(applicationAllList, np)
	}

	sort.Slice(applicationAllList, func(i, j int) bool {
		return applicationAllList[i].Sort < applicationAllList[j].Sort
	})

	return applicationMap, applicationAllList
}
