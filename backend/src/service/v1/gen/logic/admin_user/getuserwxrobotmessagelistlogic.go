package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserWxrobotMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserWxrobotMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWxrobotMessageListLogic {
	return &GetUserWxrobotMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserWxrobotMessageListLogic) GetUserWxrobotMessageList(req *types.AdminGetWxrobotMessageListReq) (resp *types.AdminGetWxrobotMessageListResp, err error) {
	var wxrobotMessage []db.WxrobotMessage
	var count int64

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	wxrobotMessageModel := db.NewWxrobotMessageModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		wxrobotMessage, err = wxrobotMessageModel.GetList(l.ctx, req.Webhook, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = wxrobotMessageModel.GetCount(l.ctx, req.Webhook, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		wxrobotMessage, err = wxrobotMessageModel.GetList(l.ctx, req.Webhook, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = wxrobotMessageModel.GetCount(l.ctx, req.Webhook, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.AdminWxrobotMessage, 0, len(wxrobotMessage))
	for _, w := range wxrobotMessage {
		respList = append(respList, types.AdminWxrobotMessage{
			Webhook:  w.Webhook,
			Text:     w.Text,
			AtAll:    w.AtAll,
			SenderId: w.SenderId,
			Success:  w.Success,
			ErrorMsg: w.ErrorMsg.String,
			CreateAt: w.CreateAt.Unix(),
		})
	}

	return &types.AdminGetWxrobotMessageListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetWxrobotMessageListData{
			Count:   count,
			Message: respList,
		},
	}, nil
}
