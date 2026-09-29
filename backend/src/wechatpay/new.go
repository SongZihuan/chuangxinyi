package wechatpay

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
	PayWayNative = "微信Native支付"
	PayWayH5     = "微信H5支付"
	PayWayJSAPI  = "微信JSAPI支付"
)

var Insufficient = errors.NewClass("insufficient") // 余额不足
var BadCNY = errors.NewClass("bad cny")
var MustBindFuwuhao = errors.NewClass("must bind fuwuhao")

func NewPagePay(ctx context.Context, user *db.User, subject string, cny int64, couponsID int64) (string, string, time.Time, errors.WTError) {
	if !config.BackendConfig.WeChatPay.UseNativePay {
		return "", "", time.Time{}, errors.Errorf("not ok")
	}

	if cny < 0 || cny > 100000 { // 小于等于0，大于等于1k
		return "", "", time.Time{}, BadCNY.New()
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		payModel := db.NewPayModel(mysql.MySQLConn)
		_, err := payModel.FindByPayID(ctx, WeChatPayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", "", time.Time{}, errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := WeChatPayID(OutTradeNoUUID.String())

	var get int64
	get = cny

	timeExpire := time.Now().Add(time.Second * time.Duration(config.BackendConfig.Coin.TimeExpireSec))

	bm := make(gopay.BodyMap)
	bm.Set("appid", config.BackendConfig.WeChatPay.AppID).
		Set("description", subject).
		Set("out_trade_no", OutTradeNo).
		Set("time_expire", timeExpire.Format(time.RFC3339)).
		Set("notify_url", config.BackendConfig.WeChatPay.ReturnURL)

	am := make(gopay.BodyMap)
	am.Set("total", cny).Set("currency", "CNY")

	bm.Set("amount", am)

	gd := make(gopay.BodyMap)
	gd.Set("goods_name", config.BackendConfig.Coin.Name).
		Set("merchant_goods_id", config.BackendConfig.Coin.ID).
		Set("quantity", cny).
		Set("unit_price", config.BackendConfig.Coin.Price)

	d := make(gopay.BodyMap)
	d.Set("goods_detail", []gopay.BodyMap{gd})

	bm.Set("detail", d)

	res, err := WeChatPayClient.V3TransactionNative(ctx, bm)
	if err != nil {
		return "", "", time.Time{}, errors.WarpQuick(err)
	} else if res.Code != 0 {
		return "", "", time.Time{}, errors.Errorf("%s", res.Error)
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	_, err = payModel.Insert(ctx, &db.Pay{
		UserId:   user.Id,
		WalletId: user.WalletId,
		PayId:    OutTradeNo,
		Subject:  subject,
		PayWay:   PayWayNative,
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

	audit.NewUserAudit(user.Id, "微信支付（%s）发起成功（%.2f）", PayWayNative, float64(cny)/100.0)
	return res.Response.CodeUrl, OutTradeNo, timeExpire, nil
}

func NewPageH5(ctx context.Context, user *db.User, subject string, cny int64, couponsID int64, h5Type string, ip string) (string, string, time.Time, errors.WTError) {
	var err error
	if !config.BackendConfig.WeChatPay.UseH5Pay {
		return "", "", time.Time{}, errors.Errorf("not ok")
	}

	if cny < 0 || cny > 100000 { // 小于等于0，大于等于1k
		return "", "", time.Time{}, BadCNY.New()
	}

	if h5Type != "Android" && h5Type != "iOS" {
		h5Type = "Wap"
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		payModel := db.NewPayModel(mysql.MySQLConn)
		_, err := payModel.FindByPayID(ctx, WeChatPayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", "", time.Time{}, errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := WeChatPayID(OutTradeNoUUID.String())

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
	bm.Set("appid", config.BackendConfig.WeChatPay.AppID).
		Set("description", subject).
		Set("out_trade_no", OutTradeNo).
		Set("time_expire", timeExpire.Format(time.RFC3339)).
		Set("notify_url", config.BackendConfig.WeChatPay.ReturnURL)

	am := make(gopay.BodyMap)
	am.Set("total", cny).Set("currency", "CNY")

	bm.Set("amount", am)

	gd := make(gopay.BodyMap)
	gd.Set("goods_name", config.BackendConfig.Coin.Name).
		Set("merchant_goods_id", config.BackendConfig.Coin.ID).
		Set("quantity", cny).
		Set("unit_price", config.BackendConfig.Coin.Price)

	d := make(gopay.BodyMap)
	d.Set("goods_detail", []gopay.BodyMap{gd})

	bm.Set("detail", d)

	s := make(gopay.BodyMap)
	s.Set("payer_client_ip", ip)

	hi := make(gopay.BodyMap)
	hi.Set("type", h5Type)

	s.Set("h5_info", hi)
	bm.Set("scene_info", s)

	res, err := WeChatPayClient.V3TransactionH5(ctx, bm)
	if err != nil {
		return "", "", time.Time{}, errors.WarpQuick(err)
	} else if res.Code != 0 {
		return "", "", time.Time{}, errors.Errorf("%s", res.Error)
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	_, err = payModel.Insert(ctx, &db.Pay{
		UserId:   user.Id,
		WalletId: user.WalletId,
		PayId:    OutTradeNo,
		Subject:  subject,
		PayWay:   PayWayH5,
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

	audit.NewUserAudit(user.Id, "微信支付（%s）发起成功（%.2f）", PayWayH5, float64(cny)/100.0)
	return res.Response.H5Url, OutTradeNo, timeExpire, nil
}

type JSAPIPayParams struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

func NewPageJsAPI(ctx context.Context, user *db.User, subject string, cny int64, couponsID int64) (string, JSAPIPayParams, string, time.Time, errors.WTError) {
	var err error

	if !config.BackendConfig.WeChatPay.UseH5Pay {
		return "", JSAPIPayParams{}, "", time.Time{}, errors.Errorf("not ok")
	}

	if cny < 0 || cny > 100000 { // 小于等于0，大于等于1k
		return "", JSAPIPayParams{}, "", time.Time{}, BadCNY.New()
	}

	wechatModel := db.NewWechatModel(mysql.MySQLConn)

	u, err := wechatModel.FindByUserID(ctx, user.Id)
	if errors.Is(err, db.ErrNotFound) {
		return "", JSAPIPayParams{}, "", time.Time{}, MustBindFuwuhao.New()
	} else if err != nil {
		return "", JSAPIPayParams{}, "", time.Time{}, errors.WarpQuick(err)
	} else if !u.Fuwuhao.Valid {
		return "", JSAPIPayParams{}, "", time.Time{}, MustBindFuwuhao.New()
	}

	OutTradeNoUUID, success := redis.GenerateUUIDMore(ctx, "pay", time.Minute*5, func(ctx context.Context, u uuid.UUID) bool {
		payModel := db.NewPayModel(mysql.MySQLConn)
		_, err := payModel.FindByPayID(ctx, WeChatPayID(u.String()))
		if errors.Is(err, db.ErrNotFound) {
			return true
		}

		return false
	})
	if !success {
		return "", JSAPIPayParams{}, "", time.Time{}, errors.Errorf("generate outtradeno fail")
	}

	OutTradeNo := WeChatPayID(OutTradeNoUUID.String())

	var get int64
	if couponsID != 0 {
		get, err = coupons.Recharge(ctx, couponsID, cny)
		if err != nil {
			return "", JSAPIPayParams{}, "", time.Time{}, errors.WarpQuick(err)
		}
	} else {
		get = cny
	}

	timeExpire := time.Now().Add(time.Second * time.Duration(config.BackendConfig.Coin.TimeExpireSec))

	bm := make(gopay.BodyMap)
	bm.Set("appid", config.BackendConfig.WeChatPay.AppID).
		Set("description", subject).
		Set("out_trade_no", OutTradeNo).
		Set("time_expire", timeExpire.Format(time.RFC3339)).
		Set("notify_url", config.BackendConfig.WeChatPay.ReturnURL)

	am := make(gopay.BodyMap)
	am.Set("total", cny).Set("currency", "CNY")

	bm.Set("amount", am)

	gd := make(gopay.BodyMap)
	gd.Set("goods_name", config.BackendConfig.Coin.Name).
		Set("merchant_goods_id", config.BackendConfig.Coin.ID).
		Set("quantity", cny).
		Set("unit_price", config.BackendConfig.Coin.Price)

	d := make(gopay.BodyMap)
	d.Set("goods_detail", []gopay.BodyMap{gd})

	bm.Set("detail", d)

	p := make(gopay.BodyMap)
	p.Set("openid", u.Fuwuhao.String)

	bm.Set("payer", p)

	res, err := WeChatPayClient.V3TransactionJsapi(ctx, bm)
	if err != nil {
		return "", JSAPIPayParams{}, "", time.Time{}, errors.WarpQuick(err)
	} else if res.Code != 0 {
		return "", JSAPIPayParams{}, "", time.Time{}, errors.Errorf("%s", res.Error)
	}

	resSign, err := WeChatPayClient.PaySignOfJSAPI(config.BackendConfig.FuWuHao.AppID, res.Response.PrepayId)
	if err != nil {
		return "", JSAPIPayParams{}, "", time.Time{}, errors.WarpQuick(err)
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	_, err = payModel.Insert(ctx, &db.Pay{
		UserId:   user.Id,
		WalletId: user.WalletId,
		PayId:    OutTradeNo,
		Subject:  subject,
		PayWay:   PayWayJSAPI,
		Cny:      cny,
		Get:      get,
		CouponsId: sql.NullInt64{
			Valid: couponsID != 0,
			Int64: couponsID,
		},
		TradeStatus: db.PayWait,
	})
	if err != nil {
		return "", JSAPIPayParams{}, "", time.Time{}, errors.WarpQuick(err)
	}

	audit.NewUserAudit(user.Id, "微信支付（%s）发起成功（%.2f）", PayWayJSAPI, float64(cny)/100.0)
	return res.Response.PrepayId, JSAPIPayParams{
		AppId:     resSign.AppId,
		TimeStamp: resSign.TimeStamp,
		NonceStr:  resSign.NonceStr,
		Package:   resSign.Package,
		SignType:  resSign.SignType,
		PaySign:   resSign.PaySign,
	}, OutTradeNo, timeExpire, nil
}

func NewRefund(ctx context.Context, user *db.User, pay *db.Pay) errors.WTError {
	if !config.BackendConfig.WeChatPay.UseReturnPay {
		return errors.Errorf("not ok")
	}

	if pay.PayWay != PayWayNative && pay.PayWay != PayWayH5 && pay.PayWay != PayWayJSAPI {
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
		Set("out_refund_no", pay.PayId).
		Set("reason", config.BackendConfig.Coin.RefundReason).
		Set("notify_url", config.BackendConfig.WeChatPay.ReturnURL)

	am := make(gopay.BodyMap)
	am.Set("refund", pay.Cny).
		Set("total", pay.Cny).
		Set("currency", "CNY")

	bm.Set("amount", am)

	gd := make(gopay.BodyMap)
	gd.Set("goods_name", config.BackendConfig.Coin.Name).
		Set("merchant_goods_id", config.BackendConfig.Coin.ID).
		Set("refund_amount", pay.Cny).
		Set("unit_price", config.BackendConfig.Coin.Price).
		Set("refund_quantity", pay.Cny)

	bm.Set("goods_detail", []gopay.BodyMap{gd})

	res, wxErr := WeChatPayClient.V3Refund(ctx, bm)
	if err != nil {
		logger.Logger.Error("微信支付退款失败：%s", wxErr.Error())
		_, _ = balance.PayRefundFail(ctx, user, pay, payStatus)
		return errors.WarpQuick(wxErr)
	} else if res.Code != 0 {
		if res.Code == 403 {
			logger.Logger.Error("微信支付退款余额不足，请注意。退款金额：%.2f", float64(pay.Cny)/100.0)
		} else {
			logger.Logger.Error("微信支付退款失败：%d", res.Code)
		}
		_, _ = balance.PayRefundFail(ctx, user, pay, payStatus)
		return errors.Errorf("system error")
	}

	mysqlErr := payModel.Update(ctx, pay)
	if mysqlErr != nil {
		return errors.WarpQuick(mysqlErr)
	}

	audit.NewUserAudit(pay.UserId, "微信退款（%s）发起成功（%.2f）", pay.PayWay, float64(pay.Cny)/100.0)
	return nil
}

func NewRefundInside(ctx context.Context, user *db.User, pay *db.Pay) errors.WTError {
	if !config.BackendConfig.WeChatPay.UseReturnPay {
		return errors.Errorf("not ok")
	}

	if pay.PayWay != PayWayNative && pay.PayWay != PayWayH5 && pay.PayWay != PayWayJSAPI {
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

	audit.NewUserAudit(pay.UserId, "微信单边退款（%s）发起成功（%.2f）", pay.PayWay, float64(pay.Cny)/100.0)
	return nil
}
