package msg

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/fuwuhao"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	utils2 "gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/utils"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"

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

func (l *SendFuwuhaoLogic) SendFuwuhao(req *types.SendFuwuhaoReq, r *http.Request) (resp *types.SendMsgResp, err error) {
	wechatModel := db.NewWechatModel(mysql.MySQLConn)

	user, err := utils2.FindUser(l.ctx, req.UserID, false)
	if errors.Is(err, utils2.UserNotFound) {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
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

	w, err := wechatModel.FindByUserID(l.ctx, user.Id)
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
	} else if !w.Fuwuhao.Valid {
		return &types.SendMsgResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.SendMsgData{
				Success: false,
				Have:    false,
			},
		}, nil
	}

	val := make(map[string]string, len(req.Val))
	for _, v := range req.Val {
		val[v.Label] = v.Value
	}

	err = fuwuhao.SendVal(l.ctx, req.TemplateID, req.Url, w.Fuwuhao.String, val, web.ID)
	if err != nil {
		logger.Logger.Error("send fuwuhao msg error: %s", err.Error())
		return &types.SendMsgResp{
			Resp: respmsg.GetRespByError(l.ctx, respmsg.Success, errors.WarpQuick(err), "发送服务号模板消息失败"),
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
