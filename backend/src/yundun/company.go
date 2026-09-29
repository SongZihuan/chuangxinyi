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

func CheckCompany(name string, id string, legalPerson string) (res bool, resErr errors.WTError) {
	defer utils.Recover(logger.Logger, &resErr, "check company error")

	prefix := "https://cardnotwo.market.alicloudapi.com/company"

	params := url.Values{}
	params.Add("com", name)

	reqUrl := fmt.Sprintf("%s?%s", prefix, params.Encode())
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return false, errors.Warp(err, "fail to create http requests")
	}

	req.Header.Set("Authorization", "APPCODE "+config.BackendConfig.Aliyun.Identity.AppCode)

	client := &http.Client{}
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

	type Result struct {
		CompanyName string `json:"companyName"`
		CreditCode  string `json:"creditCode"`
		LegalPerson string `json:"faRen"`
	}

	type Body struct {
		Code    int    `json:"error_code"`
		Message string `json:"reason"`
		Result  Result `json:"result"`
	}

	data := Body{}
	err = utils.JsonUnmarshal(body, &data)
	if err != nil {
		return false, errors.Warp(err, "fail to unmarshal resp")
	}

	if data.Code == 0 {
		if data.Result.CompanyName == name && data.Result.CreditCode == id && data.Result.LegalPerson == legalPerson {
			return true, nil
		}
		return false, nil
	} else if data.Code == 50002 {
		return false, nil
	}

	return false, errors.Errorf("system error[%d]: %s", data.Code, data.Message)
}
