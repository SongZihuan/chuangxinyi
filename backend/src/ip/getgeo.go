package ip

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const EmptyCityCode = "000000"
const EmptyCountryCode = "000"
const DefaultCountryGeoCode = "156"
const UnknownGeoCode = EmptyCountryCode + EmptyCityCode
const DefaultGeoCode = DefaultCountryGeoCode + EmptyCityCode
const GeoCodeLen = 9
const CityGeoCodeLen = 6
const CountryGeoCodeLen = GeoCodeLen - CityGeoCodeLen
const LocalGeo = "内网"

func GetGeo(ctx context.Context, ip string) (code string, geo string, resErr errors.WTError) {
	defer utils.Recover(logger.Logger, &resErr, "get geo error")

	if utils.IsLocalIP(ip) {
		return UnknownGeoCode, LocalGeo, nil
	}

	res, ok := redis.GetCache(ctx, fmt.Sprintf("ip:%s", url.QueryEscape(ip)))
	if ok {
		codegeo := strings.Split(res, ";")
		if len(codegeo) == 2 && len(codegeo[0]) == GeoCodeLen {
			return codegeo[0], codegeo[1], nil
		}
	}

	prefix := "https://qryip.market.alicloudapi.com/lundear/qryip"

	params := url.Values{}
	params.Add("ip", ip)

	reqUrl := fmt.Sprintf("%s?%s", prefix, params.Encode())
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return "", "", errors.Warp(err, "fail to create http requests")
	}

	req.Header.Set("Authorization", "APPCODE "+config.BackendConfig.Aliyun.IP.AppCode)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", errors.Warp(err, "fail to send http requests")
	}
	defer utils.Close(resp.Body)

	if resp.StatusCode != 200 {
		return "", "", errors.Errorf("fail to get http response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", errors.Warp(err, "fail to read resp")
	}

	type ADInfo struct {
		Nation     string `json:"nation"`
		NationCode int64  `json:"nation_code"`
		Province   string `json:"province"`
		City       string `json:"city"`
		District   string `json:"district"`
		ADCode     int64  `json:"adcode"`
	}

	type Result struct {
		ADInfo ADInfo `json:"ad_info"`
	}

	type Body struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Result  Result `json:"result"`
	}

	data := Body{}
	err = utils.JsonUnmarshal(body, &data)
	if err != nil {
		return "", "", errors.Warp(err, "fail to unmarshal resp")
	}

	if data.Status == 0 {
		if data.Result.ADInfo.NationCode == 0 {
			data.Result.ADInfo.NationCode = 156
		}

		if data.Result.ADInfo.ADCode <= 0 {
			data.Result.ADInfo.ADCode = 0
		}

		code = fmt.Sprintf("%03d%06d", data.Result.ADInfo.NationCode, data.Result.ADInfo.ADCode)
		geo = fmt.Sprintf("%s%s%s%s", data.Result.ADInfo.Nation, data.Result.ADInfo.Province, data.Result.ADInfo.City, data.Result.ADInfo.District)

		if len(code) != 9 {
			code = DefaultGeoCode
		}

		if strings.HasPrefix(code, "15671") { // 台湾
			code = "158" + EmptyCityCode
		} else if strings.HasPrefix(code, "15681") { // 香港
			code = "344" + EmptyCityCode
		} else if strings.HasPrefix(code, "15682") { // 澳门
			code = "446" + EmptyCityCode
		}

		if config.BackendConfig.Aliyun.IP.ExpiresSecond != 0 {
			redis.SetCache(ctx, fmt.Sprintf("ip:%s", url.QueryEscape(ip)), fmt.Sprintf("%s;%s", code, geo), time.Second*time.Duration(config.BackendConfig.Aliyun.IP.ExpiresSecond))
		}

		return code, geo, nil
	} else if data.Status == 382 {
		if config.BackendConfig.Aliyun.IP.ExpiresSecond != 0 {
			redis.SetCache(ctx, fmt.Sprintf("ip:%s", url.QueryEscape(ip)), fmt.Sprintf("%s;未知", UnknownGeoCode), time.Second*time.Duration(config.BackendConfig.Aliyun.IP.ExpiresSecond))
		}

		return UnknownGeoCode, "未知", nil
	} else if data.Status == 375 {
		if config.BackendConfig.Aliyun.IP.ExpiresSecond != 0 {
			redis.SetCache(ctx, fmt.Sprintf("ip:%s", url.QueryEscape(ip)), fmt.Sprintf("%s;%s", UnknownGeoCode, LocalGeo), time.Second*time.Duration(config.BackendConfig.Aliyun.IP.ExpiresSecond))
		}

		return UnknownGeoCode, LocalGeo, nil
	}

	return "", "", errors.Errorf("system error[%d]: %s", data.Status, data.Message)
}
