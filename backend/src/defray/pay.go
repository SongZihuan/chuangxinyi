package defray

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/coupons"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"github.com/wuntsong-org/go-zero-plus/core/stores/sqlx"
	"github.com/wuntsong-org/wterrors"
	"time"
)

var DefrayNotFound = errors.NewClass("defray not found")
var DoubleDefray = errors.NewClass("double defray")
var DoubleReturn = errors.NewClass("double return")
var UserNotFount = errors.NewClass("user not found")
var MustSelfDefray = errors.NewClass("must self defray")
var InsufficientQuota = errors.NewClass("insufficient quota")

func Pay(ctx context.Context, defrayID string, user *db.User, couponsID int64, token string) (*db.Defray, int64, errors.WTError) {
	key := fmt.Sprintf("defray:%s", defrayID)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return nil, 0, errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(key)

	userModel := db.NewUserModel(mysql.MySQLConn)
	defrayModel := db.NewDefrayModel(mysql.MySQLConn)
	d, err := defrayModel.FindByDefrayID(ctx, defrayID)
	if errors.Is(err, db.ErrNotFound) {
		return nil, 0, DefrayNotFound.New()
	} else if err != nil {
		return nil, 0, errors.WarpQuick(err)
	}

	if d.MustSelfDefray && d.OwnerId.Int64 != user.Id {
		return nil, 0, MustSelfDefray.New()
	}

	if d.UserId.Valid || d.Status != db.DefrayWait {
		return nil, 0, DoubleDefray.New()
	}

	u, err := userModel.FindOneByIDWithoutDelete(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return nil, 0, UserNotFount.New()
	} else if err != nil {
		return nil, 0, errors.WarpQuick(err)
	}

	if db.IsBanned(u) {
		return nil, 0, UserNotFount.New()
	}

	var realPrice int64
	if d.Price > 0 {
		if couponsID != 0 {
			realPrice, err = coupons.Defray(ctx, couponsID, d.Price)
			if err != nil {
				return nil, 0, errors.WarpQuick(err)
			}
		} else {
			realPrice = d.Price
		}

		if d.InvitePre != 0 && u.InviteId.Valid { // 有邀请人要打折
			realPrice = int64((float64(d.InvitePre) / 100) * float64(realPrice))
			if realPrice < 0 {
				realPrice = 0
			}
		}
	} else if d.Price == 0 {
		realPrice = 0
	} else {
		return nil, 0, errors.Errorf("bad price")
	}

	now := time.Now()

	d.UserId = sql.NullInt64{
		Valid: true,
		Int64: user.Id,
	}
	d.WalletId = sql.NullInt64{
		Valid: true,
		Int64: user.WalletId,
	}
	d.DefrayAt = sql.NullTime{
		Valid: true,
		Time:  now,
	}
	d.LastReturnAt = sql.NullTime{
		Valid: true,
		Time:  now.Add(time.Duration(d.ReturnDayLimit) * time.Hour * 24),
	}
	d.RealPrice = sql.NullInt64{
		Valid: true,
		Int64: realPrice,
	}

	var dt int64 // 不足的额度
	err = mysql.MySQLConn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		_, dt, err = balance.Defray(ctx, u, d, session)
		if errors.Is(err, balance.Insufficient) {
			return Insufficient.New()
		} else if err != nil {
			return errors.WarpQuick(err)
		}

		err = waitOrDistribution(d, session)

		defrayModel := db.NewDefrayModelWithSession(session)
		err = defrayModel.Update(ctx, d)
		if err != nil {
			return errors.WarpQuick(err)
		}

		return nil
	})
	if errors.Is(err, Insufficient) {
		return nil, dt, errors.WarpQuick(err)
	} else if err != nil {
		return nil, 0, errors.WarpQuick(err)
	}

	sender.PhoneSendChange(d.UserId.Int64, "余额（用户订单消费）")
	sender.EmailSendChange(d.UserId.Int64, "余额（用户订单消费）")
	sender.MessageSendPay(d.UserId.Int64, d.Price, d.Subject)
	sender.WxrobotSendPay(d.UserId.Int64, d.Price, d.Subject)
	sender.FuwuhaoSendDefray(d)
	audit.NewUserAudit(d.UserId.Int64, "用户订单消费成功（%.2f）", float64(d.Price)/100.00)

	go NotifySuccess(d.DefrayId, token)

	return d, 0, nil
}

func Return(ctx context.Context, defray *db.Defray, reason string) errors.WTError {
	key := fmt.Sprintf("defray:%s", defray.DefrayId)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(key)

	if !defray.UserId.Valid || defray.Status != db.DefraySuccess {
		return DoubleReturn.New()
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	defrayModel := db.NewDefrayModel(mysql.MySQLConn)

	user, err := userModel.FindOneByIDWithoutDelete(ctx, defray.UserId.Int64)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("error not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if defray.WalletId.Int64 != user.WalletId {
		return errors.Errorf("error not found")
	}

	defray.Status = db.DefrayWaitReturn
	defray.ReturnReason = sql.NullString{
		Valid:  true,
		String: reason,
	}

	err = defrayModel.Update(ctx, defray)
	if err != nil {
		return errors.WarpQuick(err)
	}

	logger.Logger.WXInfo("收到用户消费退款，单号：%s", defray.DefrayId)

	sender.PhoneSendChange(defray.UserId.Int64, "余额（订单退款）")
	sender.EmailSendChange(defray.UserId.Int64, "余额（订单退款）")
	sender.MessageSendPayReturn(defray.UserId.Int64, defray.Price, defray.Subject)
	sender.WxrobotSendPayReturn(defray.UserId.Int64, defray.Price, defray.Subject)
	sender.FuwuhaoSendReturnDefray(defray)
	audit.NewUserAudit(defray.UserId.Int64, "订单退款成功（%.2f）", float64(defray.Price)/100.00)

	return nil
}

func ReturnAdmin(ctx context.Context, defray *db.Defray, reason string) errors.WTError {
	key := fmt.Sprintf("defray:%s", defray.DefrayId)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(key)

	if !defray.UserId.Valid || (defray.Status != db.DefraySuccess && defray.Status != db.DefrayWaitReturn) {
		return DoubleReturn.New()
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	defrayModel := db.NewDefrayModel(mysql.MySQLConn)

	user, err := userModel.FindOneByIDWithoutDelete(ctx, defray.UserId.Int64)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("error not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if defray.WalletId.Int64 != user.WalletId {
		return errors.Errorf("error not found")
	}

	defray.Status = db.DefrayWaitReturn

	if defray.ReturnReason.Valid && (defray.ReturnReason.String == reason || len(reason) == 0) {
		defray.ReturnReason = sql.NullString{
			Valid:  true,
			String: fmt.Sprintf("管理员审批退款：%s", defray.ReturnReason.String),
		}
	} else if len(reason) == 0 {
		return errors.Errorf("bad reson")
	} else {
		defray.ReturnReason = sql.NullString{
			Valid:  true,
			String: fmt.Sprintf("管理员审批退款：%s（%s）", reason, defray.ReturnReason.String),
		}
	}

	_, err = balance.DefrayReturn(ctx, user, defray, true)
	if errors.Is(err, balance.Insufficient) {
		return InsufficientQuota.New()
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	err = defrayModel.Update(ctx, defray)
	if err != nil {
		return errors.WarpQuick(err)
	}

	sender.PhoneSendChange(defray.UserId.Int64, "余额（订单退款）")
	sender.EmailSendChange(defray.UserId.Int64, "余额（订单退款）")
	sender.MessageSendPayReturn(defray.UserId.Int64, defray.Price, defray.Subject)
	sender.WxrobotSendPayReturn(defray.UserId.Int64, defray.Price, defray.Subject)
	sender.FuwuhaoSendReturnDefray(defray)
	audit.NewUserAudit(defray.UserId.Int64, "订单退款成功（%.2f）", float64(defray.Price)/100.00)

	go NotifyReturn(defray.DefrayId)

	return nil
}

func ReturnWebsite(ctx context.Context, defray *db.Defray, reason string, must bool) errors.WTError {
	key := fmt.Sprintf("defray:%s", defray.DefrayId)
	if !redis.AcquireLockMore(ctx, key, time.Minute*2) {
		return errors.Errorf("can not get lock")
	}
	defer redis.ReleaseLock(key)

	if !defray.UserId.Valid || (defray.Status != db.DefraySuccess && defray.Status != db.DefrayWaitReturn) {
		return DoubleReturn.New()
	}

	userModel := db.NewUserModel(mysql.MySQLConn)
	defrayModel := db.NewDefrayModel(mysql.MySQLConn)

	user, err := userModel.FindOneByIDWithoutDelete(ctx, defray.UserId.Int64)
	if errors.Is(err, db.ErrNotFound) {
		return errors.Errorf("error not found")
	} else if err != nil {
		return errors.WarpQuick(err)
	} else if defray.WalletId.Int64 != user.WalletId {
		return errors.Errorf("error not found")
	}

	defray.Status = db.DefrayWaitReturn
	defray.ReturnReason = sql.NullString{
		Valid:  true,
		String: fmt.Sprintf("外站授权退款：%s", reason),
	}

	_, err = balance.DefrayReturn(ctx, user, defray, must)
	if errors.Is(err, balance.Insufficient) {
		return InsufficientQuota.New()
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	err = defrayModel.Update(ctx, defray)
	if err != nil {
		return errors.WarpQuick(err)
	}

	sender.PhoneSendChange(defray.UserId.Int64, "余额（订单退款）")
	sender.EmailSendChange(defray.UserId.Int64, "余额（订单退款）")
	sender.MessageSendPayReturn(defray.UserId.Int64, defray.Price, defray.Subject)
	sender.WxrobotSendPayReturn(defray.UserId.Int64, defray.Price, defray.Subject)
	sender.FuwuhaoSendReturnDefray(defray)
	audit.NewUserAudit(defray.UserId.Int64, "订单退款成功（%.2f）", float64(defray.Price)/100.00)

	return nil
}
