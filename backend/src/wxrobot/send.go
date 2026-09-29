package wxrobot

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
)

type Msg struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
}

type Text struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
}

func Send(ctx context.Context, webhook string, text string, atAll bool, senderID int64, sender string) errors.WTError {
	wxrobotMessageModel := db.NewWxrobotMessageModel(mysql.MySQLConn)
	wxrobotMessage := &db.WxrobotMessage{
		Webhook:  webhook,
		Text:     fmt.Sprintf("【%s】%s", sender, text),
		AtAll:    atAll,
		Success:  true,
		SenderId: senderID,
	}

	if len(webhook) == 0 {
		return nil
	}

	t := Text{
		Content: text,
	}

	if atAll {
		t.MentionedList = append(t.MentionedList, "@all")
	}

	data := Msg{
		MsgType: "text",
		Text:    t,
	}
	dataByte, jsonErr := utils.JsonMarshal(data)
	if jsonErr != nil {
		return jsonErr
	}

	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(dataByte))
	if err != nil {
		return errors.WarpQuick(err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		wxrobotMessage.Success = false
		wxrobotMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: err.Error(),
		}
		_, _ = wxrobotMessageModel.Insert(ctx, wxrobotMessage)
		return errors.WarpQuick(err)
	}

	if resp.StatusCode != 200 {
		wxrobotMessage.Success = false
		wxrobotMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: fmt.Sprintf("get bad status code: %d", resp.StatusCode),
		}
		_, _ = wxrobotMessageModel.Insert(ctx, wxrobotMessage)
		return errors.Errorf("get bad status code")
	}

	_, _ = wxrobotMessageModel.Insert(ctx, wxrobotMessage)
	return nil
}
