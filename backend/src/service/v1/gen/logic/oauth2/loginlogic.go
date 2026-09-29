package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	jwt2 "gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/ip"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.SuccessResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web := action.GetWebsiteByUID(req.DomainUID)
	if web.Status == db.WebsiteStatusBanned {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.DomainNotFound, "外站未找到"),
		}, nil
	}

	bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
	recordModel := db.NewOauth2RecordModel(mysql.MySQLConn)
	allow, err := bannedModel.CheckAllow(l.ctx, user.Id, web.ID, db.AllowLogin)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if !allow {
		return &types.SuccessResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.NotOpenWebsite, "用户未开通该网站"),
		}, nil
	}

	userToken, err := jwt.CreateUserToken(l.ctx, user.Uid, false, user.TokenExpiration, jwt2.UserWebsiteToken, "", web.ID)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	token, err := jwt.CreateLoginToken(l.ctx, user.Uid, web.ID, userToken)
	if err != nil {
		return nil, respmsg.JWTError.WarpQuick(err)
	}

	RemoteIP, ok := l.ctx.Value("X-Real-IP").(string)
	if !ok || len(RemoteIP) == 0 {
		RemoteIP = "0.0.0.0"
	}

	Geo, ok := l.ctx.Value("X-Real-IP-Geo").(string)
	if !ok || len(Geo) == 0 {
		Geo = "未知"
	}

	GeoCode, ok := l.ctx.Value("X-Real-IP-Geo-Code").(string)
	if !ok || len(Geo) == 0 {
		GeoCode = ip.UnknownGeoCode
	}

	_, err = recordModel.Insert(l.ctx, &db.Oauth2Record{
		UserId:    user.Id,
		WebId:     web.ID,
		WebName:   web.Name,
		Ip:        RemoteIP,
		Geo:       Geo,
		GeoCode:   GeoCode,
		LoginTime: time.Now(),
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.MessageSendOauth2(user.Id, web.Name, l.ctx)
	sender.WxrobotSendOauth2(user.Id, web.Name, l.ctx)
	sender.FuwuhaoSendOauth2(user.Id, web.Name)
	audit.NewUserAudit(user.Id, "用户授权登录网站：%s", web.Name)

	return &types.SuccessResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SuccessData{
			Type:     LoginToken,
			Token:    token,
			SubToken: userToken,
		},
	}, nil
}
