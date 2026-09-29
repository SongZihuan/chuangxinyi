package admin_user

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/password"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

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

func (l *UpdatePasswordLogic) UpdatePassword(req *types.AdminUpdatePasswordReq) (resp *types.RespEmpty, err error) {
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

	passwordModel := db.NewPasswordModel(mysql.MySQLConn)
	if req.IsDelete || len(req.PasswordHash) == 0 {
		_, err = passwordModel.InsertWithDelete(l.ctx, &db.Password{
			UserId: srcUser.Id,
			PasswordHash: sql.NullString{
				Valid: false,
			},
		})
	} else {
		_, err = passwordModel.InsertWithDelete(l.ctx, &db.Password{
			UserId: srcUser.Id,
			PasswordHash: sql.NullString{
				Valid:  true,
				String: password.GetPasswordSecondHash(req.PasswordHash, srcUser.Uid),
			},
		})
	}
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新用户（%s）密码成功", srcUser.Uid)

	if db.IsBanned(srcUser) {
		audit.NewUserAudit(srcUser.Id, "用户更新密码成功")
		sender.PhoneSendChange(srcUser.Id, "密码")
		sender.EmailSendChange(srcUser.Id, "密码")
		sender.MessageSendChange(srcUser.Id, "密码")
		sender.WxrobotSendChange(srcUser.Id, "密码")
		sender.FuwuhaoSendChange(srcUser.Id, "密码")
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
