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

type Update2FALogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdate2FALogic(ctx context.Context, svcCtx *svc.ServiceContext) *Update2FALogic {
	return &Update2FALogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Update2FALogic) Update2FA(req *types.AdminDelete2FAReq) (resp *types.RespEmpty, err error) {
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

	secondFAModel := db.NewSecondfaModel(mysql.MySQLConn)
	_, err = secondFAModel.InsertWithDelete(l.ctx, &db.Secondfa{
		UserId: srcUser.Id,
		Secret: sql.NullString{
			Valid: false,
		},
	})

	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员解绑用户（%s）2FA成功", srcUser.Uid)

	if !db.IsBanned(srcUser) {
		sender.MessageSendChange(srcUser.Id, "2FA-双因素验证")
		sender.WxrobotSendChange(srcUser.Id, "2FA-双因素验证")
		audit.NewUserAudit(srcUser.Id, "用户解绑2FA成功")
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
