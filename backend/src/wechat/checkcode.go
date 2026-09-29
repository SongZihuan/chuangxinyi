package wechat

import (
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"net/url"
)

func CheckCode(code string) (string, string, string, errors.WTError) {
	prefix := "https://api.weixin.qq.com/sns/oauth2/access_token"

	params := url.Values{}
	params.Add("appid", config.BackendConfig.WeChat.AppID)
	params.Add("secret", config.BackendConfig.WeChat.AppSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	reqUrl := fmt.Sprintf("%s?%s", prefix, params.Encode())
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return "", "", "", errors.Warp(err, "fail to create http requests")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", errors.Warp(err, "fail to send http requests")
	}
	defer utils.Close(resp.Body)

	if resp.StatusCode != 200 {
		return "", "", "", errors.Errorf("fail to get http response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", errors.Warp(err, "fail to read resp")
	}

	type Body struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		UnionID     string `json:"unionid"`
	}

	var data Body
	err = utils.JsonUnmarshal(body, &data)
	if err != nil {
		return "", "", "", errors.Warp(err, "fail to unmarshal resp")
	}

	if len(data.AccessToken) == 0 {
		return "", "", "", errors.Errorf("bad code")
	}

	return data.AccessToken, data.OpenID, data.UnionID, nil
}
