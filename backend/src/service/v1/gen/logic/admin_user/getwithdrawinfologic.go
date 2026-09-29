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

type GetWithdrawInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWithdrawInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWithdrawInfoLogic {
	return &GetWithdrawInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWithdrawInfoLogic) GetWithdrawInfo(req *types.AdminGetWithdrawInfoReq) (resp *types.AdminGetWithdrawInfoResp, err error) {
	withdrawModel := db.NewWithdrawModel(mysql.MySQLConn)
	withdraw, err := withdrawModel.FindByWithdrawID(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.AdminGetWithdrawInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.WithdrawNotFound, "提现订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	optUser, err := action.GetUserEasy(l.ctx, withdraw.UserId, "")
	if errors.Is(err, action.UserEasyNotFound) {
		return &types.AdminGetWithdrawInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	payAt := int64(0)
	if withdraw.PayAt.Valid {
		payAt = withdraw.PayAt.Time.Unix()
	}

	return &types.AdminGetWithdrawInfoResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetWithdrawInfoData{
			Withdraw: types.AdminWithdrawRecord{
				UserID:            withdraw.UserId,
				WalletID:          withdraw.WalletId,
				User:              optUser,
				WithdrawID:        withdraw.WithdrawId,
				Name:              withdraw.Name,
				AlipayLoginId:     withdraw.AlipayLoginId.String,
				WechatpayNickName: withdraw.WechatpayNickname.String,
				WechatpayOpenId:   withdraw.WechatpayOpenId.String,
				WechatpayUnionId:  withdraw.WechatpayUnionId.String,
				Cny:               withdraw.Cny,
				Balance:           withdraw.Balance.Int64,
				OrderId:           withdraw.OrderId.String,
				PayFundOrderId:    withdraw.PayFundOrderId.String,
				Remark:            withdraw.Remark,
				Status:            withdraw.Status,
				WithdrawAt:        withdraw.WithdrawAt.Unix(),
				PayAt:             payAt,
			},
		},
	}, nil
}
