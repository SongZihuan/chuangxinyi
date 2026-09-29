package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/wxrobot"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type SendWXRobotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendWXRobotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendWXRobotLogic {
	return &SendWXRobotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendWXRobotLogic) SendWXRobot(req *types.AdminSendWXRobotReq) (resp *types.AdminSendMsgResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	srcUser, err := GetUser(l.ctx, req.ID, req.UID, true)
	if errors.Is(err, UserNotFound) {
		return &types.AdminSendMsgResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	wxrobotModel := db.NewWxrobotModel(mysql.MySQLConn)
	w, err := wxrobotModel.FindByUserID(l.ctx, srcUser.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.AdminSendMsgResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.AdminSendMsgData{
				Success: false,
				Have:    false,
			},
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if !w.Webhook.Valid {
		return &types.AdminSendMsgResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.AdminSendMsgData{
				Success: false,
				Have:    false,
			},
		}, nil
	}

	if web.ID == warp.UserCenterWebsite {
		err = wxrobot.Send(l.ctx, w.Webhook.String, req.Content, false, 0, config.BackendConfig.WXRobot.Sender)
	} else {
		bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
		allow, err := bannedModel.CheckAllow(l.ctx, user.Id, web.ID, db.AllowMsg)
		if err != nil || !allow {
			return &types.AdminSendMsgResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "用户关闭了通信授权许可"),
			}, nil
		}
		err = wxrobot.Send(l.ctx, w.Webhook.String, req.Content, false, web.ID, web.Name)
	}
	if err != nil {
		return &types.AdminSendMsgResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "发送信息失败"),
			Data: types.AdminSendMsgData{
				Success: false,
				Have:    true,
			},
		}, nil
	}

	audit.NewAdminAudit(user.Id, "管理员发送邮件给用户（%s）", srcUser.Uid)

	return &types.AdminSendMsgResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminSendMsgData{
			Success: true,
			Have:    true,
		},
	}, nil
}
