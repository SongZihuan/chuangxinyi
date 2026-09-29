package center

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
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/yundun"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
	"time"

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

func (l *UpdateUserNameLogic) UpdateUserName(req *types.UserUpdateUserNameReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if req.IsDelete || len(req.UserName) == 0 {
		err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			usernameModel := db.NewUsernameModelWithSession(session)
			_, err := usernameModel.InsertSafe(l.ctx, &db.Username{
				UserId: user.Id,
				Username: sql.NullString{
					Valid: false,
				},
			})
			return err
		})
	} else {
		res, checkErr := yundun.CheckName(req.UserName)
		if checkErr != nil {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByError(l.ctx, respmsg.BadUserName, checkErr, "用户名违规检测失败"),
			}, nil
		}

		if !res {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadUserName, "用户名违规"),
			}, nil
		}

		username := base64.StdEncoding.EncodeToString([]byte(req.UserName))
		lock := fmt.Sprintf("username:%s", username)
		if !redis.AcquireLock(l.ctx, lock, time.Minute*2) {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNameHasBeenRegister, "上锁失败，用户名重复"),
			}, nil
		}
		defer redis.ReleaseLock(lock)

		usernameModel := db.NewUsernameModel(mysql.MySQLConn)
		_, mysqlErr := usernameModel.FindByUsername(l.ctx, username)
		if mysqlErr != nil && !errors.Is(mysqlErr, db.ErrNotFound) {
			return nil, respmsg.MySQLSystemError.WarpQuick(mysqlErr)
		} else if mysqlErr == nil {
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNameRepeat, "数据库存在用户名，用户名重复"),
			}, nil
		}

		err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			usernameModel := db.NewUsernameModelWithSession(session)
			_, err := usernameModel.InsertSafe(l.ctx, &db.Username{
				UserId: user.Id,
				Username: sql.NullString{
					Valid:  true,
					String: username,
				},
			})
			return err
		})
	}
	if errors.Is(err, db.UserNameRepeat) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNameRepeat, "用户名重复"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.PhoneSendChange(user.Id, "用户名")
	sender.EmailSendChange(user.Id, "用户名")
	sender.MessageSendChange(user.Id, "用户名")
	sender.WxrobotSendChange(user.Id, "用户名")
	sender.FuwuhaoSendChange(user.Id, "用户名")
	audit.NewUserAudit(user.Id, "用户更新用户名成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
