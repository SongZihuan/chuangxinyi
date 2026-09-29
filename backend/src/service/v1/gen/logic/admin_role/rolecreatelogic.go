package admin_role

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"math/big"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type RoleCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleCreateLogic {
	return &RoleCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleCreateLogic) RoleCreate(req *types.CreateRoleReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	roleModel := db.NewRoleModel(mysql.MySQLConn)

	count, err := roleModel.GetCount(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if count > config.BackendConfig.MySQL.SystemResourceLimit {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "超出限额"),
		}, nil
	}

	var p big.Int
	for _, ps := range req.Policy {
		np, ok := (model.PermissionsSign())[ps]
		if ok && np.Status != db.PolicyStatusBanned {
			p = permission.AddPermission(p, np.Permission)
		}
	}

	if !db.IsRoleStatus(req.Status) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadRoleStatus, "错误的角色状态"),
		}, nil
	}

	web := action.GetWebsite(req.Belong)
	if web.ID == warp.UnknownWebsite {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadRoleBelong, "错误的角色归属"),
		}, nil
	}

	_, err = roleModel.Insert(l.ctx, &db.Role{
		Name:                 req.Name,
		Describe:             req.Describe,
		Sign:                 req.Sign,
		Status:               req.Status,
		NotDelete:            req.NotDelete,
		NotChangePermissions: req.NotChangePermissions,
		NotChangeSign:        req.NotChangeSign,
		Belong: sql.NullInt64{
			Valid: req.Belong != 0,
			Int64: req.Belong,
		},
		Permissions: p.Text(16),
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.RoleUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员创建角色（%s）成功", req.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
