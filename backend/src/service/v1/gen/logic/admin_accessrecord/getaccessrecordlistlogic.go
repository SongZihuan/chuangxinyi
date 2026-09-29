package admin_accessrecord

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetAccessRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccessRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccessRecordListLogic {
	return &GetAccessRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccessRecordListLogic) GetAccessRecordList(req *types.GetAccessRecordList) (resp *types.GetAccessRecordListResp, err error) {
	var accessRecordList []db.AccessRecord
	var count int64

	accessModel := db.NewAccessRecordModel(mysql.MySQLConn)
	accessRecordList, err = accessModel.GetList(l.ctx, "", req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err = accessModel.GetCount(l.ctx, "", req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.AccessRecord, 0, len(accessRecordList))
	for _, a := range accessRecordList {
		startAt := int64(0)
		if a.StartAt.Valid {
			startAt = a.StartAt.Time.Unix()
		}

		endAt := int64(0)
		if a.EndAt.Valid {
			endAt = a.EndAt.Time.Unix()
		}

		respList = append(respList, types.AccessRecord{
			Id:                a.Id,
			RequestIdPrefix:   a.RequestIdPrefix,
			RequestId:         fmt.Sprintf("%s-%d", a.RequestIdPrefix, a.Id),
			ServerName:        a.ServerName,
			UserId:            a.UserId.Int64,
			UserUid:           a.UserUid.String,
			UserToken:         a.UserToken.String,
			RoleId:            a.RoleId.Int64,
			RoleName:          a.RoleName.String,
			RoleSign:          a.RoleSign.String,
			WebId:             a.WebId.Int64,
			WebName:           a.WebName.String,
			RequestsWebId:     a.RequestsWebId.Int64,
			RequestsWebName:   a.RequestsWebName.String,
			Ip:                a.Ip,
			GeoCode:           a.GeoCode,
			Geo:               a.Geo,
			Scheme:            a.Scheme,
			Method:            a.Method,
			Host:              a.Host,
			Path:              a.Path,
			Query:             a.Query,
			ContentType:       a.ContentType,
			RequestsBody:      a.RequestsBody,
			ResponseBody:      a.ResponseBody.String,
			ResponseBodyError: a.ResponseBodyError.String,
			RequestsHeader:    a.RequestsHeader,
			ResponseHeader:    a.ResponseHeader.String,
			StatusCode:        a.StatusCode.Int64,
			PanicError:        a.PanicError.String,
			Message:           a.Message.String,
			UseTime:           a.UseTime.Int64,
			CreateAt:          a.CreateAt.Unix(),
			StartAt:           startAt,
			EndAt:             endAt,
		})
	}

	return &types.GetAccessRecordListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetAccessRecordListData{
			Count:  count,
			Record: respList,
		},
	}, nil
}
