package admin_accessrecord

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetTokenRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenRecordListLogic {
	return &GetTokenRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenRecordListLogic) GetTokenRecordList(req *types.GetTokenRecordList) (resp *types.GetTokenRecordListResp, err error) {
	var recordList []db.TokenRecord
	var count int64

	accessModel := db.NewTokenRecordModel(mysql.MySQLConn)
	recordList, err = accessModel.GetList(l.ctx, "", req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err = accessModel.GetCount(l.ctx, "", req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	respList := make([]types.TokenRecord, 0, len(recordList))
	for _, r := range recordList {
		respList = append(respList, types.TokenRecord{
			TokenType: r.TokenType,
			Token:     r.Token,
			Type:      r.Type,
			Data:      r.Data,
			CreateAt:  r.CreateAt.Unix(),
		})
	}

	return &types.GetTokenRecordListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetTokenRecordListData{
			Count:  count,
			Record: respList,
		},
	}, nil
}
