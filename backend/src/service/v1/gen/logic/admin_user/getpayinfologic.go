package admin_user

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/svc"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	errors "github.com/wuntsong-org/wterrors"

	"github.com/wuntsong-org/go-zero-plus/core/logx"
)

type GetPayInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPayInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPayInfoLogic {
	return &GetPayInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPayInfoLogic) GetPayInfo(req *types.AdminGetPayInfoReq) (resp *types.AdminGetPayInfoResp, err error) {
	payModel := db.NewPayModel(mysql.MySQLConn)
	pay, err := payModel.FindByPayID(l.ctx, req.ID)
	if errors.Is(err, db.ErrNotFound) {
		return &types.AdminGetPayInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PayNotFound, "订单未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	user, err := userModel.FindOneByIDWithoutDelete(l.ctx, pay.UserId)
	if errors.Is(err, db.ErrNotFound) {
		return &types.AdminGetPayInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	} else if user.WalletId != pay.WalletId {
		return &types.AdminGetPayInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.PayNotFound, "用户和订单的钱包ID不对应"),
		}, nil
	}

	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap {
		_, _ = alipay.QueryTrade(l.ctx, user, pay)
		_, _ = alipay.QueryRefund(l.ctx, pay)
	} else if pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		_, _ = wechatpay.QueryTrade(l.ctx, user, pay)
		_, _ = wechatpay.QueryRefund(l.ctx, user, pay)
	}

	payAt := int64(0)
	if pay.PayAt.Valid {
		payAt = pay.PayAt.Time.Unix()
	}

	refundAt := int64(0)
	if pay.RefundAt.Valid {
		refundAt = pay.RefundAt.Time.Unix()
	}

	optUser, err := action.GetUserEasy(l.ctx, pay.UserId, "")
	if errors.Is(err, action.UserEasyNotFound) {
		return &types.AdminGetPayInfoResp{
			Resp: respmsg.GetRespByMsg(l.ctx, respmsg.UserNotFound, "用户未找到"),
		}, nil
	} else if err != nil {
		return nil, respmsg.MySQLSystemError.WarpQuick(err)
	}

	return &types.AdminGetPayInfoResp{
		Resp: respmsg.GetRespSuccess(l.ctx),
		Data: types.AdminGetPayInfoData{
			Pay: types.AdminPayRecord{
				UserID:      pay.UserId,
				WalletID:    pay.WalletId,
				User:        optUser,
				TradeNo:     pay.TradeNo.String,
				TradeID:     pay.PayId,
				Subject:     pay.Subject,
				PayWay:      pay.PayWay,
				Cny:         pay.Cny,
				Get:         pay.Get,
				CouponsID:   pay.CouponsId.Int64,
				Balance:     pay.Balance.Int64,
				BuyerID:     pay.BuyerId.String,
				TradeStatus: pay.TradeStatus,
				Remark:      pay.Remark,
				CreateAt:    pay.CreateAt.Unix(),
				PayAt:       payAt,
				RefundAt:    refundAt,
			},
		},
	}, nil
}
