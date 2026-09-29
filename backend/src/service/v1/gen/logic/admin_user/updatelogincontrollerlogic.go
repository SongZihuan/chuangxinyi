package admin_user

import (
	"context"
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

type UpdateLoginControllerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLoginControllerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoginControllerLogic {
	return &UpdateLoginControllerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLoginControllerLogic) UpdateLoginController(req *types.AdminUpdateLoginControllerReq) (resp *types.RespEmpty, err error) {
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

	ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

	_, err = ctrlModel.InsertWithDelete(l.ctx, &db.LoginController{
		UserId:        srcUser.Id,
		AllowPhone:    req.AllowPhone,
		AllowEmail:    req.AllowEmail,
		AllowWechat:   req.AllowWeChat,
		AllowPassword: req.AllowPassword,
		Allow2Fa:      req.AllowSecondFA,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	audit.NewAdminAudit(user.Id, "管理员更新用户（%s）登录限制成功", srcUser.Uid)

	if db.IsBanned(srcUser) {
		sender.MessageSendChange(srcUser.Id, "登录限制")
		sender.WxrobotSendChange(srcUser.Id, "登录限制")
		sender.FuwuhaoSendChange(srcUser.Id, "登录限制")
		audit.NewUserAudit(srcUser.Id, "用户登录限制更新成功")
	}

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
