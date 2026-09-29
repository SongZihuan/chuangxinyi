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

type GetWithdrawListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWithdrawListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWithdrawListLogic {
	return &GetWithdrawListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWithdrawListLogic) GetWithdrawList(req *types.AdminGetWithdrawListReq) (resp *types.AdminGetWithdrawListResp, err error) {
	var withdrawList []db.Withdraw
	var count int64

	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	if req.ID == 0 && len(req.UID) == 0 {
		withdrawList, err = withdrawModel.GetList(l.ctx, 0, req.Status, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = withdrawModel.GetCount(l.ctx, 0, req.Status, req.Src, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	} else {
		user, err := GetUser(l.ctx, req.ID, req.UID, true)
		if errors.Is(err, UserNotFound) {
			return &types.AdminGetWithdrawListResp{
				Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
			}, nil
		} else if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		withdrawList, err = withdrawModel.GetList(l.ctx, user.WalletId, req.Status, req.Src, req.Page, req.PageSize, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}

		count, err = withdrawModel.GetCount(l.ctx, user.WalletId, req.Status, req.Src, req.StartTime, req.EndTime, req.TimeType)
		if err != nil {
			return nil, respmsg.MySQLSystemError.WarpQuick(err)
		}
	}

	userMap := make(map[int64]types.UserEasy, len(withdrawList))

	respList := make([]types.AdminWithdrawRecord, 0, len(withdrawList))
	for _, d := range withdrawList {
		optUser, ok := userMap[d.UserId]
		if !ok {
			optUser, err = action.GetUserEasy(l.ctx, d.UserId, "")
			if errors.Is(err, action.UserEasyNotFound) {
				continue
			} else if err != nil {
				return nil, respmsg.MySQLSystemError.WarpQuick(err)
			}

			userMap[d.UserId] = optUser
		}

		payAt := int64(0)
		if d.PayAt.Valid {
			payAt = d.PayAt.Time.Unix()
		}

		respList = append(respList, types.AdminWithdrawRecord{
			UserID:            d.UserId,
			WalletID:          d.WalletId,
			User:              optUser,
			WithdrawID:        d.WithdrawId,
			Name:              d.Name,
			AlipayLoginId:     d.AlipayLoginId.String,
			WechatpayNickName: d.WechatpayNickname.String,
			WechatpayOpenId:   d.WechatpayOpenId.String,
			WechatpayUnionId:  d.WechatpayUnionId.String,
			Cny:               d.Cny,
			Balance:           d.Balance.Int64,
			OrderId:           d.OrderId.String,
			PayFundOrderId:    d.PayFundOrderId.String,
			Remark:            d.Remark,
			Status:            d.Status,
			WithdrawAt:        d.WithdrawAt.Unix(),
			PayAt:             payAt,
		})
	}

	return &types.AdminGetWithdrawListResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetWithdrawListData{
			Count:    count,
			Withdraw: respList,
		},
	}, nil
}
