package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserFuwuhaoMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserFuwuhaoMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFuwuhaoMessageListLogic {
	return &GetUserFuwuhaoMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserFuwuhaoMessageListLogic) GetUserFuwuhaoMessageList(req *types.AdminGetFuwuhaoMessageListReq) (resp *types.AdminGetFuwuhaoMessageListResp, err error) {
	var fuwuhaoMessage []db.FuwuhaoMessage
	var count int64

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	fuwuhaoMessageModel := db.NewFuwuhaoMessageModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		fuwuhaoMessage, err = fuwuhaoMessageModel.GetList(l.ctx, req.OpenID, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = fuwuhaoMessageModel.GetCount(l.ctx, req.OpenID, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		fuwuhaoMessage, err = fuwuhaoMessageModel.GetList(l.ctx, req.OpenID, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = fuwuhaoMessageModel.GetCount(l.ctx, req.OpenID, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.AdminFuwuhaoMessage, 0, len(fuwuhaoMessage))
	for _, f := range fuwuhaoMessage {
		var val map[string]interface{}
		err = utils.JsonUnmarshal([]byte(f.Val), &val)
		if err != nil {
			val = map[string]interface{}{}
			logger.Logger.Error("utils.JsonUnmarshal error: %s", err)
		}

		valList := make([]types.LabelInterfaceValueRecord, 0, len(val))
		for l, v := range val {
			valList = append(valList, types.LabelInterfaceValueRecord{
				Label: l,
				Value: v,
			})
		}

		respList = append(respList, types.AdminFuwuhaoMessage{
			OpenID:   f.OpenId,
			Template: f.Template,
			Url:      f.Url,
			Val:      valList,
			SenderId: f.SenderId,
			Success:  f.Success,
			ErrorMsg: f.ErrorMsg.String,
			CreateAt: f.CreateAt.Unix(),
		})
	}

	return &types.AdminGetFuwuhaoMessageListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetFuwuhaoMessageListData{
			Count:   count,
			Message: respList,
		},
	}, nil
}
