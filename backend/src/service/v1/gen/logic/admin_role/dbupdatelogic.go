package admin_role

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/cron"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/urlpath"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type DBUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDBUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DBUpdateLogic {
	return &DBUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DBUpdateLogic) DBUpdate() (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	cron.RoleUpdateHandler(true)

	go func() {
		time.Sleep(10 * time.Second)
		for _, r := range model.Roles() {
			urlpath.UpdateRole(r, nil)
		}
	}()

	audit.NewAdminAudit(user.Id, "管理员触发数据库角色缓存更新")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
