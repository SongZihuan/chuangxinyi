package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/fuwuhao"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type SendFuwuhaoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendFuwuhaoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFuwuhaoLogic {
	return &SendFuwuhaoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendFuwuhaoLogic) SendFuwuhao(req *types.AdminSendFuwuhaoReq) (resp *types.AdminSendMsgResp, err error) {
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

	wechatModel := db.NewWechatModel(mysql.MySQLConn)
	w, err := wechatModel.FindByUserID(l.ctx, srcUser.Id)
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
	} else if !w.Fuwuhao.Valid {
		return &types.AdminSendMsgResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.AdminSendMsgData{
				Success: false,
				Have:    false,
			},
		}, nil
	}

	val := make(map[string]string, len(req.Val))
	for _, v := range req.Val {
		val[v.Label] = v.Value
	}

	if web.ID == warp.UserCenterWebsite {
		err = fuwuhao.SendVal(l.ctx, req.TemplateID, req.Url, w.Fuwuhao.String, val, warp.UserCenterWebsite)
	} else {
		bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
		allow, err := bannedModel.CheckAllow(l.ctx, user.Id, web.ID, db.AllowMsg)
		if err != nil || !allow {
			return &types.AdminSendMsgResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "用户关闭了通信授权许可"),
			}, nil
		}
		err = fuwuhao.SendVal(l.ctx, req.TemplateID, req.Url, w.Fuwuhao.String, val, web.ID)
	}
	if err != nil {
		logger.Logger.Error("send fuwuhao msg error: %s", err.Error())
		return &types.AdminSendMsgResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "发送信息失败"),
			Data: types.AdminSendMsgData{
				Success: false,
				Have:    true,
			},
		}, nil
	}

	audit.NewAdminAudit(user.Id, "管理员发送服务号消息给用户（%s）", srcUser.Uid)

	return &types.AdminSendMsgResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminSendMsgData{
			Success: true,
			Have:    true,
		},
	}, nil
}
