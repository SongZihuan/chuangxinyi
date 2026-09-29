package admin_announcement

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAnnouncementLogic {
	return &UpdateAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAnnouncementLogic) UpdateAnnouncement(req *types.AdmnUpdateAnnouncementReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	key := fmt.Sprintf("sort:announcement")
	if !redis.AcquireLockMore(l.ctx, key, time.Minute*2) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotGetLock, "无法上锁，因为涉及排序操作"),
		}, nil
	}
	defer redis.ReleaseLock(key)

	announcementModel := db.NewAnnouncementModel(mysql.MySQLConn)

	announcement, err := announcementModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.AnnouncementNotFound, "公告未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	var startAt time.Time
	if req.StartAt == 0 {
		startAt = time.Now()
	} else {
		startAt = time.Unix(req.StartAt, 0)
	}

	var stopAt time.Time
	if req.StartAt == 0 {
		stopAt = startAt.Add(time.Hour * 24 * 365 * 100) // 一百年后
	} else {
		stopAt = time.Unix(req.StopAt, 0)
	}

	announcement.Title = req.Title
	announcement.Content = req.Content
	announcement.StartAt = startAt
	announcement.StopAt = stopAt

	err = announcementModel.UpdateCh(l.ctx, announcement)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新公告（%s）成功", req.Title)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
