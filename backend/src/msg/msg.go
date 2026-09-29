package msg

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	errors "github.com/wuntsong-org/wterrors"
)

func SendMessage(userID int64, title string, content string, sender string, senderID int64, senderLink string) errors.WTError {
	messageModel := db.NewMessageModel(mysql.MySQLConn)
	_, err := messageModel.InsertCh(context.Background(), &db.Message{
		UserId:   userID,
		Title:    title,
		Content:  content,
		SenderId: senderID,
		Sender:   sender,
		SenderLink: sql.NullString{
			Valid:  len(senderLink) != 0,
			String: senderLink,
		},
	})
	if err != nil {
		return errors.WarpQuick(err)
	}
	return nil
}
