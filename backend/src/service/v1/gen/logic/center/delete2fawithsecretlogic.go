package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type Delete2FAWithSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelete2FAWithSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Delete2FAWithSecretLogic {
	return &Delete2FAWithSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Delete2FAWithSecretLogic) Delete2FAWithSecret(req *types.Delete2FAWithSecretReq) (resp *types.RespEmpty, err error) {
	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)

	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	secondfa, err := secondfaModel.FindByUserID(l.ctx, user.Id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			audit.NewUserAudit(user.Id, "用户试图解绑2FA，但实际上并未绑定")
			return &types.RespEmpty{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserHasNotBeenBind2FA, "用户未绑定2FA"),
			}, nil
		}
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	if !secondfa.Secret.Valid {
		audit.NewUserAudit(user.Id, "用户试图解绑2FA，但实际上并未绑定")
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserHasNotBeenBind2FA, "用户未绑定2FA"),
		}, nil
	}

	if !utils.CheckTOTP(secondfa.Secret.String, req.Code) {
		audit.NewUserAudit(user.Id, "用户试图解绑2FA，但验证不通过")
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Bad2FACode, "2FA检验失败"),
		}, nil
	}

	_, err = secondfaModel.InsertWithDelete(l.ctx, &db.Secondfa{
		UserId: user.Id,
		Secret: sql.NullString{
			Valid: false,
		},
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.MessageSendChange(user.Id, "2FA-双因素验证")
	sender.WxrobotSendChange(user.Id, "2FA-双因素验证")
	audit.NewUserAudit(user.Id, "用户解绑2FA成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
