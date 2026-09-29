package admin_announcement

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type AddAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAnnouncementLogic {
	return &AddAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAnnouncementLogic) AddAnnouncement(req *types.AdmnCreateAnnouncementReq) (resp *types.RespEmpty, err error) {
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
	sortNum, err := announcementModel.GetNewSortNumber(l.ctx)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	var startAt time.Time
	if req.StartAt == 0 {
		startAt = time.Now()
	} else {
		startAt = time.Unix(req.StartAt, 0)
	}

	var stopAt time.Time
	if req.StopAt == 0 {
		stopAt = startAt.Add(time.Hour * 24 * 365 * 100) // 一百年后
	} else {
		stopAt = time.Unix(req.StopAt, 0)
	}

	_, err = announcementModel.InsertCh(l.ctx, &db.Announcement{
		Sort:    sortNum,
		Title:   req.Title,
		Content: req.Content,
		StartAt: startAt,
		StopAt:  stopAt,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员添加公告（%s）成功", req.Title)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
