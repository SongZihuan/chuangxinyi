package admin_permission

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type MovePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMovePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MovePermissionLogic {
	return &MovePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MovePermissionLogic) MovePermission(req *types.MoveReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	key := fmt.Sprintf("sort:permission")
	if !redis.AcquireLockMore(l.ctx, key, time.Minute*2) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotGetLock, "上锁失败，因为涉及排序"),
		}, nil
	}
	defer redis.ReleaseLock(key)

	permissionModel := db.NewPolicyModel(mysql.MySQLConn)

	permission, err := permissionModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PermissionNotFound, "权限未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	near, err := permissionModel.FindNear(l.ctx, permission.Sort, req.IsUp)
	if errors.Is(err, db.ErrNotFound) {
		if req.IsUp {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotUp, "权限已经在顶部"),
			}, nil
		} else {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotDown, "权限已经在底部"),
			}, nil
		}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	tmp := permission.Sort
	permission.Sort = near.Sort
	near.Sort = tmp

	err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) (err error) {
		permissionModel := db.NewPolicyModelWithSession(session)

		err = permissionModel.Update(ctx, permission)
		if err != nil {
			return err
		}

		err = permissionModel.Update(ctx, near)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.PermissionUpdateHandler(true)
	// 不需要更新role
	audit.NewAdminAudit(user.Id, "管理员移动权限（%s）成功", permission.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
