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

type GetUserEmailMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserEmailMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmailMessageListLogic {
	return &GetUserEmailMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserEmailMessageListLogic) GetUserEmailMessageList(req *types.AdminGetEmailMessageListReq) (resp *types.AdminGetEmailMessageListResp, err error) {
	var emailMessageList []db.EmailMessage
	var count int64

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	emailMessageModel := db.NewEmailMessageModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		emailMessageList, err = emailMessageModel.GetList(l.ctx, req.Email, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = emailMessageModel.GetCount(l.ctx, req.Email, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		emailMessageList, err = emailMessageModel.GetList(l.ctx, req.Email, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = emailMessageModel.GetCount(l.ctx, req.Email, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.AdminEmailMessage, 0, len(emailMessageList))
	for _, e := range emailMessageList {
		respList = append(respList, types.AdminEmailMessage{
			Email:    e.Email,
			Subject:  e.Subject,
			Content:  e.Content,
			Sender:   e.Sender,
			SenderId: e.SenderId,
			Success:  e.Success,
			ErrorMsg: e.ErrorMsg.String,
			CreateAt: e.CreateAt.Unix(),
		})
	}

	return &types.AdminGetEmailMessageListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetEmailMessageListData{
			Count:   count,
			Message: respList,
		},
	}, nil
}
