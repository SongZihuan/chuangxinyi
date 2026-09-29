package admin_website

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddWebsiteDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddWebsiteDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddWebsiteDomainLogic {
	return &AddWebsiteDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddWebsiteDomainLogic) AddWebsiteDomain(req *types.AddWebsiteDomainReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	if web.ID != warp.UserCenterWebsite && web.ID != req.WebsiteID {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "外站不允许操作"),
		}, nil
	}

	websiteDomainModel := db.NewWebsiteDomainModel(mysql.MySQLConn)

	count, err := websiteDomainModel.GetCount(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if count > config.BackendConfig.MySQL.SystemResourceLimit {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "超出限额"),
		}, nil
	}

	website := action.GetWebsite(req.WebsiteID)
	if website.ID == warp.UnknownWebsite {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadWebsiteID, "外站未找到"),
		}, nil
	}

	_, err = websiteDomainModel.Insert(l.ctx, &db.WebsiteDomain{
		WebsiteId: req.WebsiteID,
		Domain:    req.Domain,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.WebsiteUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员添加站点（%s）新域名（%s）成功", website.Name, req.Domain)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
