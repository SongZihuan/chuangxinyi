package utils

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
)

func DeleteUser(user *db.User, status int64, limit int64) errors.WTError {
	var err error
	if limit < 0 {
		return nil
	}

	if user.IsAdmin {
		return errors.Errorf("admin can not delete")
	}

	if db.IsKeepInfoStatus(status) {
		return nil
	}

	err = mysql.MySQLConn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		userModel := db.NewUserModelWithSession(session)
		phoneModel := db.NewPhoneModelWithSession(session)
		companyModel := db.NewCompanyModelWithSession(session)
		emailModel := db.NewEmailModelWithSession(session)
		idcardModel := db.NewIdcardModelWithSession(session)
		usernameModel := db.NewUsernameModelWithSession(session)
		wechatModel := db.NewWechatModelWithSession(session)

		user.Status = status
		err = userModel.UpdateCh(ctx, user) // 需要更新status
		if err != nil {
			logger.Logger.Tag("A1")
			return err
		}

		_, err = companyModel.DeleteByUserID(ctx, user.Id)
		if err != nil {
			logger.Logger.Tag("A2")
			return err
		}

		_, err = emailModel.DeleteByUserID(ctx, user.Id)
		if err != nil {
			logger.Logger.Tag("A3")
			return err
		}

		_, err = idcardModel.DeleteByUserID(ctx, user.Id)
		if err != nil {
			logger.Logger.Tag("A4")
			return err
		}

		_, err = phoneModel.DeleteByUserID(ctx, user.Id)
		if err != nil {
			logger.Logger.Tag("A5")
			return err
		}

		_, err = usernameModel.DeleteByUserID(ctx, user.Id)
		if err != nil {
			logger.Logger.Tag("A6")
			return err
		}

		_, err = wechatModel.DeleteByUserID(ctx, user.Id)
		if err != nil {
			logger.Logger.Tag("A7")
			return err
		}

		return nil
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	_ = jwt.DeleteAllUserToken(context.Background(), user.Uid, "")
	_ = jwt.DeleteAllUserWebsiteToken(context.Background(), user.Uid)
	_ = jwt.DeleteAllUserSonToken(context.Background(), user.Uid)
	_ = jwt.DeleteAllFatherUserToken(context.Background(), user.Uid)

	userModel := db.NewUserModel(mysql.MySQLConn)
	sonList, err := userModel.GetSonList(context.Background(), user.Id, []string{"NORMAL", "REGISTER"}, "") // 先删除子用户
	if err != nil {
		return errors.WarpQuick(err)
	} else {
		for _, son := range sonList {
			_ = DeleteUser(&son, status, limit-1)
			// 忽略接下来的错误
		}
	}

	return nil
}
