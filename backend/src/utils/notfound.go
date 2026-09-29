package utils

import (
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	"net/http"
)

type Resp struct {
	RequestsID string `json:"requestsID"`
	Code       int64  `json:"code"`
	Msg        string `json:"msg,omitempty"`
}

type RespEmpty struct {
	Resp
	Data struct{} `json:"data"`
}

// 和src/code.go对应

const StatusNotFound = 9
const StatusForbidden = 187

func NotFound(w http.ResponseWriter, r *http.Request, err error, showMsg bool, requestsID string) {
	var m string
	if showMsg && err != nil {
		m = err.Error()
	}

	httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, &RespEmpty{
		Resp: Resp{
			Code:       StatusNotFound,
			Msg:        m,
			RequestsID: requestsID,
		},
	})
}

func Forbidden(w http.ResponseWriter, r *http.Request, err error, msg bool, requestsID string) {
	var m string
	if msg && err != nil {
		m = err.Error()
	}

	httpx.WriteJsonCtx(r.Context(), w, http.StatusForbidden, &RespEmpty{
		Resp: Resp{
			Code:       StatusForbidden,
			Msg:        m,
			RequestsID: requestsID,
		},
	})
}
