package policycheck

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/handler/notallow"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/record"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/gorilla/websocket"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func WebsitePolicyCheck(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, isOptions bool) {
	// 不处理options请求，不允许options请求
	if isOptions || r.Method == http.MethodOptions {
		(notallow.NotAllow{}).ServeHTTP(w, r)
		return
	}

	if !IsWebsitePolicyCheck(r) {
		xDomainDeny(w, r)
		return
	}

	record := record.GetRecord(r.Context())

	domainUID := r.Header.Get("X-Domain")
	if len(domainUID) == 0 {
		domainUID = r.URL.Query().Get("xdomain")
	}

	n := r.Header.Get("X-N")
	if len(n) == 0 {
		n = r.URL.Query().Get("xn")
	}

	timestamp := r.Header.Get("X-Timestamp")
	if len(timestamp) == 0 {
		timestamp = r.URL.Query().Get("xtimestamp")
	}

	signBase64 := strings.TrimSpace(r.Header.Get("X-Sign"))
	if len(signBase64) == 0 {
		signBase64 = r.URL.Query().Get("xsign")
	}

	if len(timestamp) == 0 || len(n) != 18 || len(signBase64) == 0 || len(domainUID) == 0 {
		logger.Logger.Error("decode sign error")
		ipCheckDeny(w, r, errors.Errorf("empty query args"), "缺少必要参数")
		return
	}

	website := action.GetWebsiteByUID(domainUID)
	if website.Status == db.WebsiteStatusBanned || website.ID == warp.UserCenterWebsite {
		ipCheckDeny(w, r, errors.Errorf("website banned"), "站点错误")
		return
	}
	record.RequestWebsite = &website

	ip, ok := r.Context().Value("X-Real-IP").(string)
	if !ok {
		ipCheckDeny(w, r, errors.Errorf("get ip fail"), "获取ip失败")
		return
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		logger.Logger.Error("decode sign error")
		ipCheckDenyIfRelease(w, r, website, next, r.Context(), errors.WarpQuick(err), "转换时间戳错误")
		return
	}

	t := time.Unix(ts, 0)
	if time.Now().After(t.Add(time.Minute * 5)) {
		logger.Logger.Error("decode sign error")
		ipCheckDenyIfRelease(w, r, website, next, r.Context(), errors.Errorf("timeout"), "请求过期")
		return
	}

	if !utils.CheckIPInList(ip, website.GetWebsiteIPStringListType(), true, true, true) {
		ipCheckDenyIfRelease(w, r, website, next, r.Context(), errors.Errorf("bad ip"), "错误的IP")
		return
	}

	pubkeyByte, err := base64.StdEncoding.DecodeString(website.PubKey)
	if err != nil {
		logger.Logger.Error("read website pubkey from base4 error: %s", err.Error())
		ipCheckDenyIfRelease(w, r, website, next, r.Context(), errors.WarpQuick(err), "从base64读取公钥错误")
		return
	}

	pubkey, err := utils.ReadRsaPublicKey(pubkeyByte)
	if err != nil {
		logger.Logger.Error("read website pubkey error: %s", err.Error())
		ipCheckDenyIfRelease(w, r, website, next, r.Context(), errors.WarpQuick(err), "读取rsa公钥错误")
		return
	}

	var signText string
	if r.Method == http.MethodGet || r.Method == http.MethodHead || websocket.IsWebSocketUpgrade(r) {
		q := r.URL.Query()
		q.Del("xsign")
		signText = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", r.Method, domainUID, timestamp, n, r.URL.Path, q.Encode())
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.Error("read body error: %s", err.Error())
			ipCheckDenyIfRelease(w, r, website, next, r.Context(), errors.WarpQuick(err), "读取请求体错误")
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body)) // 塞回去
		signText = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", r.Method, domainUID, timestamp, n, r.URL.Path, string(body))
	}

	sign, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		logger.Logger.Error("decode sign error: %s", err.Error())
		ipCheckDenyIfRelease(w, r, website, next, r.Context(), errors.WarpQuick(err), "从base64读取签名错误")
		return
	}

	err = utils.VerifyRsaHash256Sign(signText, sign, pubkey)
	if err != nil {
		ipCheckDenyIfRelease(w, r, website, next, r.Context(), errors.WarpQuick(err), "签名验证错误")
		return
	}

	matcher, ok := urlpath.CheckWebsiteUrlPath(r.URL.Path, r.Method)
	if !ok || matcher == nil {
		if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
			ipCheckDeny(w, r, errors.Errorf("bad matcher"), "错误的路由")
			return
		} else {
			w.Header().Add("X-Path-Mather", "mather-root")
		}
	} else {
		w.Header().Add("X-Path-Mather", fmt.Sprintf("mathcer-%d", matcher.ID))
	}

	if matcher != nil && !urlpath.CheckWebsiteMatherPermission(matcher, website) {
		ipCheckDeny(w, r, errors.Errorf("not permision"), "没有权限访问")
		return
	}

	next(w, r.WithContext(context.WithValue(r.Context(), "X-Src-Website", website)))
}

func ipCheckDeny(w http.ResponseWriter, r *http.Request, err errors.WTError, reason string) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByErrorWithCode(r.Context(), respmsg.WebsiteDenyCode, respmsg.BadRequestsSender, errors.WarpQuick(err), reason),
	})
}

func ipCheckDenyIfRelease(w http.ResponseWriter, r *http.Request, web warp.Website, next http.HandlerFunc, ctx context.Context, err errors.WTError, reason string) {
	if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
		ipCheckDeny(w, r, err, reason)
	} else {
		next(w, r.WithContext(context.WithValue(ctx, "X-Src-Website", web)))
	}
}
