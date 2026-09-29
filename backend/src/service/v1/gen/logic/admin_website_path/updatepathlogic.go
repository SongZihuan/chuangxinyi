package admin_website_path

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"math/big"
	"regexp"
	"strings"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdatePathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePathLogic {
	return &UpdatePathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePathLogic) UpdatePath(req *types.UpdateWebsitePathReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !db.IsWebsitePathStatus(req.Status) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPathStatus, "错误的路由状态"),
		}, nil
	}

	if !db.IsWebsitePathMode(req.Mode) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPathMode, "错误的路由模式"),
		}, nil
	}

	var p big.Int
	for _, ps := range req.Policy {
		np, ok := (model.WebsitePermissionsSign())[ps]
		if ok && np.Status != db.WebsitePolicyStatusBanned {
			p = permission.AddPermission(p, np.Permission)
		}
	}

	var mt int64
	for _, ps := range req.Method {
		np, ok := db.PathMethodStringMap[ps]
		if ok {
			mt = permission.AddPermissionInt64(mt, np)
		}
	}

	pathModel := db.NewWebsiteUrlPathModel(mysql.MySQLConn)

	path, err := pathModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PathNotFound, "路由未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if req.Mode == db.WebsitePathModeRegex {
		_, err := regexp.Compile(req.Path)
		if err != nil {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.BadPathRegex, errors.WarpQuick(err), "正则表达式的路由错误"),
			}, nil
		}
	}

	if !strings.HasPrefix(req.Path, "/") {
		req.Path = "/" + req.Path
	}

	if !strings.HasPrefix(req.Path, "/api/v1/") { // "/api/v1/"，需要后缀的/
		req.Path = "/api/v1" + req.Path
	}

	if strings.HasSuffix(req.Path, "/") {
		req.Path = req.Path[0 : len(req.Path)-1]
	}

	path.Path = req.Path
	path.Describe = req.Describe
	path.Mode = req.Mode
	path.Status = req.Status
	path.IsOrPolicy = req.IsOr
	path.Permission = p.Text(16)
	path.Method = mt

	err = pathModel.Update(l.ctx, path)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.WebsitePathUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员更新站点路由（%s）单成功", path.Path)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
