package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	errors "github.com/wuntsong-org/wterrors"
)

var UserNotFound = errors.Errorf("user not found")

func GetUser(ctx context.Context, id int64, uid string, findBanned bool) (user *db.User, err error) {
	userModel := db.NewUserModel(mysql.MySQLConn)

	if len(uid) == 0 {
		user, err = userModel.FindOneByIDWithoutDelete(ctx, id)
	} else {
		user, err = utils2.FindUser(ctx, uid, true)
	}
	if errors.Is(err, utils2.UserNotFound) || errors.Is(err, db.ErrNotFound) {
		return nil, UserNotFound
	} else if err != nil {
		return nil, err
	} else if !findBanned && db.IsBanned(user) {
		return nil, UserNotFound
	}

	return user, nil
}
