package fuwuhao

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/rand"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"net/url"
	"time"
)

var Retry = errors.NewClass("retry")

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func getAccessTokenFromWeiXin(appID string, appSecret string) (string, int64, errors.WTError) {
	q := url.Values{}
	q.Add("grant_type", "client_credential")
	q.Add("appid", appID)
	q.Add("secret", appSecret)

	reqURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?%s", q.Encode())
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", 0, errors.WarpQuick(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, errors.WarpQuick(err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, errors.Errorf("bad respmsg status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, errors.WarpQuick(err)
	}

	data := AccessToken{}
	err = utils.JsonUnmarshal(body, &data)
	if err != nil {
		return "", 0, errors.WarpQuick(err)
	}

	if data.ErrCode == 0 {
		return data.AccessToken, data.ExpiresIn, nil
	}

	return "", 0, errors.Errorf("get access token respmsg[%d]: %s", data.ErrCode, data.ErrMsg)
}

func getAccessToken(ctx context.Context, appID string, appSecret string) (string, int64, errors.WTError) {
	key := fmt.Sprintf("fuwuhao:accesstoken:%s", appID)

	baseAccessToken, err := redis.Get(ctx, key).Result()
	if err == nil {
		return baseAccessToken, 0, nil
	}

	for i := 0; i < 10; i++ {
		newAccessToken, expires, err := func() (string, int64, errors.WTError) {
			var err error

			if !redis.AcquireLock(ctx, key, 2*time.Minute) {
				return "", 0, nil
			}
			defer redis.ReleaseLock(key)

			accessToken, err := redis.Get(ctx, key).Result()
			if err == nil && accessToken != "" {
				// access token已经刷新
				return accessToken, 0, nil
			}

			newAccessToken, expires, err := getAccessTokenFromWeiXin(appID, appSecret)
			if errors.Is(err, Retry) {
				return "", 0, nil
			} else if err != nil {
				return "", 0, errors.WarpQuick(err)
			}

			err = redis.Set(ctx, key, newAccessToken, time.Duration(expires)*time.Second).Err()
			if err != nil {
				return "", 0, errors.WarpQuick(err)
			}

			return newAccessToken, expires, nil
		}()
		if err != nil {
			return "", 0, errors.WarpQuick(err)
		} else if len(newAccessToken) == 0 {
			time.Sleep(time.Duration(30+rand.GlobalRander.Intn(20)) * time.Second)
			continue
		}
		return newAccessToken, expires, nil
	}

	return "", 0, errors.Errorf("can not get access token")
}

func delAccessToken(ctx context.Context, appID string) errors.WTError {
	key := fmt.Sprintf("fuwuhao:accesstoken:%s", appID)
	err := redis.Del(ctx, key).Err()
	if err != nil {
		return errors.WarpQuick(err)
	}
	return nil
}
