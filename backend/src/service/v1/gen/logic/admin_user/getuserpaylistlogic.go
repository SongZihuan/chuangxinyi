package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetUserPayListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPayListLogic {
	return &GetUserPayListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPayListLogic) GetUserPayList(req *types.AdminGetPayListReq) (resp *types.AdminGetPayListResp, err error) {
	var payList []db.Pay
	var count int64

	payModel := db.NewPayModel(mysql.MySQLConn)
	if req.ID == 0 && len(req.UID) == 0 {
		payList, err = payModel.GetList(l.ctx, 0, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = payModel.GetCount(l.ctx, 0, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		user, err := GetUser(l.ctx, req.ID, req.UID, true)
		if errors.Is(err, UserNotFound) {
			return &types.AdminGetPayListResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		payList, err = payModel.GetList(l.ctx, user.WalletId, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = payModel.GetCount(l.ctx, user.WalletId, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	userMap := make(map[int64]types.UserEasy, len(payList))

	respList := make([]types.AdminPayRecord, 0, len(payList))
	for _, p := range payList {
		optUser, ok := userMap[p.UserId]
		if !ok {
			optUser, err = action.GetUserEasy(l.ctx, p.UserId, "")
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			userMap[p.UserId] = optUser
		}

		payAt := int64(0)
		if p.PayAt.Valid {
			payAt = p.PayAt.Time.Unix()
		}

		refundAt := int64(0)
		if p.RefundAt.Valid {
			refundAt = p.RefundAt.Time.Unix()
		}

		respList = append(respList, types.AdminPayRecord{
			UserID:      p.UserId,
			User:        optUser,
			WalletID:    p.WalletId,
			TradeNo:     p.TradeNo.String,
			TradeID:     p.PayId,
			Subject:     p.Subject,
			PayWay:      p.PayWay,
			Cny:         p.Cny,
			Get:         p.Get,
			CouponsID:   p.CouponsId.Int64,
			Balance:     p.Balance.Int64,
			BuyerID:     p.BuyerId.String,
			TradeStatus: p.TradeStatus,
			Remark:      p.Remark,
			CreateAt:    p.CreateAt.Unix(),
			PayAt:       payAt,
			RefundAt:    refundAt,
		})
	}

	return &types.AdminGetPayListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetPayListData{
			Count: count,
			Pay:   respList,
		},
	}, nil
}
