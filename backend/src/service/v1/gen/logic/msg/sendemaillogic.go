package msg

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/email"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type SendEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailLogic {
	return &SendEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendEmailLogic) SendEmail(req *types.SendEmailReq, r *http.Request) (resp *types.SendMsgResp, err error) {
	emailModel := db.NewEmailModel(mysql.MySQLConn)

	user, err := utils2.FindUser(l.ctx, req.UserID, false)
	if errors.Is(err, utils2.UserNotFound) {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	e, err := emailModel.FindByUserID(l.ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.SendMsgData{
				Success: false,
				Have:    false,
			},
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if !e.Email.Valid {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.SendMsgData{
				Success: false,
				Have:    false,
			},
		}, nil
	}

	web, ok := l.ctx.Value("X-Src-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Src-Website")
	}

	bannedModel := db.NewOauth2BanedModel(mysql.MySQLConn)
	allow, err := bannedModel.CheckAllow(r.Context(), user.Id, web.ID, db.AllowMsg)
	if err != nil || !allow {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WebsiteNotAllow, "用户关闭了通信授权许可"),
		}, nil
	}

	err = email.SendEmail(l.ctx, req.Subject, req.Content, e.Email.String, web.Name, web.ID)
	if err != nil {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "发送邮件失败"),
			Data: types.SendMsgData{
				Success: false,
				Have:    true,
			},
		}, nil
	}

	return &types.SendMsgResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.SendMsgData{
			Success: true,
			Have:    true,
		},
	}, nil
}
