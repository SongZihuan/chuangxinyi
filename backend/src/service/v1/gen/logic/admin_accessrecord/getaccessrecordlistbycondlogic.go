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

type GetAccessRecordListByCondLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccessRecordListByCondLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccessRecordListByCondLogic {
	return &GetAccessRecordListByCondLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccessRecordListByCondLogic) GetAccessRecordListByCond(req *types.GetAccessRecordListByCond) (resp *types.GetAccessRecordListByCondResp, err error) {
	accessModel := db.NewAccessRecordModel(mysql.MySQLConn)
	accessRecordList, listQuery, err := accessModel.GetListByCond(l.ctx, req.Cond, req.Page, req.PageSize)
	if err != nil {
		fmt.Println("TAG A")
		return &types.GetAccessRecordListByCondResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetAccessRecordListByCondData{
				ListQuery: listQuery,
				ListError: err.Error(),
			},
		}, nil
	}

	count, countQuey, err := accessModel.GetCountByCond(l.ctx, req.Cond)
	if err != nil {
		fmt.Println("TAG B")
		return &types.GetAccessRecordListByCondResp{
			Resp: respmsg.GetRespSuccess(l.ctx),
			Data: types.GetAccessRecordListByCondData{
				ListQuery:  listQuery,
				CountQuery: countQuey,
				CountError: err.Error(),
			},
		}, nil
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

	fmt.Println("TAG C")
	return &types.GetAccessRecordListByCondResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetAccessRecordListByCondData{
			ListQuery:  listQuery,
			CountQuery: countQuey,
			Count:      count,
			Record:     respList,
		},
	}, nil
}
