package oauth2

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

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

func (l *GetOauthRecordListLogic) GetOauthRecordList(req *types.GetOauthRecordListReq) (resp *types.GetOauthRecordListResp, err error) {
	var accessRecordList []db.Oauth2Record
	var count int64

	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	recordModel := db.NewOauth2RecordModel(mysql.MySQLConn)
	accessRecordList, err = recordModel.GetList(l.ctx, user.Id, req.WebID, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err = recordModel.GetCount(l.ctx, user.Id, req.WebID, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.Oauth2LoginRecord, 0, len(accessRecordList))
	for _, a := range accessRecordList {
		respList = append(respList, types.Oauth2LoginRecord{
			WebId:     a.WebId,
			WebName:   a.WebName,
			Ip:        a.Ip,
			Geo:       a.Geo,
			GeoCode:   a.GeoCode,
			LoginTime: a.LoginTime.Unix(),
		})
	}

	return &types.GetOauthRecordListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetOauthRecordListData{
			Count:  count,
			Record: respList,
		},
	}, nil
}
