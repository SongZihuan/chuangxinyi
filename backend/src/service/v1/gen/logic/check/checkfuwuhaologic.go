package check

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"github.com/fastwego/offiaccount/apis/oauth"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type CheckFuwuhaoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckFuwuhaoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckFuwuhaoLogic {
	return &CheckFuwuhaoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckFuwuhaoLogic) CheckFuwuhao(req *types.CheckFuwuhaoCodeReq) (resp *types.SuccessResp, err error) {
	// 不允许外站
	if req.Type != FuwuhaoToken && req.Type != UserToken && req.Type != AutoToken {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadTokenType, "错误的Token类型"),
		}, nil
	}

	access, err := oauth.GetAccessToken(config.BackendConfig.FuWuHao.AppID, config.BackendConfig.FuWuHao.Secret, req.Code)
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByErrorWithDebug(l.ctx, respmsg.BadWeChatCode, errors.WarpQuick(err), "获取access token失败", "登录失败"),
		}, nil
	}

	userInfo, err := oauth.GetUserInfo(access.AccessToken, access.Openid, "zh_CN")
	if err != nil {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByErrorWithDebug(l.ctx, respmsg.BadWeChatCode, errors.WarpQuick(err), "获取用户信息失败", "登录失败"),
		}, nil
	} else if len(userInfo.Unionid) == 0 {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadWeChatCode, "获取unionID失败", "登录失败"),
		}, nil
	}

	var tokenType string
	var token string
	var subType string
	switch req.Type {
	case FuwuhaoToken:
		tokenType = FuwuhaoToken
		token, err = jwt.CreateWeChatToken(access.AccessToken, access.Openid, userInfo.Unionid, true)
	case UserToken, AutoToken:
		var wc *db.Wechat
		wechatModel := db.NewWechatModel(mysql.MySQLConn)

		wc, err := wechatModel.FindByUnionID(l.ctx, userInfo.Unionid)
		if errors.Is(err, db.ErrNotFound) {
			if req.Type == AutoToken {
				tokenType = FuwuhaoToken
				token, err = jwt.CreateWeChatToken(access.AccessToken, access.Openid, userInfo.Unionid, true)
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
			userModel := db.NewUserModel(mysql.MySQLConn)
			secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
			ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

			go func(wc *db.Wechat) {
				wechatModel := db.NewWechatModel(mysql.MySQLConn)
				wc.Nickname = sql.NullString{
					Valid:  len(userInfo.Nickname) != 0,
					String: userInfo.Nickname,
				}
				wc.Headimgurl = sql.NullString{
					Valid:  len(userInfo.Headimgurl) != 0,
					String: userInfo.Headimgurl,
				}
				wc.Fuwuhao = sql.NullString{
					Valid:  true,
					String: userInfo.Openid,
				}
				err := wechatModel.UpdateCh(context.Background(), wc)
				if err != nil {
					logger.Logger.Error("mysql error: %s", err.Error())
				}
			}(wc)

			user, err := userModel.FindOneByIDWithoutDelete(l.ctx, wc.UserId)
			if errors.Is(err, db.ErrNotFound) {
				return &types.SuccessResp{
					Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.UserNotFound, "用户未找到", "微信未绑定任何用户"),
				}, nil
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			ctrl, err := ctrlModel.FindByUserID(l.ctx, user.Id)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			if !ctrl.AllowWechat {
				return &types.SuccessResp{
					Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.UserNotFound, "用户不允许微信登录", "微信未绑定任何用户"),
				}, nil
			}

			secondfa, err := secondfaModel.FindByUserID(l.ctx, user.Id)
			if !errors.Is(err, db.ErrNotFound) && err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			} else if ctrl.Allow2Fa && secondfa != nil && secondfa.Secret.Valid { // secondfa != nil 表示未找到
				tokenType = SecondFAToken
				token, err = jwt.CreateLogin2FAToken(user.Uid)
				if err != nil {
					return nil, respmsg.JWTError.WarpQuick(err)
				}
				audit.NewUserAudit(user.Id, "用户试图通过服务号登录，需要2FA")
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
				audit.NewUserAudit(user.Id, "用户通过服务号登录成功")
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
