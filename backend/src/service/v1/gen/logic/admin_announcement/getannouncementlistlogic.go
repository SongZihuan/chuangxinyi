package admin_announcement

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAnnouncementListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAnnouncementListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnnouncementListLogic {
	return &GetAnnouncementListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAnnouncementListLogic) GetAnnouncementList(req *types.AdminGetAnnouncementList) (resp *types.AdminGetAnnouncementListResp, err error) {
	announcementModel := db.NewAnnouncementModel(mysql.MySQLConn)
	announcementList, err := announcementModel.GetList(l.ctx, req.Src, false, req.Page, req.PageSize)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := announcementModel.GetCount(l.ctx, req.Src, false)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.AdminAnnouncement, 0, len(announcementList))
	for _, a := range announcementList {
		respList = append(respList, types.AdminAnnouncement{
			ID:      a.Id,
			Title:   a.Title,
			Content: a.Content,
			StartAt: a.StartAt.Unix(),
			StopAt:  a.StopAt.Unix(),
			Sort:    a.Sort,
		})
	}

	return &types.AdminGetAnnouncementListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetAnnouncementListData{
			Count:        count,
			Announcement: respList,
		},
	}, nil
}
