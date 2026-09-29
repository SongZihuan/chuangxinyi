package fuwuhao

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/fastwego/offiaccount/apis/user"
	"github.com/wuntsong-org/wterrors"
	"net/url"
)

func Bind(ctx context.Context, openID string) bool {
	params := url.Values{}
	params.Add("openid", openID)
	params.Add("lang", "zh_CN")

	resp, err := user.GetUserInfo(OffiAccount, params)
	if err != nil {
		logger.Logger.Error("get user info error: %s", err.Error())
		return false
	}

	data := struct {
		UnionID string `json:"unionID"`
	}{}

	err = utils.JsonUnmarshal(resp, &data)
	if err != nil {
		logger.Logger.Error("get user info error: %s", err.Error())
		return false
	}

	if len(data.UnionID) == 0 {
		return false
	}

	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	w, err := wechatModel.FindByUnionID(ctx, data.UnionID)
	if errors.Is(err, db.ErrNotFound) {
		return false
	} else if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return false
	}

	w.Id = 0
	w.Fuwuhao = sql.NullString{
		Valid:  true,
		String: openID,
	}

	_, err = wechatModel.InsertWithDelete(ctx, w)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return false
	}

	go func() {
		_ = SendBindSuccess(context.Background(), openID, w.UserId)
	}()

	audit.NewUserAudit(w.UserId, "用户已成功绑定服务号")

	return true
}
