package admin_menu

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	errors "github.com/wuntsong-org/wterrors"
	"math/big"
	"time"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.MenuUpdateReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !db.IsMenuStatus(req.Status) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadMenuStatus, "入参status错误"),
		}, nil
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

	var sort int64
	if req.FatherID != menu.FatherId.Int64 {
		fatherID := req.FatherID
		fatherCount := config.BackendConfig.Admin.MenuDepth - 1

		for fatherCount > 0 && fatherID != 0 {
			fatherCount--
			if fatherID == req.ID { // 循环引用
				break
			}

			m, ok := model.Menus()[fatherID]
			if !ok {
				fatherID = 0
				break
			}

			fatherID = m.FatherID
		}

		if fatherID != 0 {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.MenuNotFount, "循环引用"),
			}, nil
		}

		sort, err = menuModel.GetNewSortNumber(l.ctx, req.FatherID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		sort = menu.Sort
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

	var oldP big.Int
	_, ok = oldP.SetString(menu.Policy, 16)
	if !ok {
		oldP = *big.NewInt(0)
	}

	menu.Describe = req.Describe
	menu.FatherId = sql.NullInt64{
		Valid: req.FatherID != 0,
		Int64: req.FatherID,
	}
	menu.Sort = sort
	menu.Name = req.Name
	menu.Path = req.Path
	menu.Title = req.Title
	menu.Icon = req.Icon
	menu.Redirect = sql.NullString{
		Valid:  len(req.Redirect) != 0,
		String: req.Redirect,
	}
	menu.Superior = req.Superior
	menu.Category = req.Category
	menu.Component = req.Component
	menu.ComponentAlias = req.ComponentAlias
	menu.MetaLink = sql.NullString{
		Valid:  len(req.MetaLink) != 0,
		String: req.MetaLink,
	}
	menu.Type = req.Type
	menu.IsLink = req.IsLink
	menu.IsHide = req.IsHide
	menu.IsKeepalive = req.IsKeepalive
	menu.IsAffix = req.IsAffix
	menu.IsIframe = req.IsIframe
	menu.BtnPower = req.BtnPower
	menu.Status = req.Status
	menu.Policy = p.Text(16)
	menu.IsOrPolicy = req.IsOr
	menu.SubPolicy = sp

	err = menuModel.Update(l.ctx, menu)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	go func() {
		time.Sleep(time.Second * 10)

		for _, r := range model.Roles() {
			if permission.HasOnePermission(r.PolicyPermissions, p) || permission.HasAllPermission(r.PolicyPermissions, oldP) {
				urlpath.UpdateRole(r, nil)
			}
		}
	}()

	cron.MenuUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员更新菜单（%s）单成功", menu.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
