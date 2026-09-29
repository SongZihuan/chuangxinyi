package center

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

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

func (l *UpdateLoginControllerLogic) UpdateLoginController(req *types.UpdateLoginControllerReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	ctrlModel := db.NewLoginControllerModel(mysql.MySQLConn)

	_, err = ctrlModel.InsertWithDelete(l.ctx, &db.LoginController{
		UserId:        user.Id,
		AllowPhone:    req.AllowPhone,
		AllowEmail:    req.AllowEmail,
		AllowWechat:   req.AllowWeChat,
		AllowPassword: req.AllowPassword,
		Allow2Fa:      req.AllowSecondFA,
	})
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.MessageSendChange(user.Id, "登录限制")
	sender.WxrobotSendChange(user.Id, "登录限制")
	sender.FuwuhaoSendChange(user.Id, "登录限制")
	audit.NewUserAudit(user.Id, "用户登录限制更新成功")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
