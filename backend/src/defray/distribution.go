package defray

import (
	"context"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/back"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

func startDistribution() errors.WTError {
	go func() {
		for {
			func() {
				defer utils.Recover(logger.Logger, nil, "distribution fail")

				key := "defray:distribution"
				if !redis.AcquireLock(context.Background(), key, time.Minute*2) {
					return
				}
				defer redis.ReleaseLock(key)

				defrayModel := db.NewDefrayModel(mysql.MySQLConn)
				waitList, err := defrayModel.GetWaitDistribution(context.Background(), 1000)
				if err != nil {
					logger.Logger.Error("mysql resp: %s", err.Error())
					return
				}

				for _, d := range waitList {
					go func(id string) {
						key := fmt.Sprintf("defray:%s", id)
						if !redis.AcquireLockMore(context.Background(), key, time.Minute*2) {
							return
						}
						defer redis.ReleaseLock(key)

						err := mysql.MySQLConn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
							defrayModel := db.NewDefrayModelWithSession(session)
							d, err := defrayModel.FindByDefrayID(context.Background(), id)
							if err != nil {
								return err
							}

							err = distribution(d, session)
							if err != nil {
								return err
							}

							err = defrayModel.Update(context.Background(), d)
							if err != nil {
								return errors.WarpQuick(err)
							}

							return nil
						})
						if err != nil {
							logger.Logger.Error("distribution fail: %s", err.Error())
							return
						}
					}(d.DefrayId)
				}

			}()

			time.Sleep(30 * time.Second)
		}
	}()

	return nil
}

func distribution(d *db.Defray, session sqlx.Session) errors.WTError {
	if !d.UserId.Valid || !d.LastReturnAt.Valid || d.LastReturnAt.Time.After(time.Now()) { // 没到分销时间
		return errors.Errorf("not defray")
	}

	if d.HasDistribution || d.ReturnAt.Valid { // 已经分销或者不能分销
		return nil
	}

	userModel := db.NewUserModelWithSession(session)
	user, err := userModel.FindOneByIDWithoutDelete(context.Background(), d.UserId.Int64)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if user != nil {
		return errors.WarpQuick(err)
	}

	if d.Price > 0 {
		err := _distribution(context.Background(), session, d, user, d.Price, d.DistributionLevel3, d.DistributionLevel2, d.DistributionLevel1)
		if err != nil {
			return err
		}
	}

	d.HasDistribution = true
	return nil
}

func _distribution(ctx context.Context, mysql sqlx.Session, defray *db.Defray, user *db.User, price int64, pre ...int64) errors.WTError {
	if len(pre) == 0 {
		return nil
	}

	if !user.InviteId.Valid {
		return nil
	}

	userModel := db.NewUserModelWithSession(mysql)

	inviteUser, err := userModel.FindOneByIDWithoutDelete(ctx, user.InviteId.Int64)
	if errors.Is(err, db.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if db.IsBanned(inviteUser) {
		return nil
	}

	var p int64
	if pre[0] > 0 {
		p = int64((float64(pre[0]) / 100) * float64(price))
	} else {
		p = 0
	}

	if p > 0 {
		_, err = back.NewBackWithSession(ctx, p, "分销收益", fmt.Sprintf("消费（%s）的分销收益", defray.Subject), inviteUser, defray.CanWithdraw, defray.SupplierId, mysql)
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return _distribution(ctx, mysql, defray, inviteUser, price, pre[1:]...)
}

func waitDistribution(d *db.Defray, session sqlx.Session) errors.WTError {
	userModel := db.NewUserModelWithSession(mysql.MySQLConn)
	user, err := userModel.FindOneByIDWithoutDelete(context.Background(), d.UserId.Int64)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("user not found")
	} else if user != nil {
		return errors.WarpQuick(err)
	}

	if d.Price > 0 {
		err := _waitDistribution(context.Background(), session, d, user, d.Price, d.DistributionLevel3, d.DistributionLevel2, d.DistributionLevel1)
		if err != nil {
			logger.Logger.Error("distribution error: %s", err.Error())
		}
	}

	return nil
}

func _waitDistribution(ctx context.Context, session sqlx.Session, defray *db.Defray, user *db.User, price int64, pre ...int64) errors.WTError {
	if len(pre) == 0 {
		return nil
	}

	if !user.InviteId.Valid {
		return nil
	}

	userModel := db.NewUserModelWithSession(mysql.MySQLConn)

	inviteUser, err := userModel.FindOneByIDWithoutDelete(ctx, user.InviteId.Int64)
	if errors.Is(err, db.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if db.IsBanned(inviteUser) {
		return nil
	}

	var p int64
	if pre[0] > 0 {
		p = int64((float64(pre[0]) / 100) * float64(price))
	} else {
		p = 0
	}

	if p > 0 {
		_, err = balance.WaitBack(ctx, session, inviteUser, defray.CanWithdraw, p, fmt.Sprintf("消费（%s）的分销收益", defray.Subject))
		if err != nil {
			return errors.WarpQuick(err)
		}
	}

	return _waitDistribution(ctx, session, defray, inviteUser, price, pre[1:]...)
}

func waitOrDistribution(d *db.Defray, session sqlx.Session) errors.WTError {
	if d.Price <= 0 {
		d.HasDistribution = true
		return nil
	}

	if d.ReturnDayLimit == 0 {
		return distribution(d, session)
	} else {
		return waitDistribution(d, session)
	}
}
