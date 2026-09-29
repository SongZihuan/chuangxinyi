package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetOauthRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOauthRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOauthRecordListLogic {
	return &GetOauthRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOauthRecordListLogic) GetOauthRecordList(req *types.AdminGetOauthRecordListReq) (resp *types.AdminGetOauthRecordListResp, err error) {
	var recordList []db.Oauth2Record
	var count int64

	recordModel := db.NewOauth2RecordModel(mysql.MySQLConn)
	if req.ID == 0 && len(req.UID) == 0 {
		recordList, err = recordModel.GetList(l.ctx, 0, req.WebID, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = recordModel.GetCount(l.ctx, 0, req.WebID, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		user, err := GetUser(l.ctx, req.ID, req.UID, true)
		if errors.Is(err, UserNotFound) {
			return &types.AdminGetOauthRecordListResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		recordList, err = recordModel.GetList(l.ctx, user.Id, req.WebID, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = recordModel.GetCount(l.ctx, user.Id, req.WebID, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.AdminOauth2LoginRecord, 0, len(recordList))
	for _, a := range recordList {
		respList = append(respList, types.AdminOauth2LoginRecord{
			UserID:    a.UserId,
			WebId:     a.WebId,
			WebName:   a.WebName,
			Ip:        a.Ip,
			Geo:       a.Geo,
			GeoCode:   a.GeoCode,
			LoginTime: a.LoginTime.Unix(),
		})
	}

	return &types.AdminGetOauthRecordListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetOauthRecordListData{
			Count:  count,
			Record: respList,
		},
	}, nil
}
