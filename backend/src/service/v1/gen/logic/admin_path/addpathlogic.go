package admin_path

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
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

type AddPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPathLogic {
	return &AddPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPathLogic) AddPath(req *types.CreatePathReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !db.IsPathMode(req.Mode) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPathMode, "错误的路由模式"),
		}, nil
	}

	if !db.IsPathStatus(req.Mode) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPathStatus, "错误的路由状态"),
		}, nil
	}

	if !db.IsPathAdminMode(req.AdminMode) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPathMode, "错误的路由管理员模式"),
		}, nil
	}

	if !db.IsPathCorsModel(req.CorsMode) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPathMode, "错误的路由跨域模式"),
		}, nil
	}

	if !db.IsPathBusyMode(req.BusyMode) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPathMode, "错误的路由频繁限制模式"),
		}, nil
	}

	if !db.IsCaptchaMode(req.CaptchaMode) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPathMode, "错误的路由人机验证模式"),
		}, nil
	}

	var p big.Int
	for _, ps := range req.Policy {
		np, ok := model.PermissionsSign()[ps]
		if ok && np.Status != db.PolicyStatusBanned {
			p = permission.AddPermission(p, np.Permission)
		}
	}

	var sp int64
	for _, ps := range req.SubPolicy {
		np, ok := jwt.UserSubTokenStringMap[ps]
		if ok {
			sp = permission.AddPermissionInt64(sp, np)
		}
	}

	var mt int64
	for _, ps := range req.Method {
		np, ok := db.PathMethodStringMap[ps]
		if ok {
			mt = permission.AddPermissionInt64(mt, np)
		}
	}

	if req.Mode == db.PathModeRegex {
		_, err := regexp.Compile(req.Path)
		if err != nil {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.BadPathRegex, errors.WarpQuick(err), "正则表达式的路由错误"),
			}, nil
		}
	}

	pathModel := db.NewUrlPathModel(mysql.MySQLConn)

	count, err := pathModel.GetCount(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if count > config.BackendConfig.MySQL.SystemResourceLimit {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "超出限额"),
		}, nil
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

	_, err = pathModel.Insert(l.ctx, &db.UrlPath{
		Describe:       req.Describe,
		Path:           req.Path,
		Mode:           req.Mode,
		Status:         req.Status,
		IsOrPolicy:     req.IsOr,
		Permission:     p.Text(16),
		SubPolicy:      sp,
		Method:         mt,
		Authentication: req.Authentication,
		DoubleCheck:    req.DoubleCheck,
		CorsMode:       req.CorsMode,
		AdminMode:      req.AdminMode,
		BusyMode:       req.BusyMode,
		BusyCount:      req.BusyCount,
		CaptchaMode:    req.CaptchaMode,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.PathUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员新增路由（%s）成功", req.Path)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
