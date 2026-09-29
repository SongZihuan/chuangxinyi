package admin_menu

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
	"math/big"
	"time"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type MoveMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMoveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MoveMenuLogic {
	return &MoveMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MoveMenuLogic) MoveMenu(req *types.MoveReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	menuModel := db.NewMenuModel(mysql.MySQLConn)

	menu, err := menuModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.MenuNotFount, "菜单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	key := fmt.Sprintf("sort:menu:%d", menu.FatherId.Int64)
	if !redis.AcquireLockMore(l.ctx, key, time.Minute*2) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotGetLock, "不能上锁，因为涉及排序操作"),
		}, nil
	}
	defer redis.ReleaseLock(key)

	near, err := menuModel.FindNear(l.ctx, menu.FatherId.Int64, menu.Sort, req.IsUp)
	if errors.Is(err, db.ErrNotFound) {
		if req.IsUp {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotUp, "菜单已经在顶部"),
			}, nil
		} else {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotDown, "菜单已经在底部"),
			}, nil
		}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	tmp := menu.Sort
	menu.Sort = near.Sort
	near.Sort = tmp

	err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) (err error) {
		menuModel := db.NewMenuModelWithSession(session)

		err = menuModel.Update(ctx, menu)
		if err != nil {
			return err
		}

		err = menuModel.Update(ctx, near)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		var p1, p2 big.Int
		_, ok = p1.SetString(menu.Policy, 16)

		if ok {
			for _, r := range model.Roles() {
				if permission.HasOnePermission(r.PolicyPermissions, p1) {
					urlpath.UpdateRole(r, nil)
				}
			}
		}

		_, ok = p2.SetString(near.Policy, 16)
		if ok {
			for _, r := range model.Roles() {
				if permission.HasOnePermission(r.PolicyPermissions, p2) {
					urlpath.UpdateRole(r, nil)
				}
			}
		}
	}()

	cron.MenuUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员移动菜单（%s）成功", menu.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
