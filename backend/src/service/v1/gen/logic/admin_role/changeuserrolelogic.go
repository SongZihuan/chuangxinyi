package admin_role

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	"github.com/wuntsong-org/go-zero-plus/core/logx"
	errors "github.com/wuntsong-org/wterrors"
)

type ChangeUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserRoleLogic {
	return &ChangeUserRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeUserRoleLogic) ChangeUserRole(req *types.ChangeRoleReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	var srcUser *db.User

	userModel := db.NewUserModel(mysql.MySQLConn)

	if len(req.UID) == 0 {
		srcUser, err = userModel.FindOneByIDWithoutDelete(l.ctx, req.ID)
	} else {
		srcUser, err = utils2.FindUser(l.ctx, req.UID, true)
	}
	if errors.Is(err, utils2.UserNotFound) || errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	role := action.GetRole(req.ID, srcUser.IsAdmin)

	if srcUser.IsAdmin {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.RootAdminCannotChangeRole, "根管理员不能改变角色"),
		}, nil
	}

	srcUser.RoleId = role.ID
	err = userModel.UpdateChWithoutStatus(l.ctx, srcUser)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	urlpath.ChangeRole(srcUser.Id, role.ID)
	audit.NewAdminAudit(user.Id, "管理员更新用户角色成功：%s", srcUser.Uid)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
