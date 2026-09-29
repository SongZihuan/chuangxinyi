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
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	errors "github.com/wuntsong-org/wterrors"
	"math/big"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMenuLogic {
	return &AddMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMenuLogic) AddMenu(req *types.CreateMenuReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !db.IsMenuStatus(req.Status) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadMenuStatus, "入参status错误"),
		}, nil
	}

	key := fmt.Sprintf("sort:menu:%d", req.FatherID)
	if !redis.AcquireLockMore(l.ctx, key, time.Minute*2) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotGetLock, "无法上锁，因为涉及排序"),
		}, nil
	}
	defer redis.ReleaseLock(key)

	fatherID := req.FatherID
	fatherCount := config.BackendConfig.Admin.MenuDepth - 1

	for fatherCount > 0 && fatherID != 0 {
		fatherCount--
		m, ok := (model.Menus())[fatherID]
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

	menuModel := db.NewMenuModel(mysql.MySQLConn)

	count, err := menuModel.GetCount(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if count > config.BackendConfig.MySQL.SystemResourceLimit {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "超出限额"),
		}, nil
	}

	sortNum, err := menuModel.GetNewSortNumber(l.ctx, req.FatherID)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if req.FatherID != 0 {
		_, err := menuModel.FindOneWithoutDelete(l.ctx, req.FatherID)
		if errors.Is(err, db.ErrNotFound) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.MenuNotFount, "菜单父亲丢失"),
			}, nil
		}
	}

	var p big.Int
	for _, ps := range req.Policy {
		np, ok := (model.PermissionsSign())[ps]
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

	_, err = menuModel.Insert(l.ctx, &db.Menu{
		Describe: req.Describe,
		Sort:     sortNum,
		FatherId: sql.NullInt64{
			Valid: req.FatherID != 0,
			Int64: req.FatherID,
		},
		Name:  req.Name,
		Path:  req.Path,
		Title: req.Title,
		Icon:  req.Icon,
		Redirect: sql.NullString{
			Valid:  len(req.Redirect) != 0,
			String: req.Redirect,
		},
		Superior:       req.Superior,
		Category:       req.Category,
		Component:      req.Component,
		ComponentAlias: req.ComponentAlias,
		MetaLink: sql.NullString{
			Valid:  len(req.MetaLink) != 0,
			String: req.MetaLink,
		},
		Type:        req.Type,
		IsLink:      req.IsLink,
		IsHide:      req.IsHide,
		IsKeepalive: req.IsKeepalive,
		IsAffix:     req.IsAffix,
		IsIframe:    req.IsIframe,
		BtnPower:    req.BtnPower,
		Status:      req.Status,
		IsOrPolicy:  req.IsOr,
		Policy:      p.Text(16),
		SubPolicy:   sp,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	go func() {
		time.Sleep(10 * time.Second) // 睡眠10s等更新同步
		for _, r := range model.Roles() {
			if permission.HasOnePermission(r.PolicyPermissions, p) {
				urlpath.UpdateRole(r, nil)
			}
		}
	}()

	cron.MenuUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员新增菜单（%s）成功", req.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
