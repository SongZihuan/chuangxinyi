package admin_announcement

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
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

type MoveAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMoveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MoveAnnouncementLogic {
	return &MoveAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MoveAnnouncementLogic) MoveAnnouncement(req *types.MoveReq) (resp *types.RespEmpty, err error) {
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

	near, err := announcementModel.FindNear(l.ctx, announcement.Sort, req.IsUp)
	if errors.Is(err, db.ErrNotFound) {
		if req.IsUp {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotUp, "公告在顶部"),
			}, nil
		} else {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.CanNotDown, "公告在底部"),
			}, nil
		}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	tmp := announcement.Sort
	announcement.Sort = near.Sort
	near.Sort = tmp

	err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) (err error) {
		announcementModel := db.NewAnnouncementModelWithSession(session)

		err = announcementModel.UpdateCh(ctx, announcement)
		if err != nil {
			return err
		}

		err = announcementModel.UpdateCh(ctx, near)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员移动公告成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
