package sender

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/email"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"github.com/wuntsong-org/wterrors"
)

func getEmail(userID int64) string {
	emailModel := db.NewEmailModel(mysql.MySQLConn)
	e, err := emailModel.FindByUserID(context.Background(), userID)
	if errors.Is(err, db.ErrNotFound) {
		return ""
	} else if err != nil {
		return ""
	}

	return e.Email.String
}

func EmailSendChange(userID int64, project string) {
	go func() {
		e := getEmail(userID)
		if len(e) == 0 {
			return
		}

		err := email.SendChange(project, e)
		if err != nil {
			logger.Logger.Error("send email error: %s", err.Error())
		}
	}()
}

func EmailSendEmailChange(oldEmail, newEmail string) {
	go func() {
		if len(oldEmail) == 0 {
			return
		}

		if len(newEmail) == 0 {
			newEmail = "无"
		}

		err := email.SendEmailChange(newEmail, oldEmail)
		if err != nil {
			logger.Logger.Error("send email error: %s", err.Error())
		}
	}()
}

func EmailSendDelete(userID int64) {
	go func() {
		e := getEmail(userID)
		if len(e) == 0 {
			return
		}

		err := email.SendDelete(e)
		if err != nil {
			logger.Logger.Error("send email error: %s", err.Error())
		}
	}()
}

func EmailSendBind(e string) {
	go func() {
		if len(e) == 0 {
			return
		}

		err := email.SendBind(e)
		if err != nil {
			logger.Logger.Error("send email error: %s", err.Error())
		}
	}()
}
