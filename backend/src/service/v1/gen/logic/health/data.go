package health

import (
	"gitee.com/wuntsong-auth/backend/src/ip"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"net/http"
)

func GetPingData(r *http.Request) types.PingData {
	realIp, ok := r.Context().Value("X-Real-IP").(string)
	if !ok {
		realIp = "未知"
	}

	geo, ok := r.Context().Value("X-Real-IP-Geo").(string)
	if !ok {
		geo = "未知"
	}

	geoCode, ok := r.Context().Value("X-Real-IP-Geo-Code").(string)
	if !ok {
		geoCode = ip.UnknownGeoCode
	}

	return types.PingData{
		RealIP:  realIp,
		Geo:     geo,
		GeoCode: geoCode,
	}
}
