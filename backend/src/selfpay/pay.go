package selfpay

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/alipay"
	"gitee.com/wuntsong-auth/backend/src/audit"
	"gitee.com/wuntsong-auth/backend/src/balance"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/coupons"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/redis"
	"gitee.com/wuntsong-auth/backend/src/sender"
	"gitee.com/wuntsong-auth/backend/src/wechatpay"
	"github.com/wuntsong-org/wterrors"
	"time"
)

var BadStatus = errors.NewClass("bad status")

func Pay(ctx context.Context, user *db.User, pay *db.Pay, get int64) errors.WTError {
	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap || pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		return BadPayWay.New()
	}

	if pay.TradeStatus != db.PayWait {
		return BadStatus.New()
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return errors.Errorf("pay lock error")
	}
	defer redis.ReleaseLock(keyPay)

	pay.PayAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	if get >= 0 {
		pay.Get = get
	} else { // 使用优惠券
		var get int64
		if pay.CouponsId.Valid && pay.Get == pay.Cny {
			var err error
			get, err = coupons.Recharge(ctx, pay.CouponsId.Int64, pay.Cny)
			if err != nil {
				get = pay.Get
			}
		} else {
			get = pay.Get
		}

		pay.Get = get
	}

	var err error
	_, err = balance.Pay(ctx, user, pay)
	if err != nil {
		logger.Logger.Error("add balance error: %s", err.Error())
		return errors.WarpQuick(err)
	}

	payModel := db.NewPayModel(mysql.MySQLConn)
	err = payModel.Update(ctx, pay)
	if err != nil {
		logger.Logger.Error("mysql error: %s", err.Error())
		return errors.WarpQuick(err)
	}

	logger.Logger.WXInfo("用户录入自支付：%.2f", float64(pay.Cny)/100.0)
	_ = LogMsg(true, "用户录入自支付：%.2f", float64(pay.Cny)/100.0)
	sender.PhoneSendChange(pay.UserId, "余额（用户自充值）")
	sender.EmailSendChange(pay.UserId, "余额（用户自充值）")
	sender.MessageSendRecharge(pay.UserId, pay.Cny, "用户自充值")
	sender.WxrobotSendRecharge(pay.UserId, pay.Cny, "用户自充值")
	sender.FuwuhaoSendRecharge(pay)
	audit.NewUserAudit(pay.UserId, "用户自支付审核成功（%.2f）", float64(pay.Cny)/100.00)

	return nil
}

func PayFail(ctx context.Context, pay *db.Pay) errors.WTError {
	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap || pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		return BadPayWay.New()
	}

	if pay.TradeStatus != db.PayWait {
		return BadStatus.New()
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return errors.Errorf("pay lock error")
	}
	defer redis.ReleaseLock(keyPay)

	pay.TradeStatus = db.PayClose

	payModel := db.NewPayModel(mysql.MySQLConn)
	mysqlErr := payModel.Update(ctx, pay)
	if mysqlErr != nil {
		logger.Logger.Error("mysql error: %s", mysqlErr.Error())
		return errors.WarpQuick(mysqlErr)
	}

	audit.NewUserAudit(pay.UserId, "用户自支付审核失败（%.2f）", float64(pay.Cny)/100.00)

	return nil
}

func Refund(ctx context.Context, user *db.User, pay *db.Pay) errors.WTError { // 不限制支付方式
	if pay.TradeStatus == db.PayWait || pay.TradeStatus == db.PayClose || pay.TradeStatus == db.PaySuccessRefund || pay.TradeStatus == db.PaySuccessRefundInside {
		return BadStatus.New()
	}

	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap {
		if pay.TradeStatus == db.PayWaitRefund { // 证明已经开始退款，人工不能再干预
			return BadStatus.New()
		}

		// 不需要上锁
		err := alipay.NewRefund(ctx, user, pay)
		if err != nil {
			return err
		}
	} else if pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		if pay.TradeStatus == db.PayWaitRefund { // 证明已经开始退款，人工不能再干预
			return BadStatus.New()
		}

		// 不需要上锁
		err := wechatpay.NewRefund(ctx, user, pay)
		if err != nil {
			return err
		}
	} else { // 人工
		keyPay := fmt.Sprintf("pay:%s", pay.PayId)
		if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
			return errors.Errorf("pay lock error")
		}
		defer redis.ReleaseLock(keyPay)

		if pay.TradeStatus == db.PayWaitRefund { // 此前已经申请并执行过扣款
			// 什么都不做
		} else {
			_, err := balance.PayRefund(ctx, user, pay)
			if errors.Is(err, balance.Insufficient) {
				return Insufficient.New()
			} else if err != nil {
				return errors.WarpQuick(err)
			}

			sender.PhoneSendChange(pay.UserId, "余额（管理员人工审核退款）")
			sender.EmailSendChange(pay.UserId, "余额（管理员人工审核退款）")
			sender.MessageSendRefundPay(pay.UserId, pay.Cny, "管理员人工审核退款")
			sender.WxrobotSendRefundPay(pay.UserId, pay.Cny, "管理员人工审核退款")
			sender.FuwuhaoSendRefundPay(pay)
		}

		pay.TradeStatus = db.PaySuccessRefund

		payModel := db.NewPayModel(mysql.MySQLConn)
		mysqlErr := payModel.Update(ctx, pay)
		if mysqlErr != nil {
			logger.Logger.Error("mysql error: %s", mysqlErr.Error())
			return errors.WarpQuick(mysqlErr)
		}

		audit.NewUserAudit(pay.UserId, "管理员人工审核退款成功（%.2f）", float64(pay.Cny)/100.00)
		logger.Logger.WXInfo("管理员人工审核退款：%.2f", float64(pay.Cny)/100.0)
	}

	return nil
}

func RefundInside(ctx context.Context, user *db.User, pay *db.Pay) errors.WTError { // 不限制支付方式
	if pay.TradeStatus == db.PayWait || pay.TradeStatus == db.PayClose || pay.TradeStatus == db.PaySuccessRefund || pay.TradeStatus == db.PaySuccessRefundInside {
		return BadStatus.New()
	}

	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap {
		if pay.TradeStatus == db.PayWaitRefund { // 证明已经开始退款，人工不能再干预
			return BadStatus.New()
		}

		// 不需要上锁
		err := alipay.NewRefundInside(ctx, user, pay)
		if err != nil {
			return err
		}
	} else if pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		if pay.TradeStatus == db.PayWaitRefund { // 证明已经开始退款，人工不能再干预
			return BadStatus.New()
		}

		// 不需要上锁
		err := wechatpay.NewRefundInside(ctx, user, pay)
		if err != nil {
			return err
		}
	} else { // 人工
		keyPay := fmt.Sprintf("pay:%s", pay.PayId)
		if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
			return errors.Errorf("pay lock error")
		}
		defer redis.ReleaseLock(keyPay)

		if pay.TradeStatus == db.PayWaitRefund { // 此前已经申请并执行过扣款
			// 什么都不做
		} else {
			_, err := balance.PayRefund(ctx, user, pay)
			if errors.Is(err, balance.Insufficient) {
				return Insufficient.New()
			} else if err != nil {
				return errors.WarpQuick(err)
			}

			sender.PhoneSendChange(pay.UserId, "余额（管理员人工审核退款）")
			sender.EmailSendChange(pay.UserId, "余额（管理员人工审核退款）")
			sender.MessageSendRefundPay(pay.UserId, pay.Cny, "管理员人工审核退款")
			sender.WxrobotSendRefundPay(pay.UserId, pay.Cny, "管理员人工审核退款")
			sender.FuwuhaoSendRefundPay(pay)
		}

		pay.TradeStatus = db.PaySuccessRefund // 人工没有inside模式

		payModel := db.NewPayModel(mysql.MySQLConn)
		mysqlErr := payModel.Update(ctx, pay)
		if mysqlErr != nil {
			logger.Logger.Error("mysql error: %s", mysqlErr.Error())
			return errors.WarpQuick(mysqlErr)
		}

		audit.NewUserAudit(pay.UserId, "管理员人工审核退款成功（%.2f）", float64(pay.Cny)/100.00)
		logger.Logger.WXInfo("管理员人工审核退款：%.2f", float64(pay.Cny)/100.0)
	}

	return nil
}

func RefundFail(ctx context.Context, pay *db.Pay) errors.WTError { // 不限制支付方式
	if pay.TradeStatus == db.PayWait || pay.TradeStatus == db.PayClose || pay.TradeStatus == db.PaySuccessRefund || pay.TradeStatus == db.PaySuccessRefundInside {
		return BadStatus.New()
	}

	if pay.PayWay == alipay.PayWayPC || pay.PayWay == alipay.PayWayWap || pay.PayWay == wechatpay.PayWayNative || pay.PayWay == wechatpay.PayWayH5 || pay.PayWay == wechatpay.PayWayJSAPI {
		if pay.TradeStatus == db.PayWaitRefund { // 证明已经开始退款，人工不能再干预
			return BadStatus.New()
		}
	}

	if pay.TradeStatus != db.PayWaitRefund {
		return errors.Errorf("bad status")
	}

	keyPay := fmt.Sprintf("pay:%s", pay.PayId)
	if !redis.AcquireLockMore(ctx, keyPay, time.Minute*2) {
		return errors.Errorf("pay lock error")
	}
	defer redis.ReleaseLock(keyPay)

	pay.TradeStatus = db.PayCloseRefund

	payModel := db.NewPayModel(mysql.MySQLConn)
	mysqlErr := payModel.Update(ctx, pay)
	if mysqlErr != nil {
		logger.Logger.Error("mysql error: %s", mysqlErr.Error())
		return errors.WarpQuick(mysqlErr)
	}

	audit.NewUserAudit(pay.UserId, "管理员人工审核退款失败（%.2f）", float64(pay.Cny)/100.00)

	return nil
}

func LogMsg(atall bool, text string, args ...any) errors.WTError {
	return logger.WxRobotSendNotRecord(config.BackendConfig.WXRobot.PayLog, fmt.Sprintf(text, args...), atall)
}
