package admin_user

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmailLogic {
	return &UpdateEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEmailLogic) UpdateEmail(req *types.AdminUpdateEmailReq) (resp *types.RespEmpty, err error) {
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

	emailModel := db.NewEmailModel(mysql.MySQLConn)

	e, err := emailModel.FindByUserID(l.ctx, srcUser.Id)
	if errors.Is(err, db.ErrNotFound) {
		e = &db.Email{
			Email: sql.NullString{
				String: "",
			},
		}
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if req.IsDelete || len(req.Email) == 0 {
		req.Email = ""
		_, err = emailModel.InsertWithDelete(l.ctx, &db.Email{
			UserId: srcUser.Id,
			Email: sql.NullString{
				Valid: false,
			},
			IsDelete: db.IsBanned(srcUser),
		})
	} else {
		_, err = emailModel.InsertWithDelete(l.ctx, &db.Email{
			UserId: srcUser.Id,
			Email: sql.NullString{
				Valid:  true,
				String: req.Email,
			},
			IsDelete: db.IsBanned(srcUser),
		})
	}
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新用户（%s）邮箱成功", srcUser.Uid)

	if !db.IsBanned(srcUser) {
		sender.PhoneSendChange(srcUser.Id, "邮箱")
		sender.EmailSendEmailChange(e.Email.String, req.Email) // 不需要判断e.Email是否为空，内部会判断
		sender.EmailSendBind(req.Email)                        // 不需要判断email是否为空，内部会判断
		sender.MessageSendChange(srcUser.Id, "邮箱")
		sender.WxrobotSendChange(srcUser.Id, "邮箱")
		sender.FuwuhaoSendChange(srcUser.Id, "邮箱")
		audit.NewUserAudit(srcUser.Id, "用户邮箱更新成功")
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
