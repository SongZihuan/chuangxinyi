package center

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

type GetMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageListLogic) GetMessageList(req *types.GetMessageListReq) (resp *types.GetMessageListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	var messageList []db.Message
	var count int64

	messageModel := db.NewMessageModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		messageList, err = messageModel.GetList(l.ctx, user.Id, req.Src, req.JustRead, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = messageModel.GetCount(l.ctx, user.Id, req.Src, req.JustRead, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		messageList, err = messageModel.GetList(l.ctx, user.Id, req.Src, req.JustRead, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = messageModel.GetCount(l.ctx, user.Id, req.Src, req.JustRead, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.Message, 0, len(messageList))
	for _, m := range messageList {
		readAt := int64(0)
		if m.ReadAt.Valid {
			readAt = m.ReadAt.Time.Unix()
		}

		respList = append(respList, types.Message{
			ID:         m.Id,
			Title:      m.Title,
			Content:    m.Content,
			Sender:     m.Sender,
			SenderLink: m.SenderLink.String,
			CreateAt:   m.CreateAt.Unix(),
			ReadAt:     readAt,
		})
	}

	return &types.GetMessageListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetMessageListData{
			Count:   count,
			Message: respList,
		},
	}, nil
}
