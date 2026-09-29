package admin_role

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
	"math/big"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type RoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateLogic {
	return &RoleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleUpdateLogic) RoleUpdate(req *types.UpdateRoleReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	roleModel := db.NewRoleModel(mysql.MySQLConn)
	role, err := roleModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.RoleNotFound, "角色未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if role.NotChangeSign && role.Sign != req.Sign {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SystemRoleCanNotChange, "角色标识不能修改"),
		}, nil
	}

	var p big.Int
	for _, ps := range req.Policy {
		np, ok := (model.PermissionsSign())[ps]
		if ok && np.Status != db.PolicyStatusBanned {
			p = permission.AddPermission(p, np.Permission)
		}
	}

	var oldP big.Int
	oldP.SetString(role.Permissions, 16)

	if role.NotChangePermissions && oldP.Cmp(&p) != 0 {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SystemRoleCanNotChange, "角色权限不能修改"),
		}, nil
	}

	oldBelong := action.GetWebsite(role.Belong.Int64)

	if role.NotChangePermissions && oldBelong.ID != -1 && req.Belong != role.Belong.Int64 {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SystemRoleCanNotChange, "角色归属不能修改"),
		}, nil
	}

	if !db.IsRoleStatus(req.Status) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadRoleStatus, "错误的角色状态"),
		}, nil
	}

	if role.NotDelete && req.Status == db.RoleStatusBanned {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.SystemRoleCanNotChange, "角色不可禁用"),
		}, nil
	}

	web := action.GetWebsite(req.Belong)
	if web.ID == warp.UnknownWebsite {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadRoleBelong, "错误的归属站点"),
		}, nil
	}

	role.Name = req.Name
	role.Describe = req.Describe
	role.Permissions = p.Text(16)
	role.Sign = req.Sign
	role.Belong = sql.NullInt64{
		Valid: req.Belong != 0,
		Int64: req.Belong,
	}
	role.Status = req.Status

	err = roleModel.Update(l.ctx, role)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.RoleUpdateHandler(true)

	go func() {
		time.Sleep(10 * time.Second)
		urlpath.UpdateRoleByDB(role, nil)
	}()

	audit.NewAdminAudit(user.Id, "管理员更新角色（%s）成功", req.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
