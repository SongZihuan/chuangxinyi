package admin_website_permission

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPermissionLogic {
	return &AddPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPermissionLogic) AddPermission(req *types.CreateWebsitePermissionReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	key := fmt.Sprintf("sort:website:permission")
	if !redis.AcquireLockMore(l.ctx, key, time.Minute*2) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotGetLock, "上锁失败，因为涉及排序"),
		}, nil
	}
	defer redis.ReleaseLock(key)

	if !db.IsWebsitePolicyStatus(req.Status) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPolicyStatus, "错误的站点权限状态"),
		}, nil
	}

	policyModel := db.NewWebsitePolicyModel(mysql.MySQLConn)

	count, err := policyModel.GetCount(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if count > config.BackendConfig.MySQL.SystemResourceLimit {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.TooMany, "超出限额"),
		}, nil
	}

	sortNum, err := policyModel.GetNewSortNumber(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	_, err = policyModel.Insert(l.ctx, &db.WebsitePolicy{
		Describe: req.Describe,
		Sort:     sortNum,
		Sign:     req.Sign,
		Name:     req.Name,
		Status:   req.Status,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	cron.WebsitePermissionUpdateHandler(true)
	audit.NewAdminAudit(user.Id, "管理员新增站点权限（%s）成功", req.Name)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
