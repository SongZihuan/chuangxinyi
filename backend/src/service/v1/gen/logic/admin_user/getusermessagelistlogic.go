package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMessageListLogic {
	return &GetUserMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserMessageListLogic) GetUserMessageList(req *types.AdminGetMessageListReq) (resp *types.AdminGetMessageListResp, err error) {
	var messageList []db.Message
	var count int64

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	messageModel := db.NewMessageModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		if req.ID == 0 && len(req.UID) == 0 {
			messageList, err = messageModel.GetList(l.ctx, 0, req.Src, req.JustNotRead, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = messageModel.GetCount(l.ctx, 0, req.Src, req.JustNotRead, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		} else {
			user, err := GetUser(l.ctx, req.ID, req.UID, true)
			if errors.Is(err, UserNotFound) {
				return &types.AdminGetMessageListResp{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
				}, nil
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			messageList, err = messageModel.GetList(l.ctx, user.Id, req.Src, req.JustNotRead, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = messageModel.GetCount(l.ctx, user.Id, req.Src, req.JustNotRead, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		}
	} else {
		if req.ID == 0 && len(req.UID) == 0 {
			messageList, err = messageModel.GetList(l.ctx, 0, req.Src, req.JustNotRead, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = messageModel.GetCount(l.ctx, 0, req.Src, req.JustNotRead, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		} else {
			user, err := GetUser(l.ctx, req.ID, req.UID, true)
			if errors.Is(err, UserNotFound) {
				return &types.AdminGetMessageListResp{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
				}, nil
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			messageList, err = messageModel.GetList(l.ctx, user.Id, req.Src, req.JustNotRead, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = messageModel.GetCount(l.ctx, user.Id, req.Src, req.JustNotRead, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		}
	}

	respList := make([]types.AdminMessage, 0, len(messageList))
	for _, m := range messageList {
		readAt := int64(0)
		if m.ReadAt.Valid {
			readAt = m.ReadAt.Time.Unix()
		}

		respList = append(respList, types.AdminMessage{
			ID:         m.Id,
			UserID:     m.UserId,
			Title:      m.Title,
			Content:    m.Content,
			Sender:     m.Sender,
			SenderID:   m.SenderId,
			SenderLink: m.SenderLink.String,
			CreateAt:   m.CreateAt.Unix(),
			ReadAt:     readAt,
		})
	}

	return &types.AdminGetMessageListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetMessageListData{
			Count:   count,
			Message: respList,
		},
	}, nil
}
