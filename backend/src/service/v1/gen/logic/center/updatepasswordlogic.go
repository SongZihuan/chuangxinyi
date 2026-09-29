package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/password"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(req *types.UserUpdatePasswordReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	if !utils.IsSha256(req.NewPasswordHash) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.BadPassword, "密码哈希值错误"),
		}, nil
	}

	if req.IsDelete || len(req.NewPasswordHash) == 0 {
		err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			passwordModel := db.NewPasswordModelWithSession(session)
			_, err := passwordModel.InsertWithDelete(l.ctx, &db.Password{
				UserId: user.Id,
				PasswordHash: sql.NullString{
					Valid: false,
				},
			})
			return err
		})
	} else {
		err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			passwordModel := db.NewPasswordModelWithSession(session)
			_, err := passwordModel.InsertWithDelete(l.ctx, &db.Password{
				UserId: user.Id,
				PasswordHash: sql.NullString{
					Valid:  true,
					String: password.GetPasswordSecondHash(req.NewPasswordHash, user.Uid),
				},
			})
			return err
		})
	}
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "用户更新密码成功")
	sender.PhoneSendChange(user.Id, "密码")
	sender.EmailSendChange(user.Id, "密码")
	sender.MessageSendChange(user.Id, "密码")
	sender.WxrobotSendChange(user.Id, "密码")
	sender.FuwuhaoSendChange(user.Id, "密码")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
