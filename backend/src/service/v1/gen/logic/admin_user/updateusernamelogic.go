package admin_user

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateUserNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserNameLogic {
	return &UpdateUserNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserNameLogic) UpdateUserName(req *types.AdminUpdateUserNameReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	usernameModel := db.NewUsernameModel(mysql.MySQLConn)
	if req.IsDelete || len(req.UserName) == 0 {
		_, err = usernameModel.InsertSafe(l.ctx, &db.Username{
			UserId: srcUser.Id,
			Username: sql.NullString{
				Valid: false,
			},
			IsDelete: db.IsBanned(srcUser),
		})
	} else {
		username := base64.StdEncoding.EncodeToString([]byte(req.UserName))

		if !db.IsBanned(srcUser) {
			lock := fmt.Sprintf("username:%s", username)
			if !redis.AcquireLock(l.ctx, lock, time.Minute*2) {
				return &types.RespEmpty{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNameRepeat, "无法上锁，用户名重复"),
				}, nil
			}
			defer redis.ReleaseLock(lock)

			_, err = usernameModel.FindByUsername(l.ctx, username)
			if err != nil && !errors.Is(err, db.ErrNotFound) {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			} else if err == nil {
				return &types.RespEmpty{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNameRepeat, "数据库已存在用户名"),
				}, nil
			}
		}

		_, err = usernameModel.InsertSafe(l.ctx, &db.Username{
			UserId: srcUser.Id,
			Username: sql.NullString{
				Valid:  true,
				String: username,
			},
		})
	}
	if errors.Is(err, db.UserNameRepeat) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNameRepeat, "用户名重复"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新用户（%s）用户名成功", srcUser.Uid)

	if !db.IsBanned(srcUser) {
		sender.PhoneSendChange(srcUser.Id, "用户名")
		sender.EmailSendChange(srcUser.Id, "用户名")
		sender.MessageSendChange(srcUser.Id, "用户名")
		sender.WxrobotSendChange(srcUser.Id, "用户名")
		sender.FuwuhaoSendChange(srcUser.Id, "用户名")
		audit.NewUserAudit(srcUser.Id, "用户更新用户名成功")
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
