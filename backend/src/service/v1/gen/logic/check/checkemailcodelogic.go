package check

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/checkcode"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
	"strconv"
	"strings"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CheckEmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckEmailCodeLogic {
	return &CheckEmailCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckEmailCodeLogic) CheckEmailCode(req *types.CheckEmailCodeReq, r *http.Request) (resp *types.SuccessResp, err error) {
	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	if web.ID == warp.UserCenterWebsite {
		if req.Type != EmailCheckToken && req.Type != UserToken {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTokenType, "错误的token类型"),
			}, nil
		}
	} else {
		if req.Type != EmailCheckToken {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTokenType, "错误的token类型"),
			}, nil
		}
	}

	if !utils.IsEmailAddress(req.Email) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadEmail, "错误的邮箱地址"),
		}, nil
	}

	key := fmt.Sprintf("code:email:%s", req.Email)
	if config.BackendConfig.GetMode() == config.RunModeDevelop && req.Code == "123456" {
		// 直接通过
	} else {
		res1, err := redis.Get(l.ctx, key).Result()
		if err != nil {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByErrorWithDebug(l.ctx, respmsg.OutOfDateCode, errors.WarpQuick(err), "验证码未找到", "验证码错误"),
			}, nil
		}

		res1Split := strings.Split(res1, ";")
		if len(res1Split) != 3 {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.OutOfDateCode, "验证码记录错误", "验证码错误"),
			}, nil
		}

		code := res1Split[0]
		t := res1Split[1]
		webID, err := strconv.ParseInt(res1Split[2], 10, 64)
		if err != nil {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByErrorWithDebug(l.ctx, respmsg.OutOfDateCode, errors.WarpQuick(err), "验证码webID读取错误", "验证码错误"),
			}, nil
		}

		if webID != web.ID {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.OutOfDateCode, "验证码webID不对应", "验证码错误"),
			}, nil
		}

		if t == checkcode.NormalCode && req.Type != EmailCheckToken {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadTokenType, "错误的token类型", "验证码错误"),
			}, nil
		}

		if code != req.Code {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadCode, "验证码验证失败", "验证码错误"),
			}, nil
		}
	}

	var tokenType string
	var subType string
	var token string
	switch req.Type {
	case EmailCheckToken:
		tokenType = EmailCheckToken
		token, err = jwt.CreateEmailToken(req.Email, web.ID)
		if err != nil {
			return nil, respmsg.JWTError.WarpQuick(err)
		}
	case UserToken:
		var email *db.Email
		emailModel := db.NewEmailModel(mysql.MySQLConn)

		email, err = emailModel.FindByEmail(l.ctx, req.Email)
		if errors.Is(err, db.ErrNotFound) {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.EmailNotBind, "邮箱未绑定用户", "邮箱未绑定用户"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		var user *db.User
		var secondfa *db.Secondfa
		userModel := db.NewUserModel(mysql.MySQLConn)
		secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
		ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

		user, err = userModel.FindOneByIDWithoutDelete(l.ctx, email.UserId)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return &types.SuccessResp{
					Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.EmailNotBind, "用户未找到", "邮箱未绑定用户"),
				}, nil
			}
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		ctrl, err := ctrlModel.FindByUserID(l.ctx, user.Id)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		if !ctrl.AllowEmail {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.UserNotFound, "用户不允许邮箱登录", "邮箱未绑定用户"),
			}, nil
		}

		secondfa, err = secondfaModel.FindByUserID(l.ctx, user.Id)
		if !errors.Is(err, db.ErrNotFound) && err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		} else if ctrl.Allow2Fa && secondfa != nil && secondfa.Secret.Valid { // secondfa != nil 表示未找到
			tokenType = SecondFAToken
			token, err = jwt.CreateLogin2FAToken(user.Uid)
			if err != nil {
				return nil, respmsg.JWTError.WarpQuick(err)
			}
			audit.NewUserAudit(user.Id, "用户试图通过邮件登录，需要2FA")
		} else {
			tokenType = UserToken
			sender.MessageSendLoginCenter(user.Id, l.ctx)
			sender.WxrobotSendLoginCenter(user.Id, l.ctx)
			sender.FuwuhaoSendLoginCenter(user.Id)

			if user.FatherId.Valid {
				token, err = jwt.CreateUserToken(l.ctx, user.Uid, user.Signin, user.TokenExpiration, jwt2.UserSonToken, "", 0)
				subType = jwt2.UserSonTokenString
			} else {
				token, err = jwt.CreateUserToken(l.ctx, user.Uid, user.Signin, user.TokenExpiration, jwt2.UserRootToken, "", 0)
				subType = jwt2.UserRootTokenString
			}
			if err != nil {
				return nil, respmsg.JWTError.WarpQuick(err)
			}
			audit.NewUserAudit(user.Id, "用户通过邮件登录成功")
		}
	}

	res2 := redis.Del(l.ctx, key)
	if err = res2.Err(); err != nil {
		return nil, respmsg.RedisSystemError.WarpQuick(err)
	}

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:    tokenType,
			Token:   token,
			SubType: subType,
		},
	}, nil
}
