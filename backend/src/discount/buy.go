package discount

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"github.com/wuntsong-org/wterrors"
	"time"
)

var PurchaseLimit = errors.NewClass("purchase limit")
var NeedVerify = errors.NewClass("need verify")
var NeedCompany = errors.NewClass("need company")
var NeedUserOrigin = errors.NewClass("need user origin")
var NeedCompanyOrigin = errors.NewClass("need company origin")
var NeedUserFace = errors.NewClass("need user face")
var NeedCompanyFace = errors.NewClass("need company face")

func Join(ctx context.Context, user *db.User, discount *db.Discount) (*db.DiscountBuy, errors.WTError) {
	key := fmt.Sprintf("discount:%d:%d", user.Id, discount.Id)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return nil, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(key)

	idcardModel := db.NewIdcardModel(mysql.MySQLConn)
	companyModel := db.NewCompanyModel(mysql.MySQLConn)

	idcard, err := idcardModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		idcard = nil
	} else if err != nil {
		return nil, errors.WarpQuick(err)
	}

	company, err := companyModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		company = nil
	} else if err != nil {
		return nil, errors.WarpQuick(err)
	}

	if discount.NeedVerify && idcard == nil {
		return nil, NeedVerify.New()
	}

	if discount.NeedCompany && company == nil {
		return nil, NeedCompany.New()
	}

	if discount.NeedUserOrigin && (idcard == nil || !idcard.IdcardKey.Valid) {
		return nil, NeedUserOrigin.New()
	}

	if discount.NeedCompanyOrigin && (company == nil || !company.IdcardKey.Valid) {
		return nil, NeedCompanyOrigin.New()
	}

	if discount.NeedUserFace && (idcard == nil || !idcard.FaceCheckId.Valid) {
		return nil, NeedUserFace.New()
	}

	if discount.NeedCompanyFace && (company == nil || !company.FaceCheckId.Valid) {
		return nil, NeedCompanyFace.New()
	}

	buyModel := db.NewDiscountBuyModel(mysql.MySQLConn)
	buy, err := buyModel.FindOneByUserID(ctx, user.Id, discount.Id)
	if errors.Is(err, db.ErrNotFound) {
		buy = &db.DiscountBuy{
			UserId:        user.Id,
			DiscountId:    discount.Id,
			Name:          discount.Name,
			ShortDescribe: discount.ShortDescribe,
			Days:          0,
			Month:         0,
			Year:          0,
			All:           0,
		}

		res, err := buyModel.InsertWithDelete(ctx, buy)
		if err != nil {
			return nil, errors.WarpQuick(err)
		}

		buy.Id, err = res.LastInsertId()
		if err != nil {
			return nil, errors.WarpQuick(err)
		}
	} else if err != nil {
		return nil, errors.WarpQuick(err)
	}

	var isNew bool
	var newBuy *db.DiscountBuy
	now := time.Now()

	if buy.CreateAt.Year() == now.Year() && buy.CreateAt.Month() == now.Month() && buy.CreateAt.Day() == now.Day() {
		isNew = false
		newBuy = buy
		newBuy.Days += 1
		newBuy.Month += 1
		newBuy.Year += 1
		newBuy.All += 1
	} else if buy.CreateAt.Year() == now.Year() && buy.CreateAt.Month() == now.Month() {
		isNew = true
		newBuy = &db.DiscountBuy{
			UserId:        user.Id,
			DiscountId:    discount.Id,
			Name:          discount.Name,
			ShortDescribe: discount.ShortDescribe,
			Days:          1,
			Month:         buy.Month + 1,
			Year:          buy.Year + 1,
			All:           buy.All + 1,
		}
	} else if buy.CreateAt.Year() == now.Year() {
		isNew = true
		newBuy = &db.DiscountBuy{
			UserId:        user.Id,
			DiscountId:    discount.Id,
			Name:          discount.Name,
			ShortDescribe: discount.ShortDescribe,
			Days:          1,
			Month:         1,
			Year:          buy.Year + 1,
			All:           buy.All + 1,
		}
	} else {
		isNew = true
		newBuy = &db.DiscountBuy{
			UserId:        user.Id,
			DiscountId:    discount.Id,
			Name:          discount.Name,
			ShortDescribe: discount.ShortDescribe,
			Days:          1,
			Month:         1,
			Year:          1,
			All:           buy.All + 1,
		}
	}

	if discount.DayLimit.Valid && newBuy.Days > discount.DayLimit.Int64 {
		return nil, PurchaseLimit.New()
	}

	if discount.MonthLimit.Valid && newBuy.Month > discount.MonthLimit.Int64 {
		return nil, PurchaseLimit.New()
	}

	if discount.YearLimit.Valid && newBuy.Year > discount.YearLimit.Int64 {
		return nil, PurchaseLimit.New()
	}

	if discount.Limit.Valid && newBuy.All > discount.Limit.Int64 {
		return nil, PurchaseLimit.New()
	}

	if isNew {
		res, err := buyModel.InsertWithDelete(ctx, newBuy)
		if err != nil {
			return nil, errors.WarpQuick(err)
		}

		newBuy.Id, err = res.LastInsertId()
		if err != nil {
			return nil, errors.WarpQuick(err)
		}
	} else {
		err := buyModel.Update(ctx, newBuy)
		if err != nil {
			return nil, errors.WarpQuick(err)
		}
	}

	err = Process(ctx, user, discount)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return newBuy, nil
}
