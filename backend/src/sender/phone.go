package sender

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sms"
	"github.com/wuntsong-org/wterrors"
)

func getPhone(userID int64) string {
	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	phone, err := phoneModel.FindByUserID(context.Background(), userID)
	if errors.Is(err, db.ErrNotFound) {
		return ""
	} else if err != nil {
		return ""
	}

	return phone.Phone
}

func PhoneSendChange(userID int64, project string) {
	go func() {
		phone := getPhone(userID)
		if len(phone) == 0 {
			return
		}

		err := sms.SendChange(project, phone)
		if !errors.Is(err, sms.SMSSendLimit) && err != nil {
			logger.Logger.Error("send sms error: %s", err.Error())
		}
	}()
}

func PhoneSendPhoneChange(oldPhone, newPhone string) {
	go func() {
		if len(oldPhone) == 0 {
			return
		}

		err := sms.SendChangePhone(newPhone, oldPhone)
		if !errors.Is(err, sms.SMSSendLimit) && err != nil {
			logger.Logger.Error("send sms error: %s", err.Error())
		}
	}()
}

func PhoneSendDelete(userID int64) {
	go func() {
		phone := getPhone(userID)
		if len(phone) == 0 {
			return
		}

		err := sms.SendDelete(phone)
		if !errors.Is(err, sms.SMSSendLimit) && err != nil {
			logger.Logger.Error("send sms error: %s", err.Error())
		}
	}()
}

func PhoneSendBind(p string) {
	go func() {
		if len(p) == 0 {
			return
		}

		err := sms.SendRegister(p)
		if !errors.Is(err, sms.SMSSendLimit) && err != nil {
			logger.Logger.Error("send sms error: %s", err.Error())
		}
	}()
}
