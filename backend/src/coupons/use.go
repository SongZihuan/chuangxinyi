package coupons

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/wuntsong-org/wterrors"
	"time"
)

func Recharge(ctx context.Context, couponsID int64, cny int64) (int64, errors.WTError) {
	key := fmt.Sprintf("coupons:%d", couponsID)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(key)

	couponsModel := db.NewCouponsModel(mysql.MySQLConn)
	c, err := couponsModel.FindOneWithoutDelete(ctx, couponsID)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("coupons not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	} else if c.Type != RechargeSend {
		return 0, errors.Errorf("bad coupons")
	}

	data := CouponsData{}
	err = utils.JsonUnmarshal([]byte(c.Content), &data)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	if data.Bottom > cny {
		return 0, errors.Errorf("condition not met")
	}

	c.DeleteAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	err = couponsModel.Update(ctx, c)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	return cny + data.Send, nil
}

func Defray(ctx context.Context, couponsID int64, cny int64) (int64, errors.WTError) {
	key := fmt.Sprintf("coupons:%d", couponsID)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(key)

	couponsModel := db.NewCouponsModel(mysql.MySQLConn)
	c, err := couponsModel.FindOneWithoutDelete(ctx, couponsID)
	if errors.Is(err, db.ErrNotFound) {
		return 0, errors.Errorf("coupons not found")
	} else if err != nil {
		return 0, errors.WarpQuick(err)
	}

	data := CouponsData{}
	err = utils.JsonUnmarshal([]byte(c.Content), &data)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	if data.Bottom > cny {
		return 0, errors.Errorf("condition not met")
	}

	newCny := cny
	switch c.Type {
	case FullDiscount:
		newCny -= data.Discount
	case FullPer:
		newCny = int64((float64(data.Pre) / 100.0) * float64(newCny))
	default:
		return 0, errors.Errorf("bad coupons")
	}

	if newCny < 0 {
		newCny = 0
	} else if newCny > cny {
		newCny = cny
	}

	c.DeleteAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	err = couponsModel.Update(ctx, c)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}

	return newCny, nil
}
