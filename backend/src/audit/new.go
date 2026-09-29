package audit

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
)

func NewUserAudit(userID int64, msg string, args ...any) {
	go func() {
		auditModel := db.NewAuditModel(mysql.MySQLConn)
		_, err := auditModel.Insert(context.Background(), &db.Audit{
			UserId:  userID,
			Content: fmt.Sprintf(msg, args...),
			From:    config.BackendConfig.User.ReadableName,
			FromId:  warp.UserCenterWebsite,
		})
		if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			return
		}
	}()
}

func NewAdminAudit(userID int64, msg string, args ...any) {
	go func() {
		auditModel := db.NewAuditModel(mysql.MySQLConn)
		_, err := auditModel.Insert(context.Background(), &db.Audit{
			UserId:  userID,
			Content: fmt.Sprintf(msg, args...),
			From:    fmt.Sprintf("%s（管理员）", config.BackendConfig.User.ReadableName),
			FromId:  warp.UserCenterWebsite,
		})
		if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			return
		}
	}()
}

func NewOtherAudit(userID int64, from string, fromID int64, msg string, args ...any) errors.WTError {
	auditModel := db.NewAuditModel(mysql.MySQLConn)
	_, err := auditModel.Insert(context.Background(), &db.Audit{
		UserId:  userID,
		Content: fmt.Sprintf(msg, args...),
		From:    from,
		FromId:  fromID,
	})
	if err != nil {
		return errors.WarpQuick(err)
	}
	return nil
}
