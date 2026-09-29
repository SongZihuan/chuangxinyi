package discount

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/back"
	"gitee.com/wuntsong-auth/backend/src/coupons"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

const SendAmount = 1  // 送额度
const SendCoupons = 2 // 送优惠券

func IsDiscountType(t int64) bool {
	return t == SendAmount || t == SendCoupons
}

type DiscountQuota struct {
	Amount      int64 `json:"amount"`
	CanWithdraw bool  `json:"canWithdraw"`

	Type     int64 `json:"type"`     // RechargeSend=>1   FulIDiscount =>2   FullPer =>3
	Bottom   int64 `json:"bottom"`   //满足的金额要求
	Send     int64 `json:"send"`     //送的额度
	Discount int64 `json:"discount"` //减的额度
	Pre      int64 `json:"pre"`      //打折
}

func Process(ctx context.Context, user *db.User, discount *db.Discount) errors.WTError {
	data := DiscountQuota{}
	err := utils.JsonUnmarshal([]byte(discount.Quota), &data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	if discount.Type == SendAmount {
		_, err = back.NewBack(ctx, data.Amount, "优惠包返现", discount.Name, user, data.CanWithdraw, warp.UserCenterWebsite)
		if err != nil {
			return errors.WarpQuick(err)
		}
	} else if discount.Type == SendCoupons {
		switch data.Type {
		case coupons.RechargeSend:
			err = coupons.NewRechargeSend(ctx, user.Id, discount.Name, data.Bottom, data.Send)
		case coupons.FullDiscount:
			err = coupons.NewFullDiscount(ctx, user.Id, discount.Name, data.Bottom, data.Discount)
		case coupons.FullPer:
			err = coupons.NewFullPre(ctx, user.Id, discount.Name, data.Bottom, data.Pre)
		default:
			err = nil
		}
	}

	return nil
}
