package wechat

import (
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"net/url"
)

type WeChatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Headimgurl string `json:"headimgurl"`
}

func GetUserInfo(accessToken, openID string) (WeChatUserInfo, errors.WTError) {
	prefix := "https://api.weixin.qq.com/sns/userinfo"

	params := url.Values{}
	params.Add("access_token", accessToken)
	params.Add("openid", openID)

	reqUrl := fmt.Sprintf("%s?%s", prefix, params.Encode())
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return WeChatUserInfo{}, errors.Warp(err, "fail to create http requests")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return WeChatUserInfo{}, errors.Warp(err, "fail to send http requests")
	}
	defer utils.Close(resp.Body)

	if resp.StatusCode != 200 {
		return WeChatUserInfo{}, errors.Errorf("fail to get http response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeChatUserInfo{}, errors.Warp(err, "fail to read resp")
	}

	var data WeChatUserInfo
	err = utils.JsonUnmarshal(body, &data)
	if err != nil {
		return WeChatUserInfo{}, errors.Warp(err, "fail to unmarshal resp")
	}

	if len(data.OpenID) == 0 {
		return WeChatUserInfo{}, errors.Errorf("bad access token")
	}

	return data, nil
}
