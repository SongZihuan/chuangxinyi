package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type ReadMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadMessageLogic {
	return &ReadMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadMessageLogic) ReadMessage(req *types.ReadMessageReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	messageModel := db.NewMessageModel(mysql.MySQLConn)
	message, err := messageModel.FindOneWithoutDelete(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if message.UserId != user.Id {
		return &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		}, nil
	} else if web.ID != warp.UserCenterWebsite && web.ID != message.SenderId {
		return &types.RespEmpty{
			Resp: respmsg.GetRespSuccess(l.ctx),
		}, nil
	}

	message.ReadAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	err = messageModel.UpdateReadCh(l.ctx, message)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户阅读站内信：%d", message.Id)

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
