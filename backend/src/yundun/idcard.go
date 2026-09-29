package yundun

import (
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"net/url"
)

func CheckIDCard(name string, idnumber string) (res bool, resErr errors.WTError) {
	defer utils.Recover(logger.Logger, &resErr, "check id error")

	prefix := "https://id2meta.market.alicloudapi.com/id2meta"

	params := url.Values{}
	params.Add("identifyNum", idnumber)
	params.Add("userName", name)

	return sendReq(fmt.Sprintf("%s?%s", prefix, params.Encode()))
}

func CheckPhone(name string, idnumber string, phone string) (res bool, resErr errors.WTError) {
	defer utils.Recover(logger.Logger, &resErr, "check id error")

	prefix := "https://mobilecert.market.alicloudapi.com/mobile3Meta"

	params := url.Values{}
	params.Add("identifyNum", idnumber)
	params.Add("userName", name)
	params.Add("mobile", phone)

	return sendReq(fmt.Sprintf("%s?%s", prefix, params.Encode()))
}

func sendReq(httpsUrl string) (bool, errors.WTError) {
	req, err := http.NewRequest(http.MethodGet, httpsUrl, nil)
	if err != nil {
		return false, errors.Warp(err, "fail to create http requests")
	}

	client := &http.Client{}

	req.Header.Set("Authorization", "APPCODE "+config.BackendConfig.Aliyun.Identity.AppCode)

	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Warp(err, "fail to send http requests")
	}
	defer utils.Close(resp.Body)

	if resp.StatusCode != 200 {
		return false, errors.Errorf("fail to get http response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Warp(err, "fail to read resp")
	}

	type Data struct {
		BizCode string `json:"bizCode"`
	}

	type Body struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}

	data := Body{}
	err = utils.JsonUnmarshal(body, &data)
	if err != nil {
		return false, errors.Warp(err, "fail to unmarshal resp")
	}

	if data.Code == "200" {
		return data.Data.BizCode == "1", nil
	}

	return false, errors.Errorf("system error: %s", data.Message)
}
