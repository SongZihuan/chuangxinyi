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

type GetUserSmsMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserSmsMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSmsMessageListLogic {
	return &GetUserSmsMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSmsMessageListLogic) GetUserSmsMessageList(req *types.AdminGetSmsMessageListReq) (resp *types.AdminGetSmsMessageListResp, err error) {
	var smsMessage []db.SmsMessage
	var count int64

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	smsMessageModel := db.NewSmsMessageModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		smsMessage, err = smsMessageModel.GetList(l.ctx, req.Phone, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = smsMessageModel.GetCount(l.ctx, req.Phone, req.StartTime, req.EndTime, req.TimeType, req.SenderID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		smsMessage, err = smsMessageModel.GetList(l.ctx, req.Phone, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = smsMessageModel.GetCount(l.ctx, req.Phone, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.AdminSmsMessage, 0, len(smsMessage))
	for _, s := range smsMessage {
		var tp map[string]interface{}
		err = utils.JsonUnmarshal([]byte(s.TemplateParam), &tp)
		if err != nil {
			tp = map[string]interface{}{}
			logger.Logger.Error("utils.JsonUnmarshal error: %s", err)
		}

		tpList := make([]types.LabelInterfaceValueRecord, 0, len(tp))
		for l, v := range tp {
			tpList = append(tpList, types.LabelInterfaceValueRecord{
				Label: l,
				Value: v,
			})
		}

		respList = append(respList, types.AdminSmsMessage{
			Phone:         s.Phone,
			Sig:           s.Sig,
			Template:      s.Template,
			TemplateParam: tpList,
			SenderId:      s.SenderId,
			Success:       s.Success,
			ErrorMsg:      s.ErrorMsg.String,
			CreateAt:      s.CreateAt.Unix(),
		})
	}

	return &types.AdminGetSmsMessageListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetSmsMessageListData{
			Count:   count,
			Message: respList,
		},
	}, nil
}
