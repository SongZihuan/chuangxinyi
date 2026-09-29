package center

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/jwt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"

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

func (l *UpdateEmailLogic) UpdateEmail(req *types.UserUpdateEmailReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	emailModel := db.NewEmailModel(mysql.MySQLConn)

	e, err := emailModel.FindByUserID(context.Background(), user.Id)
	if errors.Is(err, db.ErrNotFound) || err != nil {
		e = &db.Email{
			Email: sql.NullString{
				String: "",
			},
		}
	}

	var email string
	if req.IsDelete || len(req.EmailToken) == 0 {
		email = ""
		err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			emailModel := db.NewEmailModelWithSession(session)
			_, err := emailModel.InsertWithDelete(l.ctx, &db.Email{
				UserId: user.Id,
				Email: sql.NullString{
					Valid: false,
				},
			})
			return err
		})
	} else {
		var emailData jwt.EmailTokenData
		emailData, err = jwt.ParserEmailToken(req.EmailToken)
		if err != nil {
			return nil, respmsg.JWTError.WarpQuick(err)
		} else if emailData.WebID != warp.UserCenterWebsite {
			return nil, respmsg.JWTError.New("bad website")
		}

		email = emailData.Email

		err = mysql.MySQLConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			emailModel := db.NewEmailModelWithSession(session)
			_, err := emailModel.InsertWithDelete(l.ctx, &db.Email{
				UserId: user.Id,
				Email: sql.NullString{
					Valid:  true,
					String: email,
				},
			})
			return err
		})
	}
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.PhoneSendChange(user.Id, "邮箱")
	sender.EmailSendEmailChange(e.Email.String, email) // 不需要判断e.Email是否为空，内部会判断
	sender.EmailSendBind(email)                        // 不需要判断email是否为空，内部会判断
	sender.MessageSendChange(user.Id, "邮箱")
	sender.WxrobotSendChange(user.Id, "邮箱")
	sender.FuwuhaoSendChange(user.Id, "邮箱")
	audit.NewUserAudit(user.Id, "用户邮箱更新成功")

	fmt.Println("TAG F")
	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
