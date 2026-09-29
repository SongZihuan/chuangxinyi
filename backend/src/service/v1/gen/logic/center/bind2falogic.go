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

type Bind2FALogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBind2FALogic(ctx context.Context, svcCtx *svc.ServiceContext) *Bind2FALogic {
	return &Bind2FALogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Bind2FALogic) Bind2FA(req *types.Bind2FAReq) (resp *types.RespEmpty, err error) {
	if !utils.IsTotpSecret(req.Secret) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Bad2FASecret, "2FA的密钥错误"),
		}, nil
	}

	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	secondfaModel := db.NewSecondfaModel(mysql.MySQLConn)
	secondfaOld, err := secondfaModel.FindByUserID(l.ctx, user.Id)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if !errors.Is(err, db.ErrNotFound) && secondfaOld.Secret.Valid {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserHasBeenBind2FA, "用户已经绑定2FA"),
		}, nil
	}

	if !utils.CheckTOTP(req.Secret, req.Code) {
		return &types.RespEmpty{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.Bad2FACode, "用户2FA检查失败"),
		}, nil
	}

	_, err = secondfaModel.InsertWithDelete(l.ctx, &db.Secondfa{
		UserId: user.Id,
		Secret: sql.NullString{
			Valid:  true,
			String: req.Secret,
		},
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.MessageSendChange(user.Id, "2FA-双因素验证")
	sender.WxrobotSendChange(user.Id, "2FA-双因素验证")
	audit.NewUserAudit(user.Id, "用户绑定2FA成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
