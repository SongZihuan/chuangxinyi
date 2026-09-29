package coupons

import (
	"context"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/utils"
	errors "github.com/wuntsong-org/wterrors"
)

const RechargeSend = 1 // 充值送
const FullDiscount = 2 // 满减
const FullPer = 3      // 满打折

func IsCouponsType(typeID int64) bool {
	return typeID == RechargeSend || typeID == FullDiscount || typeID == FullPer
}

type CouponsData struct {
	Bottom int64 `json:"bottom"`

	Send     int64 `json:"send"`     // Recharge
	Discount int64 `json:"discount"` // FullDiscount
	Pre      int64 `json:"pre"`      // FullPre
}

func NewRechargeSend(ctx context.Context, userID int64, name string, bottom int64, send int64) errors.WTError {
	couponsModel := db.NewCouponsModel(mysql.MySQLConn)

	content, jsonErr := utils.JsonMarshal(CouponsData{Bottom: bottom, Send: send})
	if jsonErr != nil {
		return jsonErr
	}

	_, err := couponsModel.Insert(ctx, &db.Coupons{
		UserId:  userID,
		Type:    RechargeSend,
		Name:    name,
		Content: string(content),
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func NewFullDiscount(ctx context.Context, userID int64, name string, bottom int64, discount int64) errors.WTError {
	if discount > bottom {
		return errors.Errorf("bad discount")
	}

	couponsModel := db.NewCouponsModel(mysql.MySQLConn)

	content, jsonErr := utils.JsonMarshal(CouponsData{Bottom: bottom, Discount: discount})
	if jsonErr != nil {
		return jsonErr
	}

	_, err := couponsModel.Insert(ctx, &db.Coupons{
		UserId:  userID,
		Type:    FullDiscount,
		Name:    name,
		Content: string(content),
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func NewFullPre(ctx context.Context, userID int64, name string, bottom int64, pre int64) errors.WTError {
	if pre > 100 {
		return errors.Errorf("bad pre")
	}

	couponsModel := db.NewCouponsModel(mysql.MySQLConn)

	content, jsonErr := utils.JsonMarshal(CouponsData{Bottom: bottom, Pre: pre})
	if jsonErr != nil {
		return jsonErr
	}

	_, err := couponsModel.Insert(ctx, &db.Coupons{
		UserId:  userID,
		Type:    FullPer,
		Name:    name,
		Content: string(content),
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
