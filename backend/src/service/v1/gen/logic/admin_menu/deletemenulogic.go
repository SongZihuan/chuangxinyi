package admin_menu

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	errors "github.com/wuntsong-org/wterrors"
	"math/big"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteReq) (resp *types.RespEmpty, err error) {
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

	if !func() bool {
		for _, m := range model.Menus() {
			if m.FatherID == menu.Id {
				return false
			}
		}
		return true
	}() {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.MenuHasSon, "菜单有儿子，不能直接删除"),
		}, nil
	}

	menu.DeleteAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	err = menuModel.Update(l.ctx, menu)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		var p big.Int
		_, ok = p.SetString(menu.Policy, 16)

		if ok {
			for _, r := range model.Roles() {
				if permission.HasOnePermission(r.PolicyPermissions, p) {
					urlpath.UpdateRole(r, nil)
				}
			}
		}
	}()

	cron.MenuUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员删除菜单（%s）成功", menu.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
