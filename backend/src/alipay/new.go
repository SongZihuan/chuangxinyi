package alipay

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/coupons"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"github.com/SuperH-0630/gopay"
	"github.com/google/uuid"
	"github.com/wuntsong-org/wterrors"
	"time"
)

const (
	PayWayPC  = "支付宝电脑支付"
	PayWayWap = "支付宝手机支付"
)

var Insufficient = errors.NewClass("insufficient")
var BadCNY = errors.NewClass("bad_cny")

func NewPagePay(ctx context.Context, user *db.User, subject string, cny int64, payMode int64, couponsID int64, returnURL string) (string, string, time.Time, errors.WTError) {
	if !config.BackendConfig.Alipay.UsePCPay {
		return "", "", time.Time{}, errors.Errorf("not ok")
	}

	if payMode != 1 && payMode != 2 && payMode != 3 {
		payMode = 2
	}

	if cny < 0 || cny > 100000 { // 小于等于0，大于等于1k
		return "", "", time.Time{}, BadCNY.New()
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		payModel := db.NewPayModel(mysql.MySQLConn)
		_, err := payModel.FindByPayID(ctx, AlipayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", "", time.Time{}, errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := AlipayID(OutTradeNoUUID.String())

	var get int64
	get = cny

	timeExpire := time.Now().Add(time.Second * time.Duration(config.BackendConfig.Coin.TimeExpireSec))

	bm := make(gopay.BodyMap)
	bm.Set("subject", subject).
		Set("out_trade_no", OutTradeNo).
		Set("total_amount", fmt.Sprintf("%.2f", float64(cny)/100.0)).
		Set("qr_pay_mode", fmt.Sprintf("%d", payMode)).
		Set("time_expire", timeExpire.Format("2006-01-02 15:04:05"))

	if len(returnURL) != 0 {
		bm.Set("return_url", returnURL)
	}

	gd := make(gopay.BodyMap)
	gd.Set("goods_id", config.BackendConfig.Coin.ID).
		Set("goods_name", config.BackendConfig.Coin.Name).
		Set("quantity", fmt.Sprintf("%d", int(cny))).
		Set("price", fmt.Sprintf("%.2f", float64(config.BackendConfig.Coin.Price)/100.00))

	bm.Set("goods_detail", []gopay.BodyMap{gd})

	url, err := AlipayClient.TradePagePay(ctx, bm)
	if err != nil {
		return "", "", time.Time{}, errors.WarpQuick(err)
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	_, err = payModel.Insert(ctx, &db.Pay{
		UserId:   user.Id,
		WalletId: user.WalletId,
		PayId:    OutTradeNo,
		Subject:  subject,
		PayWay:   PayWayPC,
		Cny:      cny,
		Get:      get,
		CouponsId: sql.NullInt64{
			Valid: couponsID != 0,
			Int64: couponsID,
		},
		TradeStatus: db.PayWait,
	})
	if err != nil {
		return "", "", time.Time{}, errors.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "支付宝支付（%s）发起成功（%.2f）", PayWayPC, float64(cny)/100.0)
	return url, OutTradeNo, timeExpire, nil
}

func NewPageWap(ctx context.Context, user *db.User, subject string, cny int64, couponsID int64, returnURL string, quiteUrl string) (string, string, time.Time, errors.WTError) {
	var err error

	if !config.BackendConfig.Alipay.UseWapPay {
		return "", "", time.Time{}, errors.Errorf("not ok")
	}

	if cny < 0 || cny > 100000 { // 小于等于0，大于等于1k
		return "", "", time.Time{}, BadCNY.New()
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		payModel := db.NewPayModel(mysql.MySQLConn)
		_, err := payModel.FindByPayID(ctx, AlipayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", "", time.Time{}, errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := AlipayID(OutTradeNoUUID.String())

	var get int64
	if couponsID != 0 {
		get, err = coupons.Recharge(ctx, couponsID, cny)
		if err != nil {
			return "", "", time.Time{}, errors.WarpQuick(err)
		}
	} else {
		get = cny
	}

	timeExpire := time.Now().Add(time.Second * time.Duration(config.BackendConfig.Coin.TimeExpireSec))

	bm := make(gopay.BodyMap)
	bm.Set("subject", subject).
		Set("out_trade_no", OutTradeNo).
		Set("total_amount", fmt.Sprintf("%.2f", float64(cny)/100.0)).
		Set("time_expire", timeExpire.Format("2006-01-02 15:04:05"))

	if len(returnURL) != 0 {
		bm.Set("return_url", returnURL)
	}

	if len(quiteUrl) != 0 {
		bm.Set("quit_url", quiteUrl)
	} else {
		bm.Set("quit_url", config.BackendConfig.Alipay.WapQuitUrl)
	}

	gd := make(gopay.BodyMap)
	gd.Set("goods_id", config.BackendConfig.Coin.ID).
		Set("goods_name", config.BackendConfig.Coin.Name).
		Set("quantity", fmt.Sprintf("%d", int(cny))).
		Set("price", fmt.Sprintf("%.2f", float64(config.BackendConfig.Coin.Price)/100.00))

	bm.Set("goods_detail", []gopay.BodyMap{gd})

	url, err := AlipayClient.TradeWapPay(ctx, bm)
	if err != nil {
		return "", "", time.Time{}, errors.WarpQuick(err)
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	_, err = payModel.Insert(ctx, &db.Pay{
		UserId:   user.Id,
		WalletId: user.WalletId,
		PayId:    OutTradeNo,
		Subject:  subject,
		PayWay:   PayWayWap,
		Cny:      cny,
		Get:      get,
		CouponsId: sql.NullInt64{
			Valid: couponsID != 0,
			Int64: couponsID,
		},
		TradeStatus: db.PayWait,
	})
	if err != nil {
		return "", "", time.Time{}, errors.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "支付宝支付（%s）发起成功（%.2f）", PayWayWap, float64(cny)/100.0)
	return url, OutTradeNo, timeExpire, nil
}

func NewRefund(ctx context.Context, user *db.User, pay *db.Pay) errors.WTError {
	if !config.BackendConfig.Alipay.UseReturnPay {
		return errors.Errorf("not ok")
	}

	if pay.PayWay != PayWayWap && pay.PayWay != PayWayPC {
		return errors.Errorf("bad pay way")
	}

	if pay.TradeStatus != db.PaySuccess && pay.TradeStatus != db.PayFinish && pay.TradeStatus != db.PayCloseRefund {
		return errors.Errorf("bad pay status")
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return errors.Errorf("can not lock pay")
	}
	defer redis.ReleaseLock(keyPay)

	payModel := db.NewPayModel(mysql.MySQLConn)

	payStatus := pay.TradeStatus
	_, err := balance.PayRefund(ctx, user, pay)
	if errors.Is(err, balance.Insufficient) {
		return Insufficient.New()
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", pay.PayId).
		Set("refund_amount", fmt.Sprintf("%.2f", float64(pay.Cny)/100.0)).
		Set("refund_reason", config.BackendConfig.Coin.RefundReason).
		Set("out_request_no", pay.PayId)

	gd := make(gopay.BodyMap)
	gd.Set("goods_id", config.BackendConfig.Coin.ID).
		Set("refund_amount", fmt.Sprintf("%.2f", float64(pay.Cny)/100.0))

	bm.Set("refund_goods_detail", []gopay.BodyMap{gd})

	bm.Set("query_options", []string{"deposit_back_info"})

	res, aliErr := AlipayClient.TradeRefund(ctx, bm)
	if aliErr != nil {
		logger.Logger.Error("支付宝支付退款失败：%s", aliErr.Error())
		_, _ = balance.PayRefundFail(ctx, user, pay, payStatus)
		return errors.WarpQuick(aliErr)
	}

	if res.Response.SubCode == "ACQ.SELLER_BALANCE_NOT_ENOUGH" {
		logger.Logger.Error("支付宝退款余额不足，请注意。退款金额：%.2f", float64(pay.Cny)/100.0)
		_, _ = balance.PayRefundFail(ctx, user, pay, payStatus)
		return errors.Errorf("system error")
	} else if res.Response.Code != "10000" {
		logger.Logger.Error("支付宝退款失败：%s (%s)", res.Response.Code, res.Response.SubCode)
		_, _ = balance.PayRefundFail(ctx, user, pay, payStatus)
		return errors.Errorf("system error")
	}

	mysqlErr := payModel.Update(ctx, pay)
	if mysqlErr != nil {
		return errors.WarpQuick(mysqlErr)
	}

	audit.NewUserAudit(pay.UserId, "支付宝退款（%s）发起成功（%.2f）", pay.PayWay, float64(pay.Cny)/100.0)
	return nil
}

func NewRefundInside(ctx context.Context, user *db.User, pay *db.Pay) errors.WTError {
	if !config.BackendConfig.Alipay.UseReturnPay {
		return errors.Errorf("not ok")
	}

	if pay.PayWay != PayWayWap && pay.PayWay != PayWayPC {
		return errors.Errorf("bad pay way")
	}

	if pay.TradeStatus != db.PaySuccess && pay.TradeStatus != db.PayFinish && pay.TradeStatus != db.PayCloseRefund {
		return errors.Errorf("bad pay status")
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return errors.Errorf("can not lock pay")
	}
	defer redis.ReleaseLock(keyPay)

	payModel := db.NewPayModel(mysql.MySQLConn)

	_, err := balance.PayRefund(ctx, user, pay)
	if errors.Is(err, balance.Insufficient) {
		return Insufficient.New()
	} else if err != nil {
		return errors.WarpQuick(err)
	}

	pay.TradeStatus = db.PaySuccessRefundInside

	mysqlErr := payModel.Update(ctx, pay)
	if mysqlErr != nil {
		return errors.WarpQuick(mysqlErr)
	}

	audit.NewUserAudit(pay.UserId, "支付宝单边退款（%s）发起成功（%.2f）", pay.PayWay, float64(pay.Cny)/100.0)
	return nil
}
