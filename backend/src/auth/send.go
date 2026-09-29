package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type RespInterface interface {
	GetCode() string
	GetMsg() string
}

type AuthResp struct {
	Code string `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

func (r AuthResp) GetCode() string {
	return r.Code
}

func (r AuthResp) GetMsg() string {
	return r.Msg
}

func SendRequests(data any, u string, domainUID string, r RespInterface) (*http.Response, errors.WTError) {
	dataByte, jsonErr := utils.JsonMarshal(data)
	if jsonErr != nil {
		return nil, jsonErr
	}

	uu, err := url.Parse(u)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(dataByte))
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	n, err := utils.GenerateUniqueNumber(18)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	signText := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", http.MethodPost, domainUID, timestamp, n, uu.Path, string(dataByte))
	sign, err := utils.SignRsaHash256Sign(signText, PriKey)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Domain", domainUID)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-N", n)
	req.Header.Set("X-Sign", base64.StdEncoding.EncodeToString(sign))
	req.Header.Set("X-RunMode", "release")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("get bad status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	err = utils.JsonUnmarshal(body, r)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	if r.GetCode() != "SUCCESS" {
		return nil, errors.Errorf("requests fail: %s", r.GetMsg())
	}

	return resp, nil
}
