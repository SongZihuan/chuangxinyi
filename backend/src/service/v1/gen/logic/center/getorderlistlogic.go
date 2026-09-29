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

type GetOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderListLogic) GetOrderList(req *types.GetOrderListReq) (resp *types.GetOrderListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	web, ok := l.ctx.Value("X-Token-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-Website")
	}

	var workOrderList []db.WorkOrder
	var count int64

	workOrderModel := db.NewWorkOrderModel(mysql.MySQLConn)

	if web.ID == warp.UserCenterWebsite {
		workOrderList, err = workOrderModel.GetList(l.ctx, user.Id, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.FromID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = workOrderModel.GetCount(l.ctx, user.Id, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType, req.FromID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		workOrderList, err = workOrderModel.GetList(l.ctx, user.Id, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = workOrderModel.GetCount(l.ctx, user.Id, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType, web.ID)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	respList := make([]types.WorkOrder, 0, len(workOrderList))
	for _, w := range workOrderList {
		replyAt := int64(0)
		if w.LastReplyAt.Valid {
			replyAt = w.LastReplyAt.Time.Unix()
		}

		finishAt := int64(0)
		if w.FinishAt.Valid {
			finishAt = w.FinishAt.Time.Unix()
		}

		respList = append(respList, types.WorkOrder{
			OrderID:     w.Uid,
			Title:       w.Title,
			From:        w.From,
			Status:      w.Status,
			CreateAt:    w.CreateAt.Unix(),
			FinishAt:    finishAt,
			LastReplyAt: replyAt,
		})
	}

	return &types.GetOrderListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetOrderListData{
			Count: count,
			Order: respList,
		},
	}, nil
}
