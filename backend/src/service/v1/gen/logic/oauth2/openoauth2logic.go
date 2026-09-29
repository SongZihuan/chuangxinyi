package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type OpenOauth2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenOauth2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenOauth2Logic {
	return &OpenOauth2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenOauth2Logic) OpenOauth2(req *types.OpenOauth2Req) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web := action.GetWebsiteByUID(req.WebID)
	if web.ID == warp.UnknownWebsite || web.ID == warp.UserCenterWebsite {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotFound, "外站未找到"),
		}, nil
	}

	bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
	_, err = bannedModel.InsertWithDelete(l.ctx, &db.Oauth2Baned{
		UserId:      user.Id,
		WebId:       web.ID,
		AllowLogin:  true,
		AllowDefray: true,
		AllowMsg:    true,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户开通站点（%s）", web.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
