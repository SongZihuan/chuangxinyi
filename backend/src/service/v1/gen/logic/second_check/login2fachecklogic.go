package second_check

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/ip"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type Login2FACheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogin2FACheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Login2FACheckLogic {
	return &Login2FACheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Login2FACheckLogic) Login2FACheck(req *types.Login2FACheckReq, r *http.Request) (resp *types.SuccessResp, err error) {
	userID, err := jwt.ParserLogin2FAToken(req.LoginToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
	ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

	user, err := userModel.FindOneByUidWithoutDelete(l.ctx, userID.UserID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &types.SuccessResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		}
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if db.IsBanned(user) {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	}

	ctrl, err := ctrlModel.FindByUserID(l.ctx, user.Id)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	var passToken string
	if ctrl.Allow2Fa {
		secondfa, err := secondfaModel.FindByUserID(l.ctx, user.Id)
		if errors.Is(err, db.ErrNotFound) {
			// 不做操作
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		} else if secondfa.Secret.Valid {
			if config.BackendConfig.GetMode() == config.RunModeDevelop && req.Code == "123456" {
				// 直接通过
			} else {
				ua := r.Header.Get("User-Agent")
				geo, ok := l.ctx.Value("X-Real-IP-Geo").(string)
				if !ok {
					return nil, respmsg.BadContextError.New("X-Real-IP-Geo")
				}

				geoCode, ok := l.ctx.Value("X-Real-IP-Geo-Code").(string)
				if !ok {
					return nil, respmsg.BadContextError.New("X-Real-IP-Geo-Code")
				}

				if len(req.PassToken) != 0 {
					if len(ua) == 0 {
						return &types.SuccessResp{
							Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadPassToken, "缺少User-Agent", "错误的通过token"),
						}, nil
					}

					if geo != ip.LocalGeo && geoCode == ip.UnknownGeoCode {
						return &types.SuccessResp{
							Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadPassToken, "地址未知", "错误的通过token"),
						}, nil
					}

					data, err := jwt.ParserSecondFAPassToken(req.PassToken)
					if err != nil {
						return &types.SuccessResp{
							Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadPassToken, "PassToken解析错误", "错误的通过token"),
						}, nil
					}

					if data.UserID != user.Uid {
						return &types.SuccessResp{
							Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadPassToken, "用户不匹配", "错误的通过token"),
						}, nil
					}

					if data.UA != ua {
						return &types.SuccessResp{
							Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadPassToken, "UA不匹配", "错误的通过token"),
						}, nil
					}

					if data.GeoCode != geoCode {
						return &types.SuccessResp{
							Resp: respmsg.GetRespByMsgWithDebug(l.ctx, respmsg.BadPassToken, "地址不匹配", "错误的通过token"),
						}, nil
					}
				} else {
					if !utils.CheckTOTP(secondfa.Secret.String, req.Code) {
						return &types.SuccessResp{
							Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Bad2FACode, "错误的2FA密钥"),
						}, nil
					}

					if req.RememberHour > 0 {
						passToken, err = jwt.CreateSecondFAPassToken(user.Uid, ua, geoCode, req.RememberHour)
						if err != nil {
							return nil, respmsg.JWTError.WarpQuick(err)
						}
					}
				}
			}
		}
	}

	var token string
	var subType string
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

	sender.MessageSendLoginCenter(user.Id, l.ctx)
	sender.WxrobotSendLoginCenter(user.Id, l.ctx)
	sender.FuwuhaoSendLoginCenter(user.Id)
	audit.NewUserAudit(user.Id, "用户2FA验证成功，登录")

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:     UserToken,
			Token:    token,
			SubType:  subType,
			SubToken: passToken,
		},
	}, nil
}
