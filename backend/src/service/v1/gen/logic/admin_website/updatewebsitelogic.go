package admin_website

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
	"math/big"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateWebsiteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWebsiteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWebsiteLogic {
	return &UpdateWebsiteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWebsiteLogic) UpdateWebsite(req *types.UpdateWebsiteReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	if web.ID != warp.UserCenterWebsite && web.ID != req.ID {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "外站不允许操作"),
		}, nil
	}

	if !db.IsWebsiteStatus(req.Status) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadWebsiteStatus, "错误的外站状态"),
		}, nil
	}

	websiteModel := db.NewWebsiteModel(mysql.MySQLConn)

	website, err := websiteModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadWebsiteID, "外站未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	keyMap := make(map[string]string, len(req.KeyMap))
	for _, v := range req.KeyMap {
		keyMap[v.Label] = v.Value
	}

	keyMapJson, err := utils.JsonMarshal(keyMap)
	if err != nil {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.BadKeyMap, errors.WarpQuick(err), "编码KeyMap错误"),
		}, nil
	}

	var p big.Int
	for _, ps := range req.Policy {
		np, ok := model.WebsitePermissionsSign()[ps]
		if ok && np.Status != db.WebsitePolicyStatusBanned {
			p = permission.AddPermission(p, np.Permission)
		}
	}

	if web.ID != warp.UserCenterWebsite && website.Permission != p.Text(16) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "外站不允许更新外站权限"),
		}, nil
	}

	website.Name = req.Name
	website.Describe = req.Describe
	website.Keymap = string(keyMapJson)
	website.Agreement = req.Agreement
	website.Permission = p.Text(16)
	website.Status = req.Status

	err = websiteModel.Update(l.ctx, website)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.WebsiteUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员更新站点（%s）", website.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
