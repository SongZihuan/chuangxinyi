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

type DeleteWechatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWechatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWechatLogic {
	return &DeleteWechatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWechatLogic) DeleteWechat(req *types.AdminDeleteWechatReq) (resp *types.RespEmpty, err error) {
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

	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	_, err = wechatModel.InsertWithDelete(l.ctx, &db.Wechat{
		UserId: srcUser.Id,
		OpenId: sql.NullString{
			Valid: false,
		},
		IsDelete: db.IsBanned(srcUser),
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员解绑用户（%s）微信成功", srcUser.Uid)

	if !db.IsBanned(srcUser) {
		sender.PhoneSendChange(user.Id, "微信")
		sender.MessageSendChange(user.Id, "微信")
		sender.WxrobotSendChange(user.Id, "微信")
		audit.NewUserAudit(user.Id, "用户更新微信绑定成功")
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
