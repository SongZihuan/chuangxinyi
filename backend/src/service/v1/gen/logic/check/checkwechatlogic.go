package check

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/wechat"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CheckWechatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckWechatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckWechatLogic {
	return &CheckWechatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckWechatLogic) CheckWechat(req *types.CheckWechatCodeReq) (resp *types.SuccessResp, err error) {
	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	if web.ID == warp.UserCenterWebsite {
		if req.Type != WeChatToken && req.Type != UserToken && req.Type != AutoToken {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTokenType, "错误的Token类型"),
			}, nil
		}
	} else {
		if req.Type != WeChatToken {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTokenType, "错误的Token类型"),
			}, nil
		}
	}

	accessToken, openID, unionID, err := wechat.CheckCode(req.Code)
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByErrorWithDebug(l.ctx, respmsg.BadWeChatCode, errors.WarpQuick(err), "获取access token失败", "微信未绑定"),
		}, nil
	}

	var tokenType string
	var token string
	var subType string
	switch req.Type {
	case WeChatToken:
		tokenType = WeChatToken
		token, err = jwt.CreateWeChatToken(accessToken, openID, unionID, false)
	case UserToken, AutoToken:
		var wc *db.Wechat
		wechatModel := db.NewWechatModel(mysql.MySQLConn)

		wc, err = wechatModel.FindByUnionID(l.ctx, unionID)
		if errors.Is(err, db.ErrNotFound) {
			if req.Type == AutoToken {
				tokenType = WeChatToken
				token, err = jwt.CreateWeChatToken(accessToken, openID, unionID, false)
				if err != nil {
					return nil, respmsg.JWTError.WarpQuick(err)
				}
			} else {
				return &types.SuccessResp{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WeChatNotBind, "微信未绑定任何用户"),
				}, nil
			}
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		} else {
			var user *db.User
			var secondfa *db.Secondfa
			userModel := db.NewUserModel(mysql.MySQLConn)
			secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
			ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

			user, err = userModel.FindOneByIDWithoutDelete(l.ctx, wc.UserId)
			if errors.Is(err, db.ErrNotFound) {
				return &types.SuccessResp{
					Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.UserNotFound, "用户未找到", "微信未绑定"),
				}, nil
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			go func(wc *db.Wechat, accessToken, openID string, wechatModel db.WechatModel) {
				userInfo, err := wechat.GetUserInfo(accessToken, openID)
				if err == nil { // 无错误是进行操作
					wc.Nickname = sql.NullString{
						Valid:  len(userInfo.Nickname) != 0,
						String: userInfo.Nickname,
					}
					wc.Headimgurl = sql.NullString{
						Valid:  len(userInfo.Headimgurl) != 0,
						String: userInfo.Headimgurl,
					}
					wc.OpenId = sql.NullString{
						Valid:  true,
						String: openID,
					}
					mysqlErr := wechatModel.UpdateCh(context.Background(), wc)
					if mysqlErr != nil {
						logger.Logger.Error("mysql error: %s", mysqlErr.Error())
					}
				}
			}(wc, accessToken, openID, wechatModel)

			ctrl, err := ctrlModel.FindByUserID(l.ctx, user.Id)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			if !ctrl.AllowWechat {
				return &types.SuccessResp{
					Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.UserNotFound, "用户不允许微信登录", "微信未绑定"),
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

				audit.NewUserAudit(user.Id, "用户试图通过微信登录，需要2FA")
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

				audit.NewUserAudit(user.Id, "用户通过微信登录成功")
			}
		}
	default:
		return nil, errors.Errorf("bad swtich case")
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
