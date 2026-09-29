package record

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
)

type Record struct {
	RequestsID     string
	User           *db.User
	UserToken      string
	Role           *warp.Role
	Website        *warp.Website
	RequestWebsite *warp.Website
	Msg            string
	Err            errors.WTError
	Stack          string
}

func GetRecord(ctx context.Context) *Record {
	res, ok := ctx.Value("X-Record").(*Record)
	if !ok {
		logger.Logger.Error("bad X-Record")
		return &Record{RequestsID: "unknown"} // 返回一个无效的，防止报错
	}
	return res
}

func GetRecordIfExists(ctx context.Context) *Record {
	res, ok := ctx.Value("X-Record").(*Record)
	if !ok {
		return &Record{RequestsID: "unknown"} // 返回一个无效的，防止报错
	}
	return res
}
