package websocket

import (
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

type WSMessage struct {
	Code        string      `json:"code"`
	Data        interface{} `json:"data"`
	Secret      interface{} `json:"secret"`
	SecretToken []int       `json:"-"`
	WebID       int64       `json:"-"`
	HasClear    bool        `json:"-"`
}

func (w *WSMessage) ClearData(tokenType int) {
	if w.HasClear {
		return
	}

	if !func() bool {
		if w.SecretToken == nil {
			return true
		}

		for _, t := range w.SecretToken {
			if t == tokenType {
				return true
			}
		}

		return false
	}() {
		w.Secret = nil
	}

	w.WebID = 0
	w.SecretToken = nil
	w.HasClear = true
}

func (w *WSMessage) ToPeersMsg() ([]byte, errors.WTError) {
	if w.HasClear {
		return nil, errors.Errorf("has beed clear")
	}

	return utils.JsonMarshal(struct {
		Code        string      `json:"code"`
		Data        interface{} `json:"data"`
		Secret      interface{} `json:"secret"`
		SecretToken []int       `json:"secretToken"`
		WebID       int64       `json:"webID"`
		HasClear    bool        `json:"-"`
	}{Code: w.Code, Data: w.Data, Secret: w.Secret, SecretToken: w.SecretToken, WebID: w.WebID, HasClear: false})
}

type WSInMessage struct {
	Code string `json:"code"`
	Data string `json:"data"`
}

type WSWebsiteMessage struct {
	Code string `json:"code"`
	Data string `json:"data"`
}

type WSPeersMessage struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

func WriteMessage[T WSMessage | WSInMessage | WSWebsiteMessage | WSPeersMessage](c chan T, msg T) {
	go func() {
		defer func() {
			recover()
			// 不需要utils.Recover
		}()
		c <- msg
	}()
}

func WritePeersMessage(code string, data any, msg WSMessage) {
	go func() { // 因为需要上锁，不想阻塞
		dataByte, err := utils.JsonMarshal(data)
		if err != nil {
			return
		}

		msgByte, err := msg.ToPeersMsg()
		if err != nil {
			return
		}

		PeersConnMapMutex.Lock()
		defer PeersConnMapMutex.Unlock()

		for _, p := range PeersConnMap {
			go func(p chan WSPeersMessage) {
				p <- WSPeersMessage{
					Code:    code,
					Data:    string(dataByte),
					Message: string(msgByte),
				}
			}(p)
		}
	}()
}

func WriteJson(data any) string {
	s, err := utils.JsonMarshal(data)
	if err != nil {
		return err.Error()
	}
	return string(s)
}
