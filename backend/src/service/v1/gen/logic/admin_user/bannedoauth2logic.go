package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type BannedOauth2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBannedOauth2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *BannedOauth2Logic {
	return &BannedOauth2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BannedOauth2Logic) BannedOauth2(req *types.AdminBannedOauth2Req) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	web := action.GetWebsite(req.WebID)
	if web.ID == warp.UnknownWebsite || web.ID == warp.UserCenterWebsite {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotFound, "外站未找到"),
		}, nil
	}

	if !req.AllowLogin && (req.AllowDefray || req.AllowMsg) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadOauthAllow, "禁止登录即禁止消费和禁止通信"),
		}, nil
	}

	bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
	_, err = bannedModel.InsertWithDelete(l.ctx, &db.Oauth2Baned{
		UserId:      srcUser.Id,
		WebId:       web.ID,
		AllowLogin:  req.AllowLogin,
		AllowDefray: req.AllowLogin && req.AllowDefray,
		AllowMsg:    req.AllowLogin && req.AllowMsg,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if !req.AllowLogin {
		err := jwt.DeleteAllWebsiteLoginToken(l.ctx, srcUser.Uid, web.ID)
		if err != nil {
			return nil, respmsg.JWTError.WarpQuick(err)
		}
	}

	audit.NewAdminAudit(user.Id, "管理员修改用户（%s）站点策略", srcUser.Uid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
