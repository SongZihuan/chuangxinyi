package policycheck

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/afs"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/defaultuid"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/ip"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/handler/notallow"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/record"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func UserPolicyCheckOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		(notallow.NotAllow{}).ServeHTTP(w, r)
		return
	}

	ctx := r.Context()

	if IsWebsitePolicyCheck(r) {
		(notallow.NotAllow{}).ServeHTTP(w, r)
		return
	}

	// 不设置默认值
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.BackendConfig.User.AllowMethod, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.BackendConfig.User.AllowHeader, ", "))

	matcher, ok := urlpath.CheckUrlPath(r.URL.Path, r.Method)
	if !ok || matcher == nil {
		if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
			(notallow.NotAllow{}).ServeHTTP(w, r)
			return
		} else {
			w.Header().Add("X-Path-Matcher", "matcher-root")
		}
	} else {
		w.Header().Add("X-Path-Matcher", fmt.Sprintf("matcher-%d", matcher.ID))
	}

	RemoteIP, ok := ctx.Value("X-Real-IP").(string)
	if !ok || len(RemoteIP) == 0 {
		(notallow.NotAllow{}).ServeHTTP(w, r)
		return
	}

	Geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok || len(Geo) == 0 {
		(notallow.NotAllow{}).ServeHTTP(w, r)
		return
	}

	if matcher != nil && !urlpath.CheckMatherPermissionOptions(matcher) {
		(notallow.NotAllow{}).ServeHTTP(w, r)
		return
	}

	if matcher != nil {
		switch matcher.CorsMode {
		default:
			(notallow.NotAllow{}).ServeHTTP(w, r)
			return
		case db.PathCorsAll:
			WriteOrigin("", w)
		case db.PathCorsCenter:
			origin := r.Header.Get("Origin")
			if IsConfigOrigin(origin) { // 同源和配置包内的请求体已经被添加
				WriteOrigin(origin, w)
			} else {
				(notallow.NotAllow{}).ServeHTTP(w, r)
				return
			}
		case db.PathCorsWebsite:
			origin := r.Header.Get("Origin")
			if IsConfigOrigin(origin) {
				WriteOrigin(origin, w)
			} else if IsAllWebOrigin(origin, nil) { // 自带web的设置
				WriteOrigin(origin, w)
			} else {
				(notallow.NotAllow{}).ServeHTTP(w, r)
				return
			}
		}
	} else if config.BackendConfig.GetModeFromRequests(r) == config.RunModeDevelop {
		WriteOrigin("", w)
	} else {
		(notallow.NotAllow{}).ServeHTTP(w, r)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(r.Context()),
	})
}

func UserPolicyCheck(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == http.MethodOptions {
		(notallow.NotAllow{}).ServeHTTP(w, r)
		return
	}

	ctx := r.Context()

	if IsWebsitePolicyCheck(r) {
		xDomainDeny(w, r)
		return
	}

	WriteOrigin("", w) // 添加默认值
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.BackendConfig.User.AllowMethod, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.BackendConfig.User.AllowHeader, ", "))

	matcher, ok := urlpath.CheckUrlPath(r.URL.Path, r.Method)
	if !ok || matcher == nil {
		if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
			policyDeny(w, r, errors.Errorf("not matcher"), "错误的路由")
			return
		} else {
			w.Header().Add("X-Path-Matcher", "matcher-root")
		}
	} else {
		w.Header().Add("X-Path-Matcher", fmt.Sprintf("matcher-%d", matcher.ID))
	}

	RemoteIP, ok := ctx.Value("X-Real-IP").(string)
	if !ok || len(RemoteIP) == 0 {
		httpx.OkJsonCtx(ctx, w, &types.RespEmpty{
			Resp: respmsg.GetRespByErrorWithCode(ctx, respmsg.PolicyDenyCode, respmsg.SystemError, respmsg.BadContextError.New("X-Real-IP")),
		})
		return
	}

	Geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok || len(Geo) == 0 {
		httpx.OkJsonCtx(ctx, w, &types.RespEmpty{
			Resp: respmsg.GetRespByErrorWithCode(ctx, respmsg.PolicyDenyCode, respmsg.SystemError, respmsg.BadContextError.New("X-Real-IP-Geo")),
		})
		return
	}

	var xToken string
	var user *db.User
	var role *warp.Role
	var web *warp.Website
	var father *db.User
	var subType int

	recordData := record.GetRecord(ctx)

	xToken = r.Header.Get("X-Token")
	if len(xToken) == 0 {
		xToken = r.URL.Query().Get("xtoken")
	}
	if len(xToken) != 0 { // 此处不是else if
		if !func() bool {
			recordData.UserToken = xToken

			userData, _, err := jwt.ParserUserToken(ctx, xToken)
			if err != nil {
				return false
			}

			_ = jwt.UpdateUserTokenGeo(ctx, xToken, true)

			if r.Method == http.MethodGet || r.Method == http.MethodHead {
				jwt.SetLogin(ctx, xToken, jwt.LoginGet)
			} else {
				jwt.SetLogin(ctx, xToken, jwt.LoginPost)
			}

			var mysqlErr error
			userModel := db.NewUserModel(mysql.MySQLConn)
			user, mysqlErr = userModel.FindOneByUidWithoutDelete(ctx, userData.UserID)
			if errors.Is(mysqlErr, db.ErrNotFound) {
				return false
			} else if mysqlErr != nil {
				return false
			}
			recordData.User = user

			if db.IsBanned(user) {
				return false
			}

			if user.Status == db.UserStatus_Register {
				return false
			}

			role = utils.GetPointer(action.GetRole(user.RoleId, user.IsAdmin))
			recordData.Role = role

			if userData.SubType == jwt2.UserWebsiteToken {
				web = utils.GetPointer(action.GetWebsite(userData.WebsiteID))
				if web.Status == db.WebsiteStatusBanned {
					return false
				}
				recordData.Website = web

				if web.ID != warp.UserCenterWebsite {
					bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
					allow, err := bannedModel.CheckAllow(ctx, user.Id, web.ID, db.AllowMsg)
					if err != nil {
						return false
					} else if !allow {
						return false
					}
				}
			} else if userData.WebsiteID != 0 {
				badUserDeny(w, r, "错误的用户Token，非外站授权token却webID不为0")
			} else {
				web = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
			}

			if userData.SubType == jwt2.UserFatherToken || userData.SubType == jwt2.UserRootFatherToken || userData.SubType == jwt2.UserUncleToken {
				if !user.FatherId.Valid {
					return false
				}

				father, mysqlErr = userModel.FindOneByUidWithoutDelete(ctx, userData.FatherID)
				if errors.Is(mysqlErr, db.ErrNotFound) {
					return false
				} else if mysqlErr != nil {
					return false
				}

				if db.IsBanned(father) {
					return false
				}

				if father.Status == db.UserStatus_Register {
					return false
				}
			} else if len(userData.FatherID) != 0 {
				badUserDeny(w, r, "错误的用户Token，非父亲授权token却fatherID不为空")
			} else {
				father = nil
			}

			if user.FatherId.Valid {
				realFather, mysqlErr := userModel.FindOneByIDWithoutDelete(ctx, user.FatherId.Int64)
				if errors.Is(mysqlErr, db.ErrNotFound) {
					return false
				} else if mysqlErr != nil {
					return false
				}

				if db.IsBanned(realFather) {
					go func() {
						_ = utils2.DeleteUser(user, realFather.Status, 1000) // 该用户也要注销
					}()
					return false
				}

				if realFather.Status == db.UserStatus_Register {
					return false
				}

			}

			subType = userData.SubType

			ctx = context.WithValue(ctx, "X-Token", xToken)
			ctx = context.WithValue(ctx, "X-Token-Father", father)
			ctx = context.WithValue(ctx, "X-Token-User", user)

			ctx = context.WithValue(ctx, "X-Token-Website", *web)
			ctx = context.WithValue(ctx, "X-Token-Type", subType)
			ctx = context.WithValue(ctx, "X-Token-Role", *role)
			return true
		}() {
			if matcher != nil && !matcher.Authentication { // options被允许忽略xtoken
				web = nil // 先不设置web 跨域的时候设置
				subType = jwt2.UserNotToken
				role = utils.GetPointer(action.GetRoleBySign(config.BackendConfig.Admin.AnonymousRole.RoleSign, false))

				ctx = context.WithValue(ctx, "X-Token-Type", subType)
				ctx = context.WithValue(ctx, "X-Token-Role", *role)
			} else if config.BackendConfig.GetModeFromRequests(r) == config.RunModeDevelop {
				web = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
				subType = jwt2.UserNotToken
				role = utils.GetPointer(action.GetRoleBySign(config.BackendConfig.Admin.AnonymousRole.RoleSign, false))

				ctx = context.WithValue(ctx, "X-Token-Type", subType)
				ctx = context.WithValue(ctx, "X-Token-Website", *web)
				ctx = context.WithValue(ctx, "X-Token-Role", *role)
			} else {
				badUserDeny(w, r, "没有鉴权")
				return
			}
		}
	} else if config.BackendConfig.GetModeWithHeaderTrue(r, "X-Default-User") && len(defaultuid.DefaultUID) != 0 {
		recordData.UserToken = "DefaultUser"

		_ = jwt.UpdateUserTokenGeo(ctx, xToken, true)

		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			jwt.SetLogin(ctx, xToken, jwt.LoginGet)
		} else {
			jwt.SetLogin(ctx, xToken, jwt.LoginPost)
		}

		var mysqlErr error
		userModel := db.NewUserModel(mysql.MySQLConn)
		user, mysqlErr = userModel.FindOneByUidWithoutDelete(ctx, defaultuid.DefaultUID)
		if errors.Is(mysqlErr, db.ErrNotFound) {
			badUserDeny(w, r, "用户在数据库中未找到")
			return
		} else if mysqlErr != nil {
			httpx.ErrorCtx(ctx, w, respmsg.MySQLSystemError.WarpQuick(mysqlErr))
			return
		}
		recordData.User = user

		if db.IsBanned(user) {
			badUserDeny(w, r, "用户在已被封禁")
			return
		}

		if user.Status == db.UserStatus_Register {
			badUserDeny(w, r, "用户未完成注册")
			return
		}

		role = utils.GetPointer(action.GetRole(user.RoleId, user.IsAdmin))
		recordData.Role = role

		web = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
		father = nil

		if user.FatherId.Valid {
			realFather, mysqlErr := userModel.FindOneByIDWithoutDelete(ctx, user.FatherId.Int64)
			if errors.Is(mysqlErr, db.ErrNotFound) {
				badUserDeny(w, r, "用户父亲未找到")
				return
			} else if mysqlErr != nil {
				httpx.ErrorCtx(ctx, w, respmsg.MySQLSystemError.WarpQuick(mysqlErr))
				return
			}

			if db.IsBanned(realFather) {
				go func() {
					_ = utils2.DeleteUser(user, realFather.Status, 1000) // 该用户也要注销
				}()
				badUserDeny(w, r, "用户在已被封禁")
				return
			}

			if realFather.Status == db.UserStatus_Register {
				badUserDeny(w, r, "用户父亲未完成注册")
				return
			}

		}

		subType = jwt2.UserRootToken

		ctx = context.WithValue(ctx, "X-Token", xToken)
		ctx = context.WithValue(ctx, "X-Token-Father", father)
		ctx = context.WithValue(ctx, "X-Token-User", user)

		ctx = context.WithValue(ctx, "X-Token-Website", *web)
		ctx = context.WithValue(ctx, "X-Token-Type", subType)
		ctx = context.WithValue(ctx, "X-Token-Role", *role)
	} else if matcher != nil && !matcher.Authentication { // options被允许忽略xtoken
		web = nil // 先不设置web 跨域的时候设置
		subType = jwt2.UserNotToken
		role = utils.GetPointer(action.GetRoleBySign(config.BackendConfig.Admin.AnonymousRole.RoleSign, false))

		ctx = context.WithValue(ctx, "X-Token-Type", subType)
		ctx = context.WithValue(ctx, "X-Token-Role", *role)
	} else if config.BackendConfig.GetModeFromRequests(r) == config.RunModeDevelop {
		web = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
		subType = jwt2.UserNotToken
		role = utils.GetPointer(action.GetRoleBySign(config.BackendConfig.Admin.AnonymousRole.RoleSign, false))

		ctx = context.WithValue(ctx, "X-Token-Type", subType)
		ctx = context.WithValue(ctx, "X-Token-Website", *web)
		ctx = context.WithValue(ctx, "X-Token-Role", *role)
	} else {
		notUserDeny(w, r, "没有鉴权")
		return
	}

	if matcher != nil {
		if !matcher.Authentication {
			// 使用匿名权限
			if !urlpath.CheckMatherPermission(matcher, action.GetRoleBySign(config.BackendConfig.Admin.AnonymousRole.RoleSign, false), jwt2.UserNotToken) {
				badUserDeny(w, r, "没访问权限")
				return
			}
		} else {
			if !urlpath.CheckMatherPermission(matcher, *role, subType) {
				badUserDeny(w, r, "没访问权限")
				return
			}
		}
	} else if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
		badUserDeny(w, r, "没访问权限")
		return
	}

	if matcher != nil {
		switch matcher.CorsMode {
		default:
			corsDeny(w, r)
			return
		case db.PathCorsAll:
			if web == nil {
				web = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
				ctx = context.WithValue(ctx, "X-Token-Website", *web)
			}

			WriteOrigin("", w)
		case db.PathCorsCenter:
			if web != nil && web.ID != warp.UserCenterWebsite {
				corsDeny(w, r)
				return
			}

			origin := r.Header.Get("Origin")
			if !IsConfigOrigin(origin) { // 同源和配置包内的请求体已经被添加
				corsDeny(w, r)
				return
			}

			WriteOrigin(origin, w)

			if web == nil {
				web = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
				ctx = context.WithValue(ctx, "X-Token-Website", *web)
			}
		case db.PathCorsWebsite:
			if web == nil {
				origin := r.Header.Get("Origin")
				if IsConfigOrigin(origin) {
					WriteOrigin(origin, w)

					web = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
					ctx = context.WithValue(ctx, "X-Token-Website", *web)
				} else if IsAllWebOrigin(origin, &web) { // 自带web的设置
					WriteOrigin(origin, w)
					ctx = context.WithValue(ctx, "X-Token-Website", *web) // 补充设置X-Token-Website
				} else {
					corsDeny(w, r)
					return
				}
			} else {
				origin := r.Header.Get("Origin")
				if web.ID == warp.UserCenterWebsite && IsConfigOrigin(origin) {
					WriteOrigin(origin, w)
				} else if web.ID != warp.UserCenterWebsite && IsWebOrigin(origin, *web) {
					WriteOrigin(origin, w)
				} else {
					corsDeny(w, r)
					return
				}
			}
		}
	} else if config.BackendConfig.GetModeFromRequests(r) == config.RunModeDevelop {
		WriteOrigin("", w)
	} else {
		corsDeny(w, r)
		return
	}

	if matcher != nil {
		if matcher.DoubleCheck && !config.BackendConfig.GetModeWithHeaderTrue(r, "X-Skip-Double-Check") {
			if user == nil {
				doubleCheckDeny(w, r, "需要二次验证")
				return
			}

			secondFAToken := r.Header.Get("X-2FA-Token")
			phoneToken := r.Header.Get("X-Phone-Token")
			emailToken := r.Header.Get("X-Email-Token")

			if len(secondFAToken) != 0 {
				uid, err := jwt.ParserCheck2FAToken(secondFAToken)
				if err != nil {
					doubleCheckDeny(w, r, "2FA token错误")
					return
				} else if uid.WebID != warp.UserCenterWebsite {
					doubleCheckDeny(w, r, "2FA token错误")
					return
				}

				if user.Uid != uid.UserID {
					doubleCheckDeny(w, r, "2FA token用户对不上")
					return
				}
			} else if len(phoneToken) != 0 {
				phoneNumber, err := jwt.ParserPhoneToken(phoneToken)
				if err != nil {
					doubleCheckDeny(w, r, "phone token错误")
					return
				} else if phoneNumber.WebID != warp.UserCenterWebsite {
					doubleCheckDeny(w, r, "phone token错误")
					return
				}

				phoneModel := db.NewPhoneModel(mysql.MySQLConn)
				phone, mysqlErr := phoneModel.FindByUserID(r.Context(), user.Id)
				if mysqlErr != nil {
					doubleCheckDeny(w, r, "手机对应的用户找不到")
					return
				}

				if phone.Phone != phoneNumber.Phone {
					doubleCheckDeny(w, r, "手机对应错误")
					return
				}
			} else if len(emailToken) != 0 {
				emailAddress, err := jwt.ParserEmailToken(emailToken)
				if err != nil {
					doubleCheckDeny(w, r, "email token错误")
					return
				} else if emailAddress.WebID != warp.UserCenterWebsite {
					doubleCheckDeny(w, r, "email token错误")
					return
				}

				emailModel := db.NewEmailModel(mysql.MySQLConn)
				email, mysqlErr := emailModel.FindByUserID(r.Context(), user.Id)
				if mysqlErr != nil {
					doubleCheckDeny(w, r, "邮箱对应的用户找不到")
					return
				}

				if email.Email.String != emailAddress.Email {
					doubleCheckDeny(w, r, "邮箱对应错误")
					return
				}
			} else {
				doubleCheckDeny(w, r, "没有二次验证")
				return
			}
		}
	} else if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
		doubleCheckDeny(w, r, "需要二次验证")
		return
	}

	var belongWeb *warp.Website

	if matcher != nil {
		switch matcher.AdminMode {
		case db.PathNotAdmin:
			belongWeb = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
		case db.PathWebsiteAdmin:
			if user == nil {
				websiteAdminNotAllowLogin(w, r)
				return
			}

			belongWeb = utils.GetPointer(action.GetWebsite(role.Belong))
			if belongWeb.Status == db.WebsiteStatusBanned {
				websiteAdminNotAllowLogin(w, r)
				return
			}
		case db.PathCenterAdmin:
			if user == nil {
				websiteAdminNotAllowLogin(w, r)
				return
			}

			if role.Belong != 0 || web.ID != role.Belong {
				websiteAdminNotAllowLogin(w, r)
				return
			}

			belongWeb = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
		default:
			websiteAdminNotAllowLogin(w, r)
			return
		}
	} else {
		belongWeb = utils.GetPointer(action.GetWebsite(warp.UserCenterWebsite))
	}

	ctx = context.WithValue(ctx, "X-Belong-Website", *belongWeb)

	if Geo == ip.LocalGeo {
		// 跳过检查
	} else if matcher != nil {
		if !config.BackendConfig.GetModeWithHeaderTrue(r, "X-Skip-Busy") {
			key := r.URL.Path
			switch matcher.BusyMode {
			case db.PathBusyModeUser:
				if user == nil {
					key = fmt.Sprintf("userbusy:%s:non", RemoteIP)
				} else {
					key = fmt.Sprintf("userbusy:%s:non", user.Uid)
				}
			case db.PathBusyModeIP:
				key = fmt.Sprintf("userbusy:%s:non", RemoteIP)
			default:
				busyDeny(w, r, errors.Errorf("too busy"), "操作频繁")
				return
			}
			times, err := redis.Incr(r.Context(), key).Result()
			if err != nil {
				busyDeny(w, r, errors.WarpQuick(err), "redis自增错误")
				return
			}

			if times > matcher.BusyCount {
				busyDeny(w, r, errors.Errorf("too busy"), "操作频繁")
				return
			}

			if times == 1 {
				_ = redis.Expire(r.Context(), key, time.Second*1).Err()
			} else { // 保险操作
				ttl, err := redis.TTL(r.Context(), key).Result()
				if err != nil || ttl > time.Second*1 || ttl == redis.KeepTTL {
					_ = redis.Expire(r.Context(), key, time.Second*1).Err()
				}
			}
		}
	} else if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
		busyDeny(w, r, errors.Errorf("too busy"), "操作频繁")
		return
	}

	if matcher != nil {
		if !config.BackendConfig.GetModeWithHeaderTrue(r, "X-Skip-AFS") {
			var slider bool
			var silence int64
			if config.BackendConfig.Aliyun.AFS.CAPTCHAStatus { // 跳过验证
				slider = CAPTCHA(r)
			} else {
				slider = true
			}

			if config.BackendConfig.Aliyun.AFS.SilenceCAPTCHAStatus { // 跳过验证
				silence = SilenceCAPTCHA(r)
			} else {
				silence = afs.Pass
			}

			switch matcher.CaptchaMode {
			case db.CaptchaModeNone:
				// 什么都不做
			case db.CaptchaModeSliderOnly:
				if !slider {
					robotDeny(w, r)
					return
				}
			case db.CaptchaModeSilenceOnly:
				if silence == afs.Banned {
					robotDeny(w, r)
					return
				} else if silence == afs.CheckAgain {
					robotSecondCheck(w, r)
					return
				}
			case db.CaptchaModeOn:
				if !slider && silence == afs.Banned {
					robotDeny(w, r)
					return
				} else if !slider && silence == afs.CheckAgain {
					robotSecondCheck(w, r)
					return
				}
			default:
				robotDeny(w, r)
				return
			}
		}
	} else if config.BackendConfig.GetModeFromRequests(r) != config.RunModeDevelop {
		robotDeny(w, r)
		return
	}

	next(w, r.WithContext(ctx))
}

func busyDeny(w http.ResponseWriter, r *http.Request, err errors.WTError, args ...any) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByErrorWithCode(r.Context(), respmsg.PolicyDenyCode, respmsg.RequestsTooBusy, errors.WarpQuick(err), args...),
	})
}

func corsDeny(w http.ResponseWriter, r *http.Request) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.CORSDenyCode, respmsg.CorsNotAllow, "跨域"),
	})
}

func websiteAdminNotAllowLogin(w http.ResponseWriter, r *http.Request) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.TokenDenyCode, respmsg.WebsiteNotAllow, "路由禁止外站登录"),
	})
}

func websiteNotAllowLogin(w http.ResponseWriter, r *http.Request) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.TokenDenyCode, respmsg.WebsiteNotAllow, "用户禁止外站登录"),
	})
}

func badUserDenyWithError(w http.ResponseWriter, r *http.Request, err errors.WTError, args ...any) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByErrorWithCode(r.Context(), respmsg.TokenDenyCode, respmsg.BadUserToken, errors.WarpQuick(err), args...),
	})
}

func badUserDeny(w http.ResponseWriter, r *http.Request, m string, args ...any) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.TokenDenyCode, respmsg.BadUserToken, m, args...),
	})
}

func notUserDeny(w http.ResponseWriter, r *http.Request, m string, args ...any) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.NotTokenDenyCode, respmsg.BadUserToken, m, args...),
	})
}

func doubleCheckDeny(w http.ResponseWriter, r *http.Request, msg string, args ...any) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.DoubleCheckDenyCode, respmsg.PolicyDeny, msg, args...),
	})
}

func policyDeny(w http.ResponseWriter, r *http.Request, err errors.WTError, reason string) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByErrorWithCode(r.Context(), respmsg.PolicyDenyCode, respmsg.PolicyDeny, errors.WarpQuick(err), reason),
	})
}

func xDomainDeny(w http.ResponseWriter, r *http.Request) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.WebsiteDenyCode, respmsg.PolicyDeny, "错误的参数：xdomain"),
	})
}

func robotDeny(w http.ResponseWriter, r *http.Request) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.RobotDenyCode, respmsg.RobotDeny, "人机验证失败"),
	})
}

func robotSecondCheck(w http.ResponseWriter, r *http.Request) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByMsgWithCode(r.Context(), respmsg.RobotDenyCode, respmsg.CAPTCHASecondCheck, "人机验证触发二次验证"),
	})
}

func IsWebOrigin(origin string, web warp.Website) bool {
	if len(origin) == 0 {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	hostname := u.Hostname()
	for _, d := range web.Domain {
		if d.Domain == hostname {
			return true
		}
	}

	return false
}

func WriteOrigin(origin string, w http.ResponseWriter) {
	if len(origin) == 0 {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
}

func IsConfigOrigin(origin string) bool {
	if len(origin) == 0 { // 同源请求
		return true
	}

	for _, o := range config.BackendConfig.User.Origin {
		if origin == o {
			return true
		}
	}

	return false
}

func IsAllWebOrigin(origin string, webVal **warp.Website) bool {
	if len(origin) == 0 {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	hostname := u.Hostname()

	for _, web := range model.WebsiteList() {
		if web.Status == db.WebsiteStatusBanned {
			continue
		}
		for _, d := range web.Domain {
			if d.Domain == hostname {
				if webVal != nil {
					*webVal = utils.GetPointer(web)
				}
				return true
			}
		}
	}

	return false
}

func CAPTCHA(r *http.Request) bool {
	token := r.Header.Get("X-CAPTCHA-Token")
	sig := r.Header.Get("X-CAPTCHA-Sig")
	sessionId := r.Header.Get("X-CAPTCHA-SessionId")
	scene := r.Header.Get("X-CAPTCHA-Scene")
	remoteIP, ok := r.Context().Value("X-Real-IP").(string)
	if !ok {
		return false
	}

	if len(token) == 0 || len(sig) == 0 || len(sessionId) == 0 || len(scene) == 0 {
		return false
	}

	res := afs.CheckCAPTCHA(sessionId, token, sig, scene, config.BackendConfig.Aliyun.AFS.CAPTCHAAppKey, remoteIP)
	if !res {
		return false
	}

	return true
}

func SilenceCAPTCHA(r *http.Request) int64 {
	nvc := r.Header.Get("X-CAPTCHA-Nvc")

	if len(nvc) == 0 {
		return afs.Banned
	}

	res := afs.CheckSilenceCAPTCHA(nvc)
	if res == afs.Banned {
		return afs.Banned
	} else if res == afs.CheckAgain {
		return afs.CheckAgain
	}

	return afs.Pass
}
