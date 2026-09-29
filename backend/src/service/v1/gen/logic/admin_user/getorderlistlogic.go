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

func (l *GetOrderListLogic) GetOrderList(req *types.AdminGetOrderListReq) (resp *types.AdminGetOrderListResp, err error) {
	var orderList []db.WorkOrder
	var count int64

	web, ok := l.ctx.Value("X-Belong-Website").(warp.Website)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Belong-Website")
	}

	orderModel := db.NewWorkOrderModel(mysql.MySQLConn)
	if web.ID == warp.UserCenterWebsite {
		if req.ID == 0 && len(req.UID) == 0 {
			orderList, err = orderModel.GetList(l.ctx, 0, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.FromID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = orderModel.GetCount(l.ctx, 0, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType, req.FromID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		} else {
			user, err := GetUser(l.ctx, req.ID, req.UID, true)
			if errors.Is(err, UserNotFound) {
				return &types.AdminGetOrderListResp{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
				}, nil
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			orderList, err = orderModel.GetList(l.ctx, user.Id, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, req.FromID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = orderModel.GetCount(l.ctx, user.Id, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType, req.FromID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		}
	} else {
		if req.ID == 0 && len(req.UID) == 0 {
			orderList, err = orderModel.GetList(l.ctx, 0, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = orderModel.GetCount(l.ctx, 0, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		} else {
			user, err := GetUser(l.ctx, req.ID, req.UID, true)
			if errors.Is(err, UserNotFound) {
				return &types.AdminGetOrderListResp{
					Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
				}, nil
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			orderList, err = orderModel.GetList(l.ctx, user.Id, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			count, err = orderModel.GetCount(l.ctx, user.Id, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType, web.ID)
			if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}
		}
	}

	respList := make([]types.AdminWorkOrder, 0, len(orderList))
	for _, w := range orderList {
		replyAt := int64(0)
		if w.LastReplyAt.Valid {
			replyAt = w.LastReplyAt.Time.Unix()
		}

		finishAt := int64(0)
		if w.FinishAt.Valid {
			finishAt = w.FinishAt.Time.Unix()
		}

		respList = append(respList, types.AdminWorkOrder{
			UserID:      w.UserId,
			OrderID:     w.Uid,
			Title:       w.Title,
			From:        w.From,
			Status:      w.Status,
			FromID:      w.FromId,
			Remark:      w.Remark,
			CreateAt:    w.CreateAt.Unix(),
			FinishAt:    finishAt,
			LastReplyAt: replyAt,
		})
	}

	return &types.AdminGetOrderListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetOrderListData{
			Count: count,
			Order: respList,
		},
	}, nil
}
