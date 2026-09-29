package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/sms"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type SendSMSLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendSMSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSMSLogic {
	return &SendSMSLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendSMSLogic) SendSMS(req *types.AdminSendSMSReq) (resp *types.AdminSendMsgResp, err error) {
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

	phoneModel := db.NewPhoneModel(mysql.MySQLConn)
	p, err := phoneModel.FindByUserID(l.ctx, srcUser.Id)
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
	}

	data := make(map[string]string, len(req.Data))
	for _, v := range req.Data {
		data[v.Label] = v.Value
	}

	if web.ID == warp.UserCenterWebsite {
		err = sms.Send(l.ctx, data, req.Sig, req.TemplateID, p.Phone, warp.UserCenterWebsite)
	} else {
		bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
		allow, err := bannedModel.CheckAllow(l.ctx, user.Id, web.ID, db.AllowMsg)
		if err != nil || !allow {
			return &types.AdminSendMsgResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "用户关闭了通信授权许可"),
			}, nil
		}
		err = sms.Send(l.ctx, data, req.Sig, req.TemplateID, p.Phone, web.ID)
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

	audit.NewAdminAudit(user.Id, "管理员发送短信给用户（%s）", srcUser.Uid)

	return &types.AdminSendMsgResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminSendMsgData{
			Success: true,
			Have:    true,
		},
	}, nil
}
