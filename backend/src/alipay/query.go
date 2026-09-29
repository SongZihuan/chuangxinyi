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
	"gitee.com/wuntsong-auth/backend/src/sender"
	"github.com/SuperH-0630/gopay"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

func QueryTrade(ctx context.Context, user *db.User, pay *db.Pay) (int64, errors.WTError) {
	if !config.BackendConfig.Alipay.UsePCPay && !config.BackendConfig.Alipay.UseWapPay {
		return 0, errors.Errorf("not ok")
	}

	if pay.PayWay != PayWayWap && pay.PayWay != PayWayPC {
		return 0, errors.Errorf("bad pay way")
	}

	if pay.TradeStatus != db.PayWait && pay.TradeStatus != db.PaySuccess && pay.TradeStatus != db.PayWaitRefund && pay.TradeStatus != db.PayCloseRefund {
		return pay.TradeStatus, nil
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return db.PayWait, nil
	}
	defer redis.ReleaseLock(keyPay)

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", pay.PayId)

	res, err := AlipayClient.TradeQuery(ctx, bm)
	if err != nil {
		return db.PayWait, nil
	}

	if res.Response.Code != "10000" {
		return db.PayWait, nil
	}

	switch res.Response.TradeStatus {
	default:
		pay.TradeStatus = db.PayWait
	case "WAIT_BUYER_PAY":
		pay.TradeStatus = db.PayWait
	case "TRADE_CLOSED":
		pay.TradeStatus = db.PayClose
	case "TRADE_SUCCESS": // 除了PayWait不改变状态
		if len(res.Response.BuyerUserId) != 0 {
			pay.BuyerId = sql.NullString{
				Valid:  true,
				String: res.Response.BuyerUserId,
			}
		} else if len(res.Response.BuyerOpenId) != 0 {
			pay.BuyerId = sql.NullString{
				Valid:  true,
				String: res.Response.BuyerOpenId,
			}
		}

		payTime, err := time.Parse("2006-01-02 15:04:05", res.Response.SendPayDate)
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
			logger.Logger.WXInfo("支付宝支付成功（通过检查）：%.2f", float64(pay.Cny)/100.0)
			_ = LogMsg(true, "支付宝支付成功（通过检查）：%.2f", float64(pay.Cny)/100.0)

			sender.PhoneSendChange(pay.UserId, "余额（支付宝充值）")
			sender.EmailSendChange(pay.UserId, "余额（支付宝充值）")
			sender.MessageSendRecharge(pay.UserId, pay.Cny, "支付宝")
			sender.WxrobotSendRecharge(pay.UserId, pay.Cny, "支付宝")
			sender.FuwuhaoSendRecharge(pay)
			audit.NewUserAudit(pay.UserId, "支付宝充值已到账（%.2f）", float64(pay.Cny)/100.0)
		}
	case "TRADE_FINISHED":
		pay.TradeStatus = db.PayFinish
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	err = payModel.Update(ctx, pay)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return 0, errors.WarpQuick(err)
	}

	return pay.TradeStatus, nil
}

func QueryRefund(ctx context.Context, pay *db.Pay) (int64, errors.WTError) {
	if !config.BackendConfig.Alipay.UseReturnPay {
		return 0, errors.Errorf("not ok")
	}

	if pay.PayWay != PayWayWap && pay.PayWay != PayWayPC {
		return 0, errors.Errorf("bad pay way")
	}

	if pay.TradeStatus != db.PayWaitRefund {
		return pay.TradeStatus, nil
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return db.PayWaitRefund, nil
	}
	defer redis.ReleaseLock(keyPay)

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", pay.PayId)
	bm.Set("out_request_no", pay.PayId)

	res, err := AlipayClient.TradeFastPayRefundQuery(ctx, bm)
	if err != nil {
		return db.PayWaitRefund, nil
	}

	if res.Response.Code != "10000" {
		return db.PayWaitRefund, nil
	}

	if res.Response.RefundStatus == "REFUND_SUCCESS" {
		logger.Logger.WXInfo("支付宝退款成功（通过检查）：%.2f", float64(pay.Cny)/100.0)
		_ = LogMsg(true, "支付宝退款成功（通过检查）：%.2f", float64(pay.Cny)/100.0)
		pay.TradeStatus = db.PaySuccessRefund

		sender.PhoneSendChange(pay.UserId, "余额（支付宝退款）")
		sender.EmailSendChange(pay.UserId, "余额（支付宝退款）")
		sender.MessageSendRefundPay(pay.UserId, pay.Cny, "支付宝")
		sender.WxrobotSendRefundPay(pay.UserId, pay.Cny, "支付宝")
		sender.FuwuhaoSendRefundPay(pay)
		audit.NewUserAudit(pay.UserId, "支付宝退款已到账（%.2f）", float64(pay.Cny)/100.0)
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	err = payModel.Update(ctx, pay)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return 0, errors.WarpQuick(err)
	}

	return pay.TradeStatus, nil
}
