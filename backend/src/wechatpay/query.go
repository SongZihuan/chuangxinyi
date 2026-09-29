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
	"gitee.com/wuntsong-auth/backend/src/sender"
	"github.com/SuperH-0630/gopay"
	"github.com/SuperH-0630/gopay/wechat/v3"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

func QueryTrade(ctx context.Context, user *db.User, pay *db.Pay) (int64, errors.WTError) {
	if !config.BackendConfig.WeChatPay.UseNativePay && !config.BackendConfig.WeChatPay.UseH5Pay {
		return 0, errors.Errorf("not ok")
	}

	if pay.PayWay != PayWayNative && pay.PayWay != PayWayH5 && pay.PayWay != PayWayJSAPI {
		return 0, errors.Errorf("bad pay way")
	}

	if pay.TradeStatus != db.PayWait && pay.TradeStatus != db.PaySuccess {
		return pay.TradeStatus, nil
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return db.PayWait, nil
	}
	defer redis.ReleaseLock(keyPay)

	res, err := WeChatPayClient.V3TransactionQueryOrder(ctx, wechat.OutTradeNo, pay.PayId)
	if err != nil {
		return db.PayWait, nil
	}

	if res.Code != 0 {
		return db.PayWait, nil
	}

	switch res.Response.TradeState {
	default:
		pay.TradeStatus = db.PayWait
	case "NOTPAY":
		pay.TradeStatus = db.PayWait
	case "CLOSED":
		pay.TradeStatus = db.PayClose
	case "SUCCESS": // 只有PayWait才改变状态
		pay.BuyerId = sql.NullString{
			Valid:  true,
			String: res.Response.Payer.Openid,
		}
		pay.TradeNo = sql.NullString{
			Valid:  true,
			String: res.Response.TransactionId,
		}

		payTime, err := time.Parse(time.RFC3339, res.Response.SuccessTime)
		if err == nil {
			pay.PayAt = sql.NullTime{
				Valid: true,
				Time:  payTime,
			}
		}

		var get int64
		if pay.CouponsId.Valid {
			get, err = coupons.Recharge(ctx, pay.CouponsId.Int64, pay.Cny)
			if err != nil {
				get = pay.Get
			}
		} else {
			get = pay.Get
		}

		pay.Get = get

		if pay.TradeStatus == db.PayWait {
			_, err = balance.Pay(ctx, user, pay)
			if err != nil {
				logger.Logger.Error("add balance error: %s", err.Error())
				return 0, errors.WarpQuick(err)
			}
			logger.Logger.WXInfo("微信支付支付成功（通过检查）：%.2f", float64(pay.Cny)/100.0)
			_ = LogMsg(true, "微信支付支付成功（通过检查）：%.2f", float64(pay.Cny)/100.0)
		}

		sender.PhoneSendChange(pay.UserId, "余额（微信支付充值）")
		sender.EmailSendChange(pay.UserId, "余额（微信支付充值）")
		sender.MessageSendRecharge(pay.UserId, pay.Cny, "微信支付")
		sender.WxrobotSendRecharge(pay.UserId, pay.Cny, "微信支付")
		sender.FuwuhaoSendRecharge(pay)
		audit.NewUserAudit(pay.UserId, "微信支付充值已到账（%.2f）", float64(pay.Cny)/100.0)
	case "REFUND":
		pay.TradeStatus = db.PaySuccessRefund
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	err = payModel.Update(ctx, pay)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return 0, errors.WarpQuick(err)
	}

	return pay.TradeStatus, nil
}

func QueryRefund(ctx context.Context, user *db.User, pay *db.Pay) (int64, errors.WTError) {
	if !config.BackendConfig.WeChatPay.UseReturnPay {
		return 0, errors.Errorf("not ok")
	}

	if pay.PayWay != PayWayNative && pay.PayWay != PayWayH5 && pay.PayWay != PayWayJSAPI {
		return 0, errors.Errorf("bad pay way")
	}

	if pay.TradeStatus != db.PayWait {
		return pay.TradeStatus, nil
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return db.PayWaitRefund, nil
	}
	defer redis.ReleaseLock(keyPay)

	res, err := WeChatPayClient.V3RefundQuery(ctx, pay.PayId, make(gopay.BodyMap))
	if err != nil {
		return db.PayWaitRefund, nil
	}

	if res.Code != 0 {
		return db.PayWaitRefund, nil
	}

	switch res.Response.SuccessTime {
	case "SUCCESS":
		if pay.TradeStatus == db.PayWaitRefund {
			logger.Logger.WXInfo("微信支付退款成功（通过检查）：%.2f", float64(pay.Cny)/100.0)
			_ = LogMsg(true, "微信支付退款成功（通过检查）：%.2f", float64(pay.Cny)/100.0)
			sender.PhoneSendChange(pay.UserId, "余额（微信支付退款）")
			sender.EmailSendChange(pay.UserId, "余额（微信支付退款）")
			sender.MessageSendRefundPay(pay.UserId, pay.Cny, "微信支付")
			sender.WxrobotSendRefundPay(pay.UserId, pay.Cny, "微信支付")
			sender.FuwuhaoSendRefundPay(pay)
			audit.NewUserAudit(pay.UserId, "微信支付退款已到账（%.2f）", float64(pay.Cny)/100.0)
		}
		pay.TradeStatus = db.PaySuccessRefund
	case "CLOSE":
		pay.TradeStatus = db.PayCloseRefund
		_, err = balance.PayRefundFail(ctx, user, pay, db.PayCloseRefund)
	case "ABNORMAL":
		pay.TradeStatus = db.PayWaitRefund
		logger.Logger.Error("微信退款异常，需要前往处理：%s", pay.PayId)
	case "PROCESSING":
		pay.TradeStatus = db.PayWaitRefund
	default:
		pay.TradeStatus = db.PayWaitRefund
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	err = payModel.Update(ctx, pay)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return 0, errors.WarpQuick(err)
	}

	return pay.TradeStatus, nil
}
