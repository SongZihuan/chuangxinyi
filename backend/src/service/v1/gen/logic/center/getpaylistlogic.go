package center

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

type GetPayListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPayListLogic {
	return &GetPayListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPayListLogic) GetPayList(req *types.GetPayListReq) (resp *types.GetPayListResp, err error) {
	user, ok := l.ctx.Value("X-Token-User").(*db.User)
	if !ok {
		return nil, respmsg.BadContextError.New("X-Token-User")
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	payList, err := payModel.GetList(l.ctx, user.WalletId, req.Src, req.Status, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	count, err := payModel.GetCount(l.ctx, user.WalletId, req.Src, req.Status, req.StartTime, req.EndTime, req.TimeType)
	if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userMap := make(map[int64]types.UserEasy, len(payList))

	respList := make([]types.PayRecord, 0, len(payList))
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

		respList = append(respList, types.PayRecord{
			User:        optUser,
			TradeNo:     p.TradeNo.String,
			TradeID:     p.PayId,
			Subject:     p.Subject,
			PayWay:      p.PayWay,
			Cny:         p.Cny,
			Get:         p.Get,
			HasCoupons:  p.CouponsId.Valid,
			Balance:     p.Balance.Int64,
			TradeStatus: p.TradeStatus,
			CreateAt:    p.CreateAt.Unix(),
			PayAt:       payAt,
			RefundAt:    refundAt,
		})
	}

	return &types.GetPayListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.GetPayListData{
			Count: count,
			Pay:   respList,
		},
	}, nil
}
