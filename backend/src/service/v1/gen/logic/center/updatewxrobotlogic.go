package center

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type UpdateWXRobotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWXRobotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWXRobotLogic {
	return &UpdateWXRobotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWXRobotLogic) UpdateWXRobot(req *types.UserUpdateWXRobotReq) (resp *types.RespEmpty, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	wxrobotModel := db.NewWxrobotModel(mysql.MySQLConn)
	if req.IsDelete || len(req.Webhook) == 0 {
		_, err = wxrobotModel.InsertWithDelete(l.ctx, &db.Wxrobot{
			UserId: user.Id,
			Webhook: sql.NullString{
				Valid: false,
			},
		})
	} else {
		_, err = wxrobotModel.InsertWithDelete(l.ctx, &db.Wxrobot{
			UserId: user.Id,
			Webhook: sql.NullString{
				Valid:  true,
				String: req.Webhook,
			},
		})
	}
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	sender.MessageSendChange(user.Id, "企业微信机器人")
	sender.WxrobotSendChange(user.Id, "企业微信机器人")
	audit.NewUserAudit(user.Id, "用户更新企业微信机器人webhook")

	return &types.RespEmpty{
		Resp: respmsg.GetRespSuccess(l.ctx),
	}, nil
}
